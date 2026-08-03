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
		vpcName := "wireguard24"
		vpcCidr := "10.8.0.0"
		vmName := "wireguard24"
		vmIP := "10.8.1.145"

		// If stack contains "test", use alternate values
		if strings.Contains(stackName, "test") {
			vpcName = "wireguard-test"
			vpcCidr = "10.15.1.0"
			vmName = "wireguard-test"
			vmIP = "10.15.1.145"
		}

		// Set WIREGUARD_DOMAIN env var from pulumi config (or env var)
		wireguardDomain := cfg.Get("wireguard_domain")
		if wireguardDomain != "" {
			if err := os.Setenv("WIREGUARD_DOMAIN", wireguardDomain); err != nil {
				return err
			}
		}

		// Set SECRET_OPERATOR_AUTHENTICATION_TOKEN from pulumi config if provided
		secretToken := cfg.Get("secret_operator_authentication_token")
		if secretToken != "" {
			if err := os.Setenv("SECRET_OPERATOR_AUTHENTICATION_TOKEN", secretToken); err != nil {
				return err
			}
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

		// Get snapshot ID from config, defaults to ubuntu-26.04 if not specified
		snapshotID := cfg.Get("wireguard_snapshot_id")
		imageName := "ubuntu-26.04"
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

		vm, err := compute.CreateWireguardVM(ctx, computeArgs, vmIP)

		if err != nil {
			return err
		}

		sshConnector := shared.WireguardProvisioner(ctx, keyPair)

		//TODO implement exporting of mutliptl ssh output with one session
		_ = compute.ProvisionVM(ctx, vmName, &model.ProvisionArgs{
			ExportName:    "wireguard.publicKey",
			SourceCompute: vm,
		}, &sshConnector)

		sshConnectorPassword := shared.WireguardPasswordProvisioner(ctx, keyPair)
		_ = compute.ProvisionVM(ctx, vmName, &model.ProvisionArgs{
			ExportName:    "wireguard.password",
			SourceCompute: vm,
		}, &sshConnectorPassword)

		return err
	})
}
