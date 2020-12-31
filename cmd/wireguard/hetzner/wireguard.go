package main

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
	"github.com/fr123k/pulumi-wireguard-aws/cmd/wireguard/config"


	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	// "github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

const size = "t2.large"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// config := config.New(ctx, "")

		vpc, err := network.CreateVPC(ctx, config.VPCArgs)
		if err != nil {
			return err
		}
		// return nil
		return compute.CreateWireguardVM(ctx, vpc)
	})
}
