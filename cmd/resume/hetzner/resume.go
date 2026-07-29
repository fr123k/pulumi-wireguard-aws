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
		vpcName := "resume"
		vpcCidr := "10.13.1.0"
		vmName := "resume"
		vmIP := "10.13.1.145"

		// If stack contains "test", use alternate values
		if strings.Contains(stackName, "test") {
			vpcName = "resume-test"
			vpcCidr = "10.14.1.0"
			vmName = "resume-test"
			vmIP = "10.14.1.145"
		}

		security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), model.VpcArg(vpcName, vpcCidr))
		security.Println()

		vpc, err := network.CreateVPC(ctx, model.VpcArg(vpcName, vpcCidr))
		if err != nil {
			return err
		}
		keyPairName := vmName + "-"

		var keyPair *model.KeyPairArgs

		keyFile := cfg.Get("ssh_key_file")

		if _, err := os.Stat(keyFile); err == nil {
			keyPair = model.NewKeyPairArgsWithPrivateKeyFile(&keyPairName, keyFile)
			fmt.Printf("Use local ssh key file %s\n", keyFile)
		} else {
			keyPair = model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
			fmt.Println("Use random ssh key")
		}

		keyPair.Name = &keyPairName
		keyPair.Username = "frank.ittermann"

		// Use pre-baked snapshot ID if configured, otherwise default to ubuntu-24.04
		imageName := cfg.Get("resume_snapshot_id")
		if imageName == "" {
			imageName = "ubuntu-24.04"
		}

		computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
		computeArgs.Name = vmName
		computeArgs.Images = []*model.ImageArgs{
			{
				Name: imageName,
			},
		}

		vm, err := compute.CreateResumeVM(ctx, computeArgs, vmIP)

		if err != nil {
			return err
		}

		sshConnector := shared.ResumeProvisioner(ctx, keyPair)

		return compute.ProvisionVM(ctx, vmName, &model.ProvisionArgs{
			ExportName:    "resume.result",
			SourceCompute: vm,
		}, &sshConnector)
	})
}