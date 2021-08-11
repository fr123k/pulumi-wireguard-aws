package compute

import (
	"fmt"
	"os"
	"strings"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const size = "t2.micro"

//CreateWireguardVM creates a wireguard ec2 aws instance

func CreateSecurityGroups(ctx *pulumi.Context, computeArgs *model.ComputeArgs) ([]*model.SecurityGroup, error) {
	for _, securityGroup := range computeArgs.SecurityGroups {
		securityGroupArgs := &ec2.SecurityGroupArgs{
			Description: pulumi.String(securityGroup.Description),
			Ingress: network.IngressRules(computeArgs.Security, securityGroup.IngressRules),
			Egress: network.EgressRules(computeArgs.Security, securityGroup.EgressRules),
			Tags: pulumi.ToStringMap(securityGroup.Tags),
		}
		if computeArgs.Vpc != nil {
			securityGroupArgs.VpcId = computeArgs.Vpc.ID()
		}
		sgExternal, err := ec2.NewSecurityGroup(ctx, securityGroup.Name, securityGroupArgs)
		if err != nil {
			return nil, err
		}

		securityGroup.State = sgExternal.CustomResourceState
	}
	return computeArgs.SecurityGroups, nil
}

func CreateWireguardVM(ctx *pulumi.Context, computeArgs *model.ComputeArgs) (*model.ComputeResult, error) {
	CreateSecurityGroups(ctx, computeArgs)

	mostRecent := true
	ami, err := aws.GetAmi(ctx, &aws.GetAmiArgs{
		Filters: []aws.GetAmiFilter{
			{
				Name:   "name",
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*"},
			},
		},
		Owners:     []string{"099720109477"},
		MostRecent: &mostRecent,
	})

	if err != nil {
		return nil, err
	}

	ami2, err := aws.GetAmiIds(ctx, &aws.GetAmiIdsArgs{
		Filters: []aws.GetAmiIdsFilter{
			{
				Name: "name",
				Values: []string{
					"wireguard-ami",
				},
			},
			{
				Name: "state",
				Values: []string{
					"available",
				},
			},
		},
		// MostRecent: &mostRecent,
		Owners: []string{
			"self",
		},
	}, nil)

	if err != nil {
		return nil, err
	}

	var amiID string
	if ami2.Ids != nil && len(ami2.Ids) > 0 {
		amiID = ami2.Ids[0]
	} else {
		amiID = ami.Id
	}

	//TODO cloud-init use only if jenkins ami doesn't exists.
	// yaml, err := getCloudInitYaml("cloud-init/cloud-init.yaml", awsKeyID, awsKeySecret)
	userDataVariables := map[string]string{
		"{{ CLIENT_PUBLICKEY }}":        "CLIENT_PUBLICKEY",
		"{{ CLIENT_IP_ADDRESS }}":       "CLIENT_IP_ADDRESS",
		"{{ MAILJET_API_CREDENTIALS }}": "MAILJET_API_CREDENTIALS",
		"{{ METADATA_URL }}":            "METADATA_URL",
	}
	userData, err := model.NewUserData("cloud-init/user-data.txt", model.TemplateVariablesEnvironment(userDataVariables))
	if err != nil {
		return nil, err
	}

	ctx.Export("cloud-init", pulumi.String(userData.Content))

	keyPair, err := ec2.NewKeyPair(ctx, *computeArgs.KeyPair.Name, &ec2.KeyPairArgs{
		KeyName:   pulumi.String(*computeArgs.KeyPair.Name),
		PublicKey: pulumi.String(*computeArgs.KeyPair.SSHKeyPair.PublicKeyStr),
	})

	if err != nil {
		return nil, err
	}

	wireguardEc2Args := &ec2.InstanceArgs{
		AssociatePublicIpAddress: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name":   pulumi.String("wireguard"),
			"JobUrl": pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
		},
		InstanceType: pulumi.String(size),
		KeyName:      keyPair.KeyName,
		Ami:          pulumi.String(amiID),
		UserData:     pulumi.String(userData.Content),

		VpcSecurityGroupIds: network.ToStringArray(computeArgs.SecurityGroups),
	}

	if computeArgs.Vpc != nil {
		wireguardEc2Args.SubnetId = computeArgs.Vpc.SubnetResults[0].ID()
	}

	server, err := ec2.NewInstance(ctx, "wireguard", wireguardEc2Args)

	ctx.Export("publicIp", server.PublicIp)
	ctx.Export("publicDns", server.PublicDns)

	return &model.ComputeResult{
		Compute: server.CustomResourceState,
	}, err
}

func ProvisionVM(ctx *pulumi.Context, provisionArgs *model.ProvisionArgs, actor actors.Connector) error {

	server, err := ec2.GetInstance(ctx, "wireguard2", provisionArgs.SourceCompute.ID(), &ec2.InstanceState{
		InstanceState: pulumi.String("running"),
	})

	if err != nil {
		return err
	}

	provision := server.PublicIp.ApplyT(func(hostip string) string {
		var result string
		if actor != nil {
			result = actor.Connect(hostip)
			defer actor.Stop()
		}
		return strings.TrimSuffix(result, "\r\n")
	})

	ctx.Export(provisionArgs.ExportName, provision)

	return nil
}

// CreateImage creates an virtual machine image from an running VM.
func CreateImage(ctx *pulumi.Context, imageArgs model.ImageArgs, actor actors.Connector) error {

	server, err := ec2.GetInstance(ctx, "wireguard2", imageArgs.SourceCompute.ID(), &ec2.InstanceState{
		InstanceState: pulumi.String("running"),
	})

	if err != nil {
		return err
	}

	provision := server.PublicIp.ApplyT(func(hostip string) string {
		var result string
		if actor != nil {
			result = actor.Connect(hostip)
			defer actor.Stop()
		}

		//TODO implement the NewAmiFromInstance logic as an actor as well

		_, err = ec2.NewAmiFromInstance(ctx, imageArgs.Name, &ec2.AmiFromInstanceArgs{
			SourceInstanceId:      imageArgs.SourceCompute.ID(),
			Name:                  pulumi.String(imageArgs.Name),
			SnapshotWithoutReboot: pulumi.Bool(false),
		}, pulumi.IgnoreChanges([]string{"sourceInstanceId"}))

		if err != nil {
			panic(fmt.Errorf("failed to create ami image: %s", err))
		}

		return result
	})

	ctx.Export("Provisioning", provision)

	return nil
}
