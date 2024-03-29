// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package iam

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetSessionContext(ctx *pulumi.Context, args *GetSessionContextArgs, opts ...pulumi.InvokeOption) (*GetSessionContextResult, error) {
	var rv GetSessionContextResult
	err := ctx.Invoke("aws:iam/getSessionContext:getSessionContext", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getSessionContext.
type GetSessionContextArgs struct {
	// ARN for an assumed role.
	Arn string `pulumi:"arn"`
}

// A collection of values returned by getSessionContext.
type GetSessionContextResult struct {
	Arn string `pulumi:"arn"`
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// IAM source role ARN if `arn` corresponds to an STS assumed role. Otherwise, `issuerArn` is equal to `arn`.
	IssuerArn string `pulumi:"issuerArn"`
	// Unique identifier of the IAM role that issues the STS assumed role.
	IssuerId string `pulumi:"issuerId"`
	// Name of the source role. Only available if `arn` corresponds to an STS assumed role.
	IssuerName string `pulumi:"issuerName"`
	// Name of the STS session. Only available if `arn` corresponds to an STS assumed role.
	SessionName string `pulumi:"sessionName"`
}
