// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package hcloud

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provides details about a specific Hetzner Cloud Datacenter.
// Use this resource to get detailed information about specific datacenter.
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
// 		opt0 := "fsn1-dc8"
// 		_, err := hcloud.GetDatacenter(ctx, &hcloud.GetDatacenterArgs{
// 			Name: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		opt1 := 4
// 		_, err = hcloud.GetDatacenter(ctx, &hcloud.GetDatacenterArgs{
// 			Id: &opt1,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func GetDatacenter(ctx *pulumi.Context, args *GetDatacenterArgs, opts ...pulumi.InvokeOption) (*GetDatacenterResult, error) {
	var rv GetDatacenterResult
	err := ctx.Invoke("hcloud:index/getDatacenter:getDatacenter", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getDatacenter.
type GetDatacenterArgs struct {
	// ID of the datacenter.
	Id *int `pulumi:"id"`
	// Name of the datacenter.
	Name *string `pulumi:"name"`
}

// A collection of values returned by getDatacenter.
type GetDatacenterResult struct {
	// (list) List of available server types.
	AvailableServerTypeIds []int `pulumi:"availableServerTypeIds"`
	// (string) Description of the datacenter.
	Description string `pulumi:"description"`
	// (int) Unique ID of the datacenter.
	Id int `pulumi:"id"`
	// (map) Physical datacenter location.
	Location map[string]interface{} `pulumi:"location"`
	// (string) Name of the datacenter.
	Name string `pulumi:"name"`
	// (list) List of server types supported by the datacenter.
	SupportedServerTypeIds []int `pulumi:"supportedServerTypeIds"`
}
