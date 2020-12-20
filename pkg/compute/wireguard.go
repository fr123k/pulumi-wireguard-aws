package compute

import (
	"os"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const size = "t2.large"

//CreateWireguardVM creates a wireguard ec2 aws instance
func CreateWireguardVM(ctx *pulumi.Context, vpc *ec2.Vpc, subnet *ec2.Subnet) error {
	sgExternal, err := ec2.NewSecurityGroup(ctx, "wireguard-external", &ec2.SecurityGroupArgs{
		Description: pulumi.String("Terraform Managed. Allow Wireguard client traffic from internet."),
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("udp"),
				FromPort:   pulumi.Int(51820),
				ToPort:     pulumi.Int(51820),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(22),
				ToPort:     pulumi.Int(22),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
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
		VpcId: vpc.ID(),
	})
	if err != nil {
		return err
	}

	sgAdmin, err := ec2.NewSecurityGroup(ctx, "wireguard-admin", &ec2.SecurityGroupArgs{
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
		VpcId: vpc.ID(),
	})
	if err != nil {
		return err
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
		return err
	}

	//TODO cloud-init use only if jenkins ami doesn't exists.
	// yaml, err := getCloudInitYaml("cloud-init/cloud-init.yaml", awsKeyID, awsKeySecret)
	yaml, err := utility.GetUserData("cloud-init/user-data.txt")

	if err != nil {
		return err
	}

	ctx.Export("cloud-init", pulumi.String(*yaml))

	publicKey, err := utility.ReadFile("keys/wireguard.pem.pub")

	if err != nil {
		return err
	}

	keyPair, err := ec2.NewKeyPair(ctx, "wireguard", &ec2.KeyPairArgs{
		KeyName:   pulumi.String("wireguard"),
		PublicKey: pulumi.String(*publicKey),
	})

	if err != nil {
		return err
	}

	server, err := ec2.NewInstance(ctx, "wireguard", &ec2.InstanceArgs{
		AssociatePublicIpAddress: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name":   pulumi.String("wireguard"),
			"JobUrl": pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
		},
		InstanceType: pulumi.String(size),
		KeyName:      keyPair.KeyName, //create the keypair with pulumi
		Ami:          pulumi.String(ami.Id),
		UserData:     pulumi.String(*yaml),
		SubnetId:     subnet.ID(),

		VpcSecurityGroupIds: pulumi.StringArray{
			sgExternal.ID(), sgAdmin.ID(),
		},
	})

	ctx.Export("publicIp", server.PublicIp)
	ctx.Export("publicDns", server.PublicDns)

	return err
}
