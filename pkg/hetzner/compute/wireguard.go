package compute

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const size = "cx11"

type Infrastructure struct {
	Server    *hcloud.Server
	ImageName *string
	UserData  *model.UserData
}

type exportsFnc = func(ctx *pulumi.Context, infra *Infrastructure)

// type exportsFnc = func(ctx *pulumi.Context, infra *Infrastructure)

func exports(ctx *pulumi.Context, infra *Infrastructure) {
	ctx.Export("publicIp", infra.Server.Ipv4Address)
	ctx.Export("publicDns", infra.Server.Ipv4Address)
}

func CreateServer(ctx *pulumi.Context, computeArgs *model.ComputeArgs, export exportsFnc) (*Infrastructure, error) {
	if computeArgs.UserData != nil {
		ctx.Export("cloud-init", pulumi.String(computeArgs.UserData.Content))
	}

	var serverKeys pulumi.StringArray
	if computeArgs.KeyPair.SSHKeyPair != nil {
		sshKey, err := hcloud.NewSshKey(ctx, *computeArgs.KeyPair.Name, &hcloud.SshKeyArgs{
			Name:      pulumi.String(*computeArgs.KeyPair.Name),
			PublicKey: pulumi.String(*computeArgs.KeyPair.SSHKeyPair.PublicKeyStr),
		})
		if err != nil {
			return nil, err
		}
		serverKeys = pulumi.StringArray{
			sshKey.ID(),
		}
	} else {
		nameSelector := fmt.Sprintf("Name=%s", *computeArgs.KeyPair.Name)
		sshKeys, err := hcloud.GetSshKeys(ctx, &hcloud.GetSshKeysArgs{
			WithSelector: &nameSelector,
		})
		if err != nil {
			return nil, err
		}
		if len(sshKeys.SshKeys) > 0 {
			serverKeys = pulumi.StringArray{
				pulumi.String(strconv.Itoa(sshKeys.SshKeys[0].Id)),
			}
		} else {
			return nil, fmt.Errorf("ssh keys not specified")
		}
	}

	serverArgs := hcloud.ServerArgs{
		//TODO handle multiple images like in the aws modul
		Image: pulumi.String(computeArgs.Images[0].Name),
		Labels: pulumi.Map{
			"Name": pulumi.String(computeArgs.Name),
		},
		Location:   pulumi.String("nbg1"),
		Name:       pulumi.String(computeArgs.Name),
		ServerType: pulumi.String(size),
		SshKeys:    serverKeys,
		UserData:   pulumi.String(computeArgs.UserData.Content),
		Backups:    pulumi.Bool(false),
	}

	if computeArgs.UserData != nil {
		serverArgs.UserData = pulumi.String(computeArgs.UserData.Content)
	}

	server, err := hcloud.NewServer(ctx, computeArgs.Name, &serverArgs, pulumi.IgnoreChanges([]string{"firewallIds"}))

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

	infra := Infrastructure{
		Server:   server,
		UserData: computeArgs.UserData,
	}

	export(ctx, &infra)

	return &infra, nil
}

// CreateWireguardVM creates a wireguard ec2 aws instance
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

	userData, err := shared.WireguardUserData()
	if err != nil {
		return nil, err
	}

	computeArgs.UserData = userData

	infra, err := CreateServer(ctx, computeArgs, exports)

	if err != nil {
		return nil, err
	}

	//TODO hetzner cloud doesn't support security rules but the same can be achieved with local firewalls with in the VM
	//     Implement firewall provisioning based on userdata script or cloud-init.

	return &model.ComputeResult{
		Compute: infra.Server.CustomResourceState,
	}, nil
}

func ProvisionVM(ctx *pulumi.Context, provisionArgs *model.ProvisionArgs, actor actors.Connector) error {
	name := fmt.Sprintf("wireguard-%s", utility.RandomSecret(12))
	server, err := hcloud.GetServer(ctx, name, provisionArgs.SourceCompute.ID(), &hcloud.ServerState{
		Status: pulumi.String("running"),
	})

	if err != nil {
		return err
	}

	provision := server.Ipv4Address.ApplyT(func(hostip string) string {
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
