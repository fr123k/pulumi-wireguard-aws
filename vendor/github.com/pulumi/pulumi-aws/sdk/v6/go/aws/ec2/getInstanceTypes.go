// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Information about EC2 Instance Types.
//
// ## Example Usage
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
//			_, err := ec2.GetInstanceTypes(ctx, &ec2.GetInstanceTypesArgs{
//				Filters: []ec2.GetInstanceTypesFilter{
//					{
//						Name: "auto-recovery-supported",
//						Values: []string{
//							"true",
//						},
//					},
//					{
//						Name: "network-info.encryption-in-transit-supported",
//						Values: []string{
//							"true",
//						},
//					},
//					{
//						Name: "instance-storage-supported",
//						Values: []string{
//							"true",
//						},
//					},
//					{
//						Name: "instance-type",
//						Values: []string{
//							"g5.2xlarge",
//							"g5.4xlarge",
//						},
//					},
//				},
//			}, nil)
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
func GetInstanceTypes(ctx *pulumi.Context, args *GetInstanceTypesArgs, opts ...pulumi.InvokeOption) (*GetInstanceTypesResult, error) {
	opts = internal.PkgInvokeDefaultOpts(opts)
	var rv GetInstanceTypesResult
	err := ctx.Invoke("aws:ec2/getInstanceTypes:getInstanceTypes", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getInstanceTypes.
type GetInstanceTypesArgs struct {
	// One or more configuration blocks containing name-values filters. See the [EC2 API Reference](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstanceTypes.html) for supported filters. Detailed below.
	Filters []GetInstanceTypesFilter `pulumi:"filters"`
}

// A collection of values returned by getInstanceTypes.
type GetInstanceTypesResult struct {
	Filters []GetInstanceTypesFilter `pulumi:"filters"`
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// List of EC2 Instance Types.
	InstanceTypes []string `pulumi:"instanceTypes"`
}

func GetInstanceTypesOutput(ctx *pulumi.Context, args GetInstanceTypesOutputArgs, opts ...pulumi.InvokeOption) GetInstanceTypesResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (GetInstanceTypesResultOutput, error) {
			args := v.(GetInstanceTypesArgs)
			opts = internal.PkgInvokeDefaultOpts(opts)
			var rv GetInstanceTypesResult
			secret, err := ctx.InvokePackageRaw("aws:ec2/getInstanceTypes:getInstanceTypes", args, &rv, "", opts...)
			if err != nil {
				return GetInstanceTypesResultOutput{}, err
			}

			output := pulumi.ToOutput(rv).(GetInstanceTypesResultOutput)
			if secret {
				return pulumi.ToSecret(output).(GetInstanceTypesResultOutput), nil
			}
			return output, nil
		}).(GetInstanceTypesResultOutput)
}

// A collection of arguments for invoking getInstanceTypes.
type GetInstanceTypesOutputArgs struct {
	// One or more configuration blocks containing name-values filters. See the [EC2 API Reference](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstanceTypes.html) for supported filters. Detailed below.
	Filters GetInstanceTypesFilterArrayInput `pulumi:"filters"`
}

func (GetInstanceTypesOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*GetInstanceTypesArgs)(nil)).Elem()
}

// A collection of values returned by getInstanceTypes.
type GetInstanceTypesResultOutput struct{ *pulumi.OutputState }

func (GetInstanceTypesResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*GetInstanceTypesResult)(nil)).Elem()
}

func (o GetInstanceTypesResultOutput) ToGetInstanceTypesResultOutput() GetInstanceTypesResultOutput {
	return o
}

func (o GetInstanceTypesResultOutput) ToGetInstanceTypesResultOutputWithContext(ctx context.Context) GetInstanceTypesResultOutput {
	return o
}

func (o GetInstanceTypesResultOutput) Filters() GetInstanceTypesFilterArrayOutput {
	return o.ApplyT(func(v GetInstanceTypesResult) []GetInstanceTypesFilter { return v.Filters }).(GetInstanceTypesFilterArrayOutput)
}

// The provider-assigned unique ID for this managed resource.
func (o GetInstanceTypesResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v GetInstanceTypesResult) string { return v.Id }).(pulumi.StringOutput)
}

// List of EC2 Instance Types.
func (o GetInstanceTypesResultOutput) InstanceTypes() pulumi.StringArrayOutput {
	return o.ApplyT(func(v GetInstanceTypesResult) []string { return v.InstanceTypes }).(pulumi.StringArrayOutput)
}

func init() {
	pulumi.RegisterOutputType(GetInstanceTypesResultOutput{})
}