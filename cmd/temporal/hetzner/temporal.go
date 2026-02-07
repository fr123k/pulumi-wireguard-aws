package main

import (
	"fmt"
	"os"
	"strings"

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

		stackName := ctx.Stack()

		// Default values for production
		vpcName := "temporal"
		vpcCidr := "10.9.1.0"
		vmName := "temporal"
		vmIP := "10.9.1.145"

		// If stack contains "test", use alternate values
		if strings.Contains(stackName, "test") {
			vpcName = "temporal-test"
			vpcCidr = "10.10.1.0"
			vmName = "temporal-test"
			vmIP = "10.10.1.145"
		}

		security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), model.VpcArg(vpcName, vpcCidr))
		security.Println()

		vpc, err := network.CreateVPC(ctx, model.VpcArg(vpcName, vpcCidr))
		if err != nil {
			return err
		}
		keyPairName := vmName + "-"

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

		// Get snapshot ID from config, defaults to ubuntu-24.04 if not specified
		snapshotID := cfg.Get("temporal_snapshot_id")
		imageName := "ubuntu-24.04"
		if snapshotID != "" {
			imageName = snapshotID
		}

		computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
		computeArgs.Name = vmName
		computeArgs.Images = []*model.ImageArgs{
			{
				Name: imageName,
			},
		}

		vm, err := compute.CreateTemporalVM(ctx, computeArgs, vmIP)

		if err != nil {
			return err
		}

		sshConnector := shared.TemporalProvisioner(ctx, keyPair)

		//TODO implement exporting of mutliptl ssh output with one session
		compute.ProvisionVM(ctx, vmName, &model.ProvisionArgs{
			ExportName:    "wireguard.publicKey",
			SourceCompute: vm,
		}, &sshConnector)

		return err
	})
}
