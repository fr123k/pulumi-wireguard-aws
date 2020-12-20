package main

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/network"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	// "github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

const size = "t2.large"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// config := config.New(ctx, "")

		// awsKeyID := config.Require("key")
		// awsKeySecret := config.Require("secret")
		vpc, subnet, err := network.CreateVPC(ctx)
		if err != nil {
			return err
		}
		return compute.CreateWireguardVM(ctx, vpc, subnet)
	})
}
