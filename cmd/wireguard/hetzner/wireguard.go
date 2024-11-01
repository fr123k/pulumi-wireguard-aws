package main

import (
	"fmt"
	"os"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), model.VPCArgsDefault)
		security.Println()

		vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
		if err != nil {
			return err
		}
		keyPairName := "wireguard24-"

		var keyPair *model.KeyPairArgs

		kefFile := cfg.Get("ssh_key_file")

		if _, err := os.Stat(kefFile); err == nil {
			//Uncomment to enable ssh access for debugging
			keyPair = model.NewKeyPairArgsWithPrivateKeyFile(&keyPairName, kefFile)
			fmt.Printf("Use local ssh key file %s\n", kefFile)
		} else {
			keyPair = model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
			fmt.Println("Use random ssh key")
		}

		keyPair.Name = &keyPairName
		keyPair.Username = "frank.ittermann"

		// keyPair.Username = "frank.ittermann"
		computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
		computeArgs.Name = "wireguard24"
		computeArgs.Images = []*model.ImageArgs{
			{
				Name: "ubuntu-24.04",
			},
		}

		vm, err := compute.CreateWireguardVM(ctx, computeArgs)

		if err != nil {
			return err
		}

		sshConnector := shared.WireguardProvisioner(ctx, keyPair)

		//TODO implement exporting of mutliptl ssh output with one session
		compute.ProvisionVM(ctx, &model.ProvisionArgs{
			ExportName:    "wireguard.publicKey",
			SourceCompute: vm,
		}, &sshConnector)

		sshConnectorPassword := shared.WireguardPasswordProvisioner(ctx, keyPair)
		compute.ProvisionVM(ctx, &model.ProvisionArgs{
			ExportName:    "wireguard.password",
			SourceCompute: vm,
		}, &sshConnectorPassword)

		return err
	})
}
