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
		security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), wireguardCfg.VPCArgs)
		security.Println()

		vpc, err := network.CreateVPC(ctx, wireguardCfg.VPCArgs)
		if err != nil {
			return err
		}

		return compute.CreateWireguardVM(ctx, vpc, security)
	})
}