// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package hcloud

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provides details about a specific Hetzner Cloud Server Type.
// Use this resource to get detailed information about specific Server Type.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-hcloud/sdk/go/hcloud"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		opt0 := "cx11"
// 		_, err := hcloud.GetServerType(ctx, &hcloud.GetServerTypeArgs{
// 			Name: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		opt1 := 1
// 		_, err = hcloud.GetServerType(ctx, &hcloud.GetServerTypeArgs{
// 			Id: &opt1,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func GetServerType(ctx *pulumi.Context, args *GetServerTypeArgs, opts ...pulumi.InvokeOption) (*GetServerTypeResult, error) {
	var rv GetServerTypeResult
	err := ctx.Invoke("hcloud:index/getServerType:getServerType", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getServerType.
type GetServerTypeArgs struct {
	// ID of the server_type.
	Id *int `pulumi:"id"`
	// Name of the server_type.
	Name *string `pulumi:"name"`
}

// A collection of values returned by getServerType.
type GetServerTypeResult struct {
	// (int) Number of cpu cores a Server of this type will have.
	Cores   int    `pulumi:"cores"`
	CpuType string `pulumi:"cpuType"`
	// (string) Description of the server_type.
	Description string `pulumi:"description"`
	// (int) Disk size a Server of this type will have in GB.
	Disk int `pulumi:"disk"`
	// (int) Unique ID of the server_type.
	Id int `pulumi:"id"`
	// (int) Memory a Server of this type will have in GB.
	Memory int `pulumi:"memory"`
	// (string) Name of the server_type.
	Name        string `pulumi:"name"`
	StorageType string `pulumi:"storageType"`
}
