// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// This resource can be useful for getting back a list of VPC Ids for a region.
//
// The following example retrieves a list of VPC Ids with a custom tag of `service` set to a value of "production".
//
// ## Example Usage
//
// The following shows outputting all VPC Ids.
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			foo, err := ec2.GetVpcs(ctx, &ec2.GetVpcsArgs{
//				Tags: map[string]interface{}{
//					"service": "production",
//				},
//			}, nil)
//			if err != nil {
//				return err
//			}
//			ctx.Export("foo", foo.Ids)
//			return nil
//		})
//	}
//
// ```
//
// An example use case would be interpolate the `ec2.getVpcs` output into `count` of an ec2.FlowLog resource.
func GetVpcs(ctx *pulumi.Context, args *GetVpcsArgs, opts ...pulumi.InvokeOption) (*GetVpcsResult, error) {
	opts = internal.PkgInvokeDefaultOpts(opts)
	var rv GetVpcsResult
	err := ctx.Invoke("aws:ec2/getVpcs:getVpcs", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getVpcs.
type GetVpcsArgs struct {
	// Custom filter block as described below.
	//
	// More complex filters can be expressed using one or more `filter` sub-blocks,
	// which take the following arguments:
	Filters []GetVpcsFilter `pulumi:"filters"`
	// Map of tags, each pair of which must exactly match
	// a pair on the desired vpcs.
	Tags map[string]string `pulumi:"tags"`
}

// A collection of values returned by getVpcs.
type GetVpcsResult struct {
	Filters []GetVpcsFilter `pulumi:"filters"`
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// List of all the VPC Ids found.
	Ids  []string          `pulumi:"ids"`
	Tags map[string]string `pulumi:"tags"`
}

func GetVpcsOutput(ctx *pulumi.Context, args GetVpcsOutputArgs, opts ...pulumi.InvokeOption) GetVpcsResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (GetVpcsResultOutput, error) {
			args := v.(GetVpcsArgs)
			opts = internal.PkgInvokeDefaultOpts(opts)
			var rv GetVpcsResult
			secret, err := ctx.InvokePackageRaw("aws:ec2/getVpcs:getVpcs", args, &rv, "", opts...)
			if err != nil {
				return GetVpcsResultOutput{}, err
			}

			output := pulumi.ToOutput(rv).(GetVpcsResultOutput)
			if secret {
				return pulumi.ToSecret(output).(GetVpcsResultOutput), nil
			}
			return output, nil
		}).(GetVpcsResultOutput)
}

// A collection of arguments for invoking getVpcs.
type GetVpcsOutputArgs struct {
	// Custom filter block as described below.
	//
	// More complex filters can be expressed using one or more `filter` sub-blocks,
	// which take the following arguments:
	Filters GetVpcsFilterArrayInput `pulumi:"filters"`
	// Map of tags, each pair of which must exactly match
	// a pair on the desired vpcs.
	Tags pulumi.StringMapInput `pulumi:"tags"`
}

func (GetVpcsOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*GetVpcsArgs)(nil)).Elem()
}

// A collection of values returned by getVpcs.
type GetVpcsResultOutput struct{ *pulumi.OutputState }

func (GetVpcsResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*GetVpcsResult)(nil)).Elem()
}

func (o GetVpcsResultOutput) ToGetVpcsResultOutput() GetVpcsResultOutput {
	return o
}

func (o GetVpcsResultOutput) ToGetVpcsResultOutputWithContext(ctx context.Context) GetVpcsResultOutput {
	return o
}

func (o GetVpcsResultOutput) Filters() GetVpcsFilterArrayOutput {
	return o.ApplyT(func(v GetVpcsResult) []GetVpcsFilter { return v.Filters }).(GetVpcsFilterArrayOutput)
}

// The provider-assigned unique ID for this managed resource.
func (o GetVpcsResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v GetVpcsResult) string { return v.Id }).(pulumi.StringOutput)
}

// List of all the VPC Ids found.
func (o GetVpcsResultOutput) Ids() pulumi.StringArrayOutput {
	return o.ApplyT(func(v GetVpcsResult) []string { return v.Ids }).(pulumi.StringArrayOutput)
}

func (o GetVpcsResultOutput) Tags() pulumi.StringMapOutput {
	return o.ApplyT(func(v GetVpcsResult) map[string]string { return v.Tags }).(pulumi.StringMapOutput)
}

func init() {
	pulumi.RegisterOutputType(GetVpcsResultOutput{})
}