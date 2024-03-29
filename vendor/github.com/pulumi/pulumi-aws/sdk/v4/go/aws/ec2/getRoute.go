// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// `ec2.Route` provides details about a specific Route.
//
// This resource can prove useful when finding the resource associated with a CIDR. For example, finding the peering connection associated with a CIDR value.
//
// ## Example Usage
//
// The following example shows how one might use a CIDR value to find a network interface id and use this to create a data source of that network interface.
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		cfg := config.New(ctx, "")
// 		subnetId := cfg.RequireObject("subnetId")
// 		opt0 := subnetId
// 		_, err := ec2.LookupRouteTable(ctx, &ec2.LookupRouteTableArgs{
// 			SubnetId: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		opt1 := "10.0.1.0/24"
// 		route, err := ec2.LookupRoute(ctx, &ec2.LookupRouteArgs{
// 			RouteTableId:         aws_route_table.Selected.Id,
// 			DestinationCidrBlock: &opt1,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		opt2 := route.NetworkInterfaceId
// 		_, err = ec2.LookupNetworkInterface(ctx, &ec2.LookupNetworkInterfaceArgs{
// 			Id: &opt2,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupRoute(ctx *pulumi.Context, args *LookupRouteArgs, opts ...pulumi.InvokeOption) (*LookupRouteResult, error) {
	var rv LookupRouteResult
	err := ctx.Invoke("aws:ec2/getRoute:getRoute", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getRoute.
type LookupRouteArgs struct {
	// EC2 Carrier Gateway ID of the Route belonging to the Route Table.
	CarrierGatewayId *string `pulumi:"carrierGatewayId"`
	// CIDR block of the Route belonging to the Route Table.
	DestinationCidrBlock *string `pulumi:"destinationCidrBlock"`
	// IPv6 CIDR block of the Route belonging to the Route Table.
	DestinationIpv6CidrBlock *string `pulumi:"destinationIpv6CidrBlock"`
	// The ID of a managed prefix list destination of the Route belonging to the Route Table.
	DestinationPrefixListId *string `pulumi:"destinationPrefixListId"`
	// Egress Only Gateway ID of the Route belonging to the Route Table.
	EgressOnlyGatewayId *string `pulumi:"egressOnlyGatewayId"`
	// Gateway ID of the Route belonging to the Route Table.
	GatewayId *string `pulumi:"gatewayId"`
	// Instance ID of the Route belonging to the Route Table.
	InstanceId *string `pulumi:"instanceId"`
	// Local Gateway ID of the Route belonging to the Route Table.
	LocalGatewayId *string `pulumi:"localGatewayId"`
	// NAT Gateway ID of the Route belonging to the Route Table.
	NatGatewayId *string `pulumi:"natGatewayId"`
	// Network Interface ID of the Route belonging to the Route Table.
	NetworkInterfaceId *string `pulumi:"networkInterfaceId"`
	// The ID of the specific Route Table containing the Route entry.
	RouteTableId string `pulumi:"routeTableId"`
	// EC2 Transit Gateway ID of the Route belonging to the Route Table.
	TransitGatewayId *string `pulumi:"transitGatewayId"`
	// VPC Peering Connection ID of the Route belonging to the Route Table.
	VpcPeeringConnectionId *string `pulumi:"vpcPeeringConnectionId"`
}

// A collection of values returned by getRoute.
type LookupRouteResult struct {
	CarrierGatewayId         string `pulumi:"carrierGatewayId"`
	DestinationCidrBlock     string `pulumi:"destinationCidrBlock"`
	DestinationIpv6CidrBlock string `pulumi:"destinationIpv6CidrBlock"`
	DestinationPrefixListId  string `pulumi:"destinationPrefixListId"`
	EgressOnlyGatewayId      string `pulumi:"egressOnlyGatewayId"`
	GatewayId                string `pulumi:"gatewayId"`
	// The provider-assigned unique ID for this managed resource.
	Id                     string `pulumi:"id"`
	InstanceId             string `pulumi:"instanceId"`
	LocalGatewayId         string `pulumi:"localGatewayId"`
	NatGatewayId           string `pulumi:"natGatewayId"`
	NetworkInterfaceId     string `pulumi:"networkInterfaceId"`
	RouteTableId           string `pulumi:"routeTableId"`
	TransitGatewayId       string `pulumi:"transitGatewayId"`
	VpcPeeringConnectionId string `pulumi:"vpcPeeringConnectionId"`
}
