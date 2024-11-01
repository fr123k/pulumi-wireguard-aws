package network

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/creasty/defaults"
)

// CreateVPC creates a aws VPC resource
func CreateVPC(ctx *pulumi.Context, vpcArgs *model.VpcArgs) (*model.VpcResult, error) {
	defaults.MustSet(vpcArgs)

	vpc, err := ec2.NewVpc(ctx, vpcArgs.Name, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(vpcArgs.Cidr),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
		InstanceTenancy:    pulumi.String(vpcArgs.InstanceTenancy),
	})
	if err != nil {
		return nil, err
	}

	// Export IDs of the created resources to the Pulumi stack
	ctx.Export("vpcId", vpc.ID())

	internetGW, err := ec2.NewInternetGateway(ctx, vpcArgs.Name, &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
	})

	ec2.NewRoute(ctx, vpcArgs.Name, &ec2.RouteArgs{
		RouteTableId:         vpc.MainRouteTableId,
		DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
		GatewayId:            internetGW.ID(),
	})

	if err != nil {
		return nil, err
	}

	subnets := make([]model.SubnetResult, len(vpcArgs.Subnets))
	for i, subnetArg := range vpcArgs.Subnets {
		subnet, err := ec2.NewSubnet(ctx, vpcArgs.Name, &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: pulumi.String(subnetArg.Cidr),
		})

		if err != nil {
			return nil, err
		}

		ctx.Export("subnetId", subnet.ID())
		subnets[i] = model.SubnetResult{
			Subnet: subnet.CustomResourceState,
		}
	}

	vpcResult := &model.VpcResult{
		Vpc:           vpc.CustomResourceState,
		SubnetResults: subnets,
	}

	return vpcResult, nil
}
