// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package hcloud

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Provides a Hetzner Cloud Network to represent a Network in the Hetzner Cloud.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
// 	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		_, err := hcloud.NewNetwork(ctx, "privNet", &hcloud.NetworkArgs{
// 			IpRange: pulumi.String("10.0.1.0/24"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type Network struct {
	pulumi.CustomResourceState

	// IP Range of the whole Network which must span all included subnets and route destinations. Must be one of the private ipv4 ranges of RFC1918.
	IpRange pulumi.StringOutput `pulumi:"ipRange"`
	// User-defined labels (key-value pairs) should be created with.
	Labels pulumi.MapOutput `pulumi:"labels"`
	// Name of the Network to create (must be unique per project).
	Name pulumi.StringOutput `pulumi:"name"`
}

// NewNetwork registers a new resource with the given unique name, arguments, and options.
func NewNetwork(ctx *pulumi.Context,
	name string, args *NetworkArgs, opts ...pulumi.ResourceOption) (*Network, error) {
	if args == nil || args.IpRange == nil {
		return nil, errors.New("missing required argument 'IpRange'")
	}
	if args == nil {
		args = &NetworkArgs{}
	}
	var resource Network
	err := ctx.RegisterResource("hcloud:index/network:Network", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetNetwork gets an existing Network resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetNetwork(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *NetworkState, opts ...pulumi.ResourceOption) (*Network, error) {
	var resource Network
	err := ctx.ReadResource("hcloud:index/network:Network", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Network resources.
type networkState struct {
	// IP Range of the whole Network which must span all included subnets and route destinations. Must be one of the private ipv4 ranges of RFC1918.
	IpRange *string `pulumi:"ipRange"`
	// User-defined labels (key-value pairs) should be created with.
	Labels map[string]interface{} `pulumi:"labels"`
	// Name of the Network to create (must be unique per project).
	Name *string `pulumi:"name"`
}

type NetworkState struct {
	// IP Range of the whole Network which must span all included subnets and route destinations. Must be one of the private ipv4 ranges of RFC1918.
	IpRange pulumi.StringPtrInput
	// User-defined labels (key-value pairs) should be created with.
	Labels pulumi.MapInput
	// Name of the Network to create (must be unique per project).
	Name pulumi.StringPtrInput
}

func (NetworkState) ElementType() reflect.Type {
	return reflect.TypeOf((*networkState)(nil)).Elem()
}

type networkArgs struct {
	// IP Range of the whole Network which must span all included subnets and route destinations. Must be one of the private ipv4 ranges of RFC1918.
	IpRange string `pulumi:"ipRange"`
	// User-defined labels (key-value pairs) should be created with.
	Labels map[string]interface{} `pulumi:"labels"`
	// Name of the Network to create (must be unique per project).
	Name *string `pulumi:"name"`
}

// The set of arguments for constructing a Network resource.
type NetworkArgs struct {
	// IP Range of the whole Network which must span all included subnets and route destinations. Must be one of the private ipv4 ranges of RFC1918.
	IpRange pulumi.StringInput
	// User-defined labels (key-value pairs) should be created with.
	Labels pulumi.MapInput
	// Name of the Network to create (must be unique per project).
	Name pulumi.StringPtrInput
}

func (NetworkArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*networkArgs)(nil)).Elem()
}