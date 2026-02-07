package compute

import (
	"regexp"
	"strings"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// isPrebakedImage detects if the image is a pre-baked Hetzner snapshot.
// Returns true if the image name is numeric (snapshot ID) or starts with "prebaked" prefix.
func isPrebakedImage(imageName string) bool {
	// Check if it's a numeric snapshot ID
	numericPattern := regexp.MustCompile(`^\d+$`)
	if numericPattern.MatchString(imageName) {
		return true
	}

	// Check if it starts with "prebaked" or "temporal-prebaked" prefix
	lowerName := strings.ToLower(imageName)
	return strings.HasPrefix(lowerName, "prebaked") || strings.HasPrefix(lowerName, "temporal-prebaked")
}

// const size = "cx11"

// type Infrastructure struct {
// 	Server    *hcloud.Server
// 	ImageName *string
// 	UserData  *model.UserData
// }

// type exportsFnc = func(ctx *pulumi.Context, infra *Infrastructure)

// // type exportsFnc = func(ctx *pulumi.Context, infra *Infrastructure)

// func exports(ctx *pulumi.Context, infra *Infrastructure) {
// 	ctx.Export("publicIp", infra.Server.Ipv4Address)
// 	ctx.Export("publicDns", infra.Server.Ipv4Address)
// }

// func CreateServer(ctx *pulumi.Context, computeArgs *model.ComputeArgs, export exportsFnc) (*Infrastructure, error) {
// 	if computeArgs.UserData != nil {
// 		ctx.Export("cloud-init", pulumi.String(computeArgs.UserData.Content))
// 	}

// 	var serverKeys pulumi.StringArray
// 	if computeArgs.KeyPair.SSHKeyPair != nil {
// 		sshKey, err := hcloud.NewSshKey(ctx, *computeArgs.KeyPair.Name, &hcloud.SshKeyArgs{
// 			Name:      pulumi.String(*computeArgs.KeyPair.Name),
// 			PublicKey: pulumi.String(*computeArgs.KeyPair.SSHKeyPair.PublicKeyStr),
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		serverKeys = pulumi.StringArray{
// 			sshKey.ID(),
// 		}
// 	} else {
// 		nameSelector := fmt.Sprintf("Name=%s", *computeArgs.KeyPair.Name)
// 		sshKeys, err := hcloud.GetSshKeys(ctx, &hcloud.GetSshKeysArgs{
// 			WithSelector: &nameSelector,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		if len(sshKeys.SshKeys) > 0 {
// 			serverKeys = pulumi.StringArray{
// 				pulumi.String(strconv.Itoa(sshKeys.SshKeys[0].Id)),
// 			}
// 		} else {
// 			return nil, fmt.Errorf("ssh keys not specified")
// 		}
// 	}

// 	serverArgs := hcloud.ServerArgs{
// 		//TODO handle multiple images like in the aws modul
// 		Image: pulumi.String(computeArgs.Images[0].Name),
// 		Labels: pulumi.StringMap{
// 			"Name": pulumi.String(computeArgs.Name),
// 		},
// 		Location:   pulumi.String("nbg1"),
// 		Name:       pulumi.String(computeArgs.Name),
// 		ServerType: pulumi.String(size),
// 		SshKeys:    serverKeys,
// 		UserData:   pulumi.String(computeArgs.UserData.Content),
// 		Backups:    pulumi.Bool(false),
// 	}

// 	if computeArgs.UserData != nil {
// 		serverArgs.UserData = pulumi.String(computeArgs.UserData.Content)
// 	}

// 	server, err := hcloud.NewServer(ctx, computeArgs.Name, &serverArgs, pulumi.IgnoreChanges([]string{
// 		"firewallIds",
// 		"allowDeprecatedImages",
// 		"deleteProtection",
// 		"ignoreRemoteFirewallIds",
// 	}))

// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = hcloud.NewServerNetwork(ctx, "srvnetwork", &hcloud.ServerNetworkArgs{
// 		ServerId:  utility.IDtoInt(server.CustomResourceState),
// 		NetworkId: computeArgs.Vpc.IDtoInt(),
// 		Ip:        pulumi.String("10.8.0.145"),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	infra := Infrastructure{
// 		Server:   server,
// 		UserData: computeArgs.UserData,
// 	}

// 	export(ctx, &infra)

// 	return &infra, nil
// }

// CreateTemporalVM creates a Temporal VM on Hetzner Cloud.
// If the image is a pre-baked snapshot (numeric ID or "prebaked" prefix), it uses the minimal
// cloud-init script. Otherwise, it uses the full cloud-init script for base Ubuntu images.
func CreateTemporalVM(ctx *pulumi.Context, computeArgs *model.ComputeArgs, vmIP string) (*model.ComputeResult, error) {
	var userData *model.UserData
	var err error

	// Detect if using a pre-baked snapshot image
	imageName := ""
	if len(computeArgs.Images) > 0 && computeArgs.Images[0] != nil {
		imageName = computeArgs.Images[0].Name
	}

	if isPrebakedImage(imageName) {
		ctx.Log.Info("Using pre-baked image, loading minimal cloud-init", nil)
		userData, err = shared.TemporalPrebakedUserData()
	} else {
		ctx.Log.Info("Using base image, loading full cloud-init", nil)
		userData, err = shared.TemporalUserData()
	}

	if err != nil {
		return nil, err
	}

	computeArgs.UserData = userData

	infra, err := CreateServer(ctx, computeArgs, vmIP, exports)
	if err != nil {
		return nil, err
	}

	//TODO hetzner cloud doesn't support security rules but the same can be achieved with local firewalls with in the VM
	//     Implement firewall provisioning based on userdata script or cloud-init.

	return &model.ComputeResult{
		Compute: infra.Server.CustomResourceState,
	}, nil
}

// func ProvisionVM(ctx *pulumi.Context, provisionArgs *model.ProvisionArgs, actor actors.Connector) error {
// 	name := fmt.Sprintf("wireguard-%s", utility.RandomSecret(12))
// 	server, err := hcloud.GetServer(ctx, name, provisionArgs.SourceCompute.ID(), &hcloud.ServerState{
// 		Status: pulumi.String("running"),
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	provision := server.Ipv4Address.ApplyT(func(hostip string) string {
// 		var result string
// 		if actor != nil {
// 			result = actor.Connect(hostip)
// 			defer actor.Stop()
// 		}
// 		return strings.TrimSuffix(result, "\r\n")
// 	})

// 	ctx.Export(provisionArgs.ExportName, provision)

// 	return nil
// }
