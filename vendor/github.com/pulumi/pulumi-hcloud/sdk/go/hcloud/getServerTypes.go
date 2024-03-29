// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package hcloud

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provides a list of available Hetzner Cloud Server Types.
func GetServerTypes(ctx *pulumi.Context, args *GetServerTypesArgs, opts ...pulumi.InvokeOption) (*GetServerTypesResult, error) {
	var rv GetServerTypesResult
	err := ctx.Invoke("hcloud:index/getServerTypes:getServerTypes", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getServerTypes.
type GetServerTypesArgs struct {
	// (list) List of unique Server Types identifiers.
	ServerTypeIds []string `pulumi:"serverTypeIds"`
}

// A collection of values returned by getServerTypes.
type GetServerTypesResult struct {
	// (list) List of all Server Types descriptions.
	Descriptions []string `pulumi:"descriptions"`
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// (list) List of Server Types names.
	Names []string `pulumi:"names"`
	// (list) List of unique Server Types identifiers.
	ServerTypeIds []string `pulumi:"serverTypeIds"`
}
