// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package aws

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetDefaultTags(ctx *pulumi.Context, args *GetDefaultTagsArgs, opts ...pulumi.InvokeOption) (*GetDefaultTagsResult, error) {
	var rv GetDefaultTagsResult
	err := ctx.Invoke("aws:index/getDefaultTags:getDefaultTags", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getDefaultTags.
type GetDefaultTagsArgs struct {
	// Blocks of default tags set on the provider. See details below.
	Tags map[string]string `pulumi:"tags"`
}

// A collection of values returned by getDefaultTags.
type GetDefaultTagsResult struct {
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// Blocks of default tags set on the provider. See details below.
	Tags map[string]string `pulumi:"tags"`
}
