package main

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	// "github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

const size = "t2.large"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// config := config.New(ctx, "")

		// awsKeyID := config.Require("key")
		// awsKeySecret := config.Require("secret")
		vpcArgs := &model.VpcArgs{
			Name: "wireguard",
			Cidr: "10.8.0.0/16",
			Subnets: []model.SubnetArgs {{
					Cidr: "10.8.0.0/24",
				},
			},
		}
		vpc, err := network.CreateVPC(ctx, vpcArgs)
		if err != nil {
			return err
		}
		// return nil
		return compute.CreateWireguardVM(ctx, vpc)
	})
}
