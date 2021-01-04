package compute

import (
	"fmt"
	"os"
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/ssh"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const size = "t2.micro"

//CreateWireguardVM creates a wireguard ec2 aws instance
func CreateWireguardVM(ctx *pulumi.Context, computeArgs *model.ComputeArgs) (*model.ComputeResult,error) {
	wireguardExtSecGroupArgs := &ec2.SecurityGroupArgs{
		Description: pulumi.String("Pulumi Managed. Allow Wireguard client traffic from internet."),
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("udp"),
				FromPort:   pulumi.Int(51820),
				ToPort:     pulumi.Int(51820),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
			network.SSHIngressRule(computeArgs.Security),
		},
		Egress: ec2.SecurityGroupEgressArray{
			ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("-1"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
		Tags: pulumi.StringMap{
			"JobUrl":         pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
			"Project":        pulumi.String("wireguard"),
			"pulumi-managed": pulumi.String("True"),
		},
	}
	if computeArgs.Vpc != nil {
		wireguardExtSecGroupArgs.VpcId = computeArgs.Vpc.ID()
	}

	sgExternal, err := ec2.NewSecurityGroup(ctx, "wireguard-external", wireguardExtSecGroupArgs)
	if err != nil {
		return nil, err
	}

	wireguardAdminSecGroupArgs := &ec2.SecurityGroupArgs{
		Description: pulumi.String("Terraform Managed. Allow admin traffic internal resources from VPN"),
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:       pulumi.String("-1"),
				FromPort:       pulumi.Int(0),
				ToPort:         pulumi.Int(0),
				SecurityGroups: pulumi.StringArray{sgExternal.ID()},
			},
			ec2.SecurityGroupIngressArgs{
				Protocol:       pulumi.String("icmp"),
				FromPort:       pulumi.Int(8),
				ToPort:         pulumi.Int(0),
				SecurityGroups: pulumi.StringArray{sgExternal.ID()},
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("-1"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
		Tags: pulumi.StringMap{
			"JobUrl":         pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
			"Project":        pulumi.String("wireguard"),
			"pulumi-managed": pulumi.String("True"),
		},
	}

	if computeArgs.Vpc != nil {
		wireguardAdminSecGroupArgs.VpcId = computeArgs.Vpc.ID()
	}

	sgAdmin, err := ec2.NewSecurityGroup(ctx, "wireguard-admin", wireguardAdminSecGroupArgs)
	if err != nil {
		return nil, err
	}

	mostRecent := true
	//TODO check if jenkins master jocker ami exists use it otherwise use this one.
	//make this behaviour configurable always use the following ami except following cases
	// 1) jenkins jocker ami exists 2) 1) && env var JENKINS_AMI=ami
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
	if (ami2.Ids != nil && len(ami2.Ids) > 0) {
		amiID = ami2.Ids[0]
	} else {
		amiID = ami.Id
	}

	//TODO cloud-init use only if jenkins ami doesn't exists.
	// yaml, err := getCloudInitYaml("cloud-init/cloud-init.yaml", awsKeyID, awsKeySecret)
	userDataVariables := map[string]string{
		"{{ CLIENT_PUBLICKEY }}": "CLIENT_PUBLICKEY",
		"{{ CLIENT_IP_ADDRESS }}": "CLIENT_IP_ADDRESS",
		"{{ MAILJET_API_CREDENTIALS }}": "MAILJET_API_CREDENTIALS",
		"{{ METADATA_URL }}": "METADATA_URL",
	}
	userData, err := model.NewUserData("cloud-init/user-data.txt", model.TemplateVariablesEnvironment(userDataVariables))
	if err != nil {
		return nil, err
	}

	ctx.Export("cloud-init", pulumi.String(userData.Content))

	publicKey, err := utility.ReadFile("keys/wireguard.pem.pub")

	if err != nil {
		return nil, err
	}

	// randSrc := rand.NewSource(time.Now().UnixNano())
	// keyPairName := fmt.Sprintf("wireguard-%d", randSrc.New(s1).Intn(100000))
	keyPairName := "wireguard"
	keyPair, err := ec2.NewKeyPair(ctx, keyPairName, &ec2.KeyPairArgs{
		KeyName:   pulumi.String(keyPairName),
		PublicKey: pulumi.String(*publicKey),
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
		
		VpcSecurityGroupIds: pulumi.StringArray{
			sgExternal.ID(), sgAdmin.ID(),
		},
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



// CreateImage creates an virtual machine image from an running VM.
func CreateImage(ctx *pulumi.Context, imageArgs model.ImageArgs) (error) {

	server, err := ec2.GetInstance(ctx, "wireguard2", imageArgs.SourceCompute.ID(), &ec2.InstanceState{
		InstanceState: pulumi.String("running"),
	})

	if err != nil {
		return err
	}

	provision := server.PublicIp.ApplyString(func(hostip string) string {
		sshClient := ssh.SSHClientConfig{
			Hostname: hostip,
			Port: 22,
			Username: "ubuntu",
			PrivateKeyFileName: "/Users/franki/private/github/pulumi-wireguard-aws/keys/wireguard.pem",
			Timeout: 2 * time.Minute,
		}

		fmt.Printf("Open SSH connection to %s", hostip)

		result, err := sshClient.SSHCommand("sudo cloud-init status --wait")
		if err != nil {
			panic(fmt.Errorf("Failed to create session: %s", err))
		}
		fmt.Printf("Result: %s", *result)

		_, err = ec2.NewAmiFromInstance(ctx, imageArgs.Name, &ec2.AmiFromInstanceArgs{
			SourceInstanceId: imageArgs.SourceCompute.ID(),
			Name: pulumi.String(imageArgs.Name),
			SnapshotWithoutReboot: pulumi.Bool(false),
		}, pulumi.IgnoreChanges([]string{"sourceInstanceId"}))

		if err != nil {
			panic(fmt.Errorf("Failed to create Ami Image: %s", err))
		}

		return *result
	})

	ctx.Export("Provisioning", provision)

	return nil
}
