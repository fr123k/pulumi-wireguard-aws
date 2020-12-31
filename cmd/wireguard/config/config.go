package config

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
)

var VPCArgs = vpcArg("wireguard", "10.8.0.0")

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
