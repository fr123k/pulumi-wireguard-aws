// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provides information for multiple EC2 Local Gateways, such as their identifiers.
//
// ## Example Usage
//
// The following example retrieves Local Gateways with a resource tag of `service` set to `production`.
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		fooLocalGateways, err := ec2.GetLocalGateways(ctx, &ec2.GetLocalGatewaysArgs{
// 			Tags: map[string]interface{}{
// 				"service": "production",
// 			},
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		ctx.Export("foo", fooLocalGateways.Ids)
// 		return nil
// 	})
// }
// ```
func GetLocalGateways(ctx *pulumi.Context, args *GetLocalGatewaysArgs, opts ...pulumi.InvokeOption) (*GetLocalGatewaysResult, error) {
	var rv GetLocalGatewaysResult
	err := ctx.Invoke("aws:ec2/getLocalGateways:getLocalGateways", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getLocalGateways.
type GetLocalGatewaysArgs struct {
	// Custom filter block as described below.
	Filters []GetLocalGatewaysFilter `pulumi:"filters"`
	// A mapping of tags, each pair of which must exactly match
	// a pair on the desired local_gateways.
	Tags map[string]string `pulumi:"tags"`
}

// A collection of values returned by getLocalGateways.
type GetLocalGatewaysResult struct {
	Filters []GetLocalGatewaysFilter `pulumi:"filters"`
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// Set of all the Local Gateway identifiers
	Ids  []string          `pulumi:"ids"`
	Tags map[string]string `pulumi:"tags"`
}
