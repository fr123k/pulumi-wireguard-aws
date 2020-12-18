package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	// "github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

const size = "t2.large"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// config := config.New(ctx, "")

		// awsKeyID := config.Require("key")
		// awsKeySecret := config.Require("secret")
		vpc, subnet, err := createVPC(ctx)
		if err != nil {
			return err
		}
		return createWireguardVM(ctx, vpc, subnet)
	})
}

func createVPC(ctx *pulumi.Context) (*ec2.Vpc, *ec2.Subnet, error) {
	vpc, err := ec2.NewVpc(ctx, "wireguard", &ec2.VpcArgs{
		CidrBlock:          pulumi.String("10.8.0.0/16"),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
		InstanceTenancy:    pulumi.String("default"),
	})
	if err != nil {
		return nil, nil, err
	}

	// Export IDs of the created resources to the Pulumi stack
	ctx.Export("VPC-ID", vpc.ID())

	internetGW, err := ec2.NewInternetGateway(ctx, "wireguard", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
	})

	ec2.NewRoute(ctx, "wireguard", &ec2.RouteArgs{
		RouteTableId: vpc.MainRouteTableId,
		DestinationCidrBlock:   pulumi.String("0.0.0.0/0"),
		GatewayId: internetGW.ID(),
	})

	if err != nil {
		return nil, nil, err
	}

	subnet, err := ec2.NewSubnet(ctx, "wireguard", &ec2.SubnetArgs{
		VpcId:     vpc.ID(),
		CidrBlock: pulumi.String("10.8.0.0/24"),
	})

	if err != nil {
		return nil, nil, err
	}

	// Export IDs of the created resources to the Pulumi stack
	ctx.Export("Subnet-ID", subnet.ID())
	return vpc, subnet, nil
}

func readFile(fileName string) (*string, error) {
	b, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		return nil, err
	}
	yaml := string(b)
	return &yaml, nil
}

func getUserData(fileName string) (*string, error) {
	data, err := readFile(fileName)
	if err != nil {
		return nil, err
	}
	yaml := parseUserData(*data)
	return &yaml, nil
}

func parseUserData(content string) string {
	clientPublicKey, ok := os.LookupEnv("CLIENT_PUBLICKEY")
	var result string
	if ok == true {
		result = strings.ReplaceAll(content, "{{ CLIENT_PUBLICKEY }}", clientPublicKey)
	} else {
		result = strings.ReplaceAll(content, "{{ CLIENT_PUBLICKEY }}", "")
	}
	return result
}

func createWireguardVM(ctx *pulumi.Context, vpc *ec2.Vpc, subnet *ec2.Subnet) error {

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
	yaml, err := getUserData("cloud-init/user-data.txt")

	if err != nil {
		return err
	}

	ctx.Export("cloud-init", pulumi.String(*yaml))

	publicKey, err := readFile("keys/wireguard.pem.pub")

	if err != nil {
		return err
	}

	keyPair, err := ec2.NewKeyPair(ctx, "wireguard", &ec2.KeyPairArgs{
		KeyName: pulumi.String("wireguard"),
		PublicKey: pulumi.String(*publicKey),
	})

	if err != nil {
		return err
	}

	server, err := ec2.NewInstance(ctx, "wireguard", &ec2.InstanceArgs{
		AssociatePublicIpAddress: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("wireguard"),
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
