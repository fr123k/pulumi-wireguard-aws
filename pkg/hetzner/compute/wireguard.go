package compute

import (
	"strings"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

const size = "cx11"

//CreateWireguardVM creates a wireguard ec2 aws instance
func CreateWireguardVM(ctx *pulumi.Context, computeArgs *model.ComputeArgs) (*model.ComputeResult, error) {

	/*
			// Enable or disable backups.
		Backups pulumi.BoolPtrInput
		// The datacenter name to create the server in.
		Datacenter pulumi.StringPtrInput
		// Name or ID of the image the server is created from.
		Image pulumi.StringInput
		// ID or Name of an ISO image to mount.
		Iso pulumi.StringPtrInput
		// If true, do not upgrade the disk. This allows downgrading the server type later.
		KeepDisk pulumi.BoolPtrInput
		// User-defined labels (key-value pairs) should be created with.
		Labels pulumi.MapInput
		// The location name to create the server in. `nbg1`, `fsn1` or `hel1`
		Location pulumi.StringPtrInput
		// Name of the server to create (must be unique per project and a valid hostname as per RFC 1123).
		Name pulumi.StringPtrInput
		// Enable and boot in to the specified rescue system. This enables simple installation of custom operating systems. `linux64` `linux32` or `freebsd64`
		Rescue pulumi.StringPtrInput
		// Name of the server type this server should be created with.
		ServerType pulumi.StringInput
		// SSH key IDs or names which should be injected into the server at creation time
		SshKeys pulumi.StringArrayInput
		// Cloud-Init user data to use during server creation
		UserData pulumi.StringPtrInput
	*/

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

	sshKey, err := hcloud.NewSshKey(ctx, *computeArgs.KeyPair.Name, &hcloud.SshKeyArgs{
		Name:      pulumi.String(*computeArgs.KeyPair.Name),
		PublicKey: pulumi.String(*computeArgs.KeyPair.SSHKeyPair.PublicKeyStr),
	})

	if err != nil {
		return nil, err
	}

	server, err := hcloud.NewServer(ctx, "wireguard", &hcloud.ServerArgs{
		Image:      pulumi.String("ubuntu-20.04"),
		Location:   pulumi.String("nbg1"),
		Name:       pulumi.String("wireguard"),
		ServerType: pulumi.String(size),
		SshKeys: pulumi.StringArray{
			sshKey.ID(),
		},
		UserData: pulumi.String(userData.Content),
	})

	if err != nil {
		return nil, err
	}

	_, err = hcloud.NewServerNetwork(ctx, "srvnetwork", &hcloud.ServerNetworkArgs{
		ServerId:  utility.IDtoInt(server.CustomResourceState),
		NetworkId: computeArgs.Vpc.IDtoInt(),
		Ip:        pulumi.String("10.8.0.145"),
	})
	if err != nil {
		return nil, err
	}

	ctx.Export("publicIp", server.Ipv4Address)
	ctx.Export("publicDns", server.Ipv4Address)

	//TODO hetzner cloud doesn't support security rules but the same can be achieved with local firewalls with in the VM
	//     Implement firewall provisioning based on userdata script or cloud-init.

	return &model.ComputeResult{
		Compute: server.CustomResourceState,
	}, err

	// sgExternal, err := ec2.NewSecurityGroup(ctx, "wireguard-external", &ec2.SecurityGroupArgs{
	// 	Description: pulumi.String("Terraform Managed. Allow Wireguard client traffic from internet."),
	// 	Ingress: ec2.SecurityGroupIngressArray{
	// 		ec2.SecurityGroupIngressArgs{
	// 			Protocol:   pulumi.String("udp"),
	// 			FromPort:   pulumi.Int(51820),
	// 			ToPort:     pulumi.Int(51820),
	// 			CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
	// 		},
	// 		ec2.SecurityGroupIngressArgs{
	// 			Protocol:   pulumi.String("tcp"),
	// 			FromPort:   pulumi.Int(22),
	// 			ToPort:     pulumi.Int(22),
	// 			CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
	// 		},
	// 	},
	// 	Egress: ec2.SecurityGroupEgressArray{
	// 		ec2.SecurityGroupEgressArgs{
	// 			Protocol:   pulumi.String("-1"),
	// 			FromPort:   pulumi.Int(0),
	// 			ToPort:     pulumi.Int(0),
	// 			CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
	// 		},
	// 	},
	// 	Tags: pulumi.StringMap{
	// 		"JobUrl":         pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
	// 		"Project":        pulumi.String("wireguard"),
	// 		"pulumi-managed": pulumi.String("True"),
	// 	},
	// 	VpcId: vpc.ID(),
	// })
	// if err != nil {
	// 	return err
	// }

	// sgAdmin, err := ec2.NewSecurityGroup(ctx, "wireguard-admin", &ec2.SecurityGroupArgs{
	// 	Description: pulumi.String("Terraform Managed. Allow admin traffic internal resources from VPN"),
	// 	Ingress: ec2.SecurityGroupIngressArray{
	// 		ec2.SecurityGroupIngressArgs{
	// 			Protocol:       pulumi.String("-1"),
	// 			FromPort:       pulumi.Int(0),
	// 			ToPort:         pulumi.Int(0),
	// 			SecurityGroups: pulumi.StringArray{sgExternal.ID()},
	// 		},
	// 		ec2.SecurityGroupIngressArgs{
	// 			Protocol:       pulumi.String("icmp"),
	// 			FromPort:       pulumi.Int(8),
	// 			ToPort:         pulumi.Int(0),
	// 			SecurityGroups: pulumi.StringArray{sgExternal.ID()},
	// 		},
	// 	},
	// 	Egress: ec2.SecurityGroupEgressArray{
	// 		ec2.SecurityGroupEgressArgs{
	// 			Protocol:   pulumi.String("-1"),
	// 			FromPort:   pulumi.Int(0),
	// 			ToPort:     pulumi.Int(0),
	// 			CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
	// 		},
	// 	},
	// 	Tags: pulumi.StringMap{
	// 		"JobUrl":         pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
	// 		"Project":        pulumi.String("wireguard"),
	// 		"pulumi-managed": pulumi.String("True"),
	// 	},
	// 	VpcId: vpc.ID(),
	// })
	// if err != nil {
	// 	return err
	// }

	// mostRecent := true
	// //TODO check if jenkins master jocker ami exists use it otherwise use this one.
	// //make this behaviour configurable always use the following ami except following cases
	// // 1) jenkins jocker ami exists 2) 1) && env var JENKINS_AMI=ami
	// ami, err := aws.GetAmi(ctx, &aws.GetAmiArgs{
	// 	Filters: []aws.GetAmiFilter{
	// 		{
	// 			Name:   "name",
	// 			Values: []string{"ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*"},
	// 		},
	// 	},
	// 	Owners:     []string{"099720109477"},
	// 	MostRecent: &mostRecent,
	// })

	// if err != nil {
	// 	return err
	// }

	// //TODO cloud-init use only if jenkins ami doesn't exists.
	// // yaml, err := getCloudInitYaml("cloud-init/cloud-init.yaml", awsKeyID, awsKeySecret)
	// yaml, err := utility.GetUserData("cloud-init/user-data.txt")

	// if err != nil {
	// 	return err
	// }

	// ctx.Export("cloud-init", pulumi.String(*yaml))

	// publicKey, err := utility.ReadFile("keys/wireguard.pem.pub")

	// if err != nil {
	// 	return err
	// }

	// keyPair, err := ec2.NewKeyPair(ctx, "wireguard", &ec2.KeyPairArgs{
	// 	KeyName:   pulumi.String("wireguard"),
	// 	PublicKey: pulumi.String(*publicKey),
	// })

	// if err != nil {
	// 	return err
	// }

	// server, err := ec2.NewInstance(ctx, "wireguard", &ec2.InstanceArgs{
	// 	AssociatePublicIpAddress: pulumi.Bool(true),
	// 	Tags: pulumi.StringMap{
	// 		"Name":   pulumi.String("wireguard"),
	// 		"JobUrl": pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
	// 	},
	// 	InstanceType: pulumi.String(size),
	// 	KeyName:      keyPair.KeyName, //create the keypair with pulumi
	// 	Ami:          pulumi.String(ami.Id),
	// 	UserData:     pulumi.String(*yaml),
	// 	SubnetId:     vpc.SubnetResults[0].ID(),

	// 	VpcSecurityGroupIds: pulumi.StringArray{
	// 		sgExternal.ID(), sgAdmin.ID(),
	// 	},
	// })

	// ctx.Export("publicIp", server.PublicIp)
	// ctx.Export("publicDns", server.PublicDns)

	// return err
}

func ProvisionVM(ctx *pulumi.Context, provisionArgs *model.ProvisionArgs, actor actors.Connector) error {

	server, err := hcloud.GetServer(ctx, "wireguard2", provisionArgs.SourceCompute.ID(), &hcloud.ServerState{
		Status: pulumi.String("running"),
	})

	if err != nil {
		return err
	}

	provision := server.Ipv4Address.ApplyString(func(hostip string) string {
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
