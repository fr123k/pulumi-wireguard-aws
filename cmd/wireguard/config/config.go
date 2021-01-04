package config

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
)

// VPCArgsDefault deffine default arguments for the VPC
var VPCArgsDefault = vpcArg("wireguard", "10.8.0.0")

func vpcArg(name string, cidr string) (*model.VpcArgs) {
	vpcArgs := &model.VpcArgs{
		Name: name,
		Cidr: fmt.Sprintf("%s/16", cidr),
		Subnets: []model.SubnetArgs {{
				Cidr: fmt.Sprintf("%s/24", cidr),
			},
		},
	}
	defaults.MustSet(vpcArgs)
	return vpcArgs
}

//Todo precalculate subnet cidrs look at https://play.golang.org/p/m8TNTtygK0
