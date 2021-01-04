package main

import (
	wireguardCfg "github.com/fr123k/pulumi-wireguard-aws/cmd/wireguard/config"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
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

		vpc, err := network.CreateVPC(ctx, wireguardCfg.VPCArgsDefault)
		if err != nil {
			return err
		}
		keyPairName := "wireguard-"
		_, err = compute.CreateWireguardVM(ctx, model.NewComputeArgsWithKeyPair(vpc, security, model.NewKeyPairArgsWithRandomName(&keyPairName)))
		return err
	})
}
