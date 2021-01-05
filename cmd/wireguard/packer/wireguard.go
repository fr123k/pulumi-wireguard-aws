package main

import (
	"time"

	wireguardCfg "github.com/fr123k/pulumi-wireguard-aws/cmd/wireguard/config"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

const size = "t2.large"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), wireguardCfg.VPCArgsDefault)
		security.Println()
		
		keyPairName := "wireguard-"
		vm, err := compute.CreateWireguardVM(ctx, model.NewComputeArgsWithSecurityAndKeyPair(security, model.NewKeyPairArgsWithRandomName(&keyPairName)))

		if err != nil {
			return err
		}

		sshConnector := actors.NewSSHConnector(
			actors.SSHConnectorArgs{
				Port: 22,
				Username: "ubuntu",
				Timeout: 2 * time.Minute,
				PrivateKeyFileName: "/Users/franki/private/github/pulumi-wireguard-aws/keys/wireguard.pem",
			},
		)

		err = compute.CreateImage(ctx, model.ImageArgs{
			Name: "wireguard-ami-new",
			SourceCompute: vm,
		},	&sshConnector)
		return err
	})
}
