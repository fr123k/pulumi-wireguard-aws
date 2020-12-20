package network

import (
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

//CreateVPC creates a aws VPC resource
func CreateVPC(ctx *pulumi.Context) (*ec2.Vpc, *ec2.Subnet, error) {
	vpc, err := ec2.NewVpc(ctx, "wireguard", &ec2.VpcArgs{
		CidrBlock:          pulumi.String("10.8.0.0/16"),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
		InstanceTenancy:    pulumi.String("default"),
	})
	if err != nil {
		return nil, nil, err
	}

	// Export IDs of the created resources to the Pulumi stack
	ctx.Export("VPC-ID", vpc.ID())

	internetGW, err := ec2.NewInternetGateway(ctx, "wireguard", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
	})

	ec2.NewRoute(ctx, "wireguard", &ec2.RouteArgs{
		RouteTableId: vpc.MainRouteTableId,
		DestinationCidrBlock:   pulumi.String("0.0.0.0/0"),
		GatewayId: internetGW.ID(),
	})

	if err != nil {
		return nil, nil, err
	}

	subnet, err := ec2.NewSubnet(ctx, "wireguard", &ec2.SubnetArgs{
		VpcId:     vpc.ID(),
		CidrBlock: pulumi.String("10.8.0.0/24"),
	})

	if err != nil {
		return nil, nil, err
	}

	// Export IDs of the created resources to the Pulumi stack
	ctx.Export("Subnet-ID", subnet.ID())
	return vpc, subnet, nil
}
