// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package iam

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// This data source can be used to fetch information about a specific
// IAM policy.
//
// ## Example Usage
// ### By ARN
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		opt0 := "arn:aws:iam::123456789012:policy/UsersManageOwnCredentials"
// 		_, err := iam.LookupPolicy(ctx, &iam.LookupPolicyArgs{
// 			Arn: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
// ### By Name
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		opt0 := "test_policy"
// 		_, err := iam.LookupPolicy(ctx, &iam.LookupPolicyArgs{
// 			Name: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupPolicy(ctx *pulumi.Context, args *LookupPolicyArgs, opts ...pulumi.InvokeOption) (*LookupPolicyResult, error) {
	var rv LookupPolicyResult
	err := ctx.Invoke("aws:iam/getPolicy:getPolicy", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getPolicy.
type LookupPolicyArgs struct {
	// The ARN of the IAM policy.
	Arn *string `pulumi:"arn"`
	// The name of the IAM policy.
	Name *string `pulumi:"name"`
	// The prefix of the path to the IAM policy. Defaults to a slash (`/`).
	PathPrefix *string `pulumi:"pathPrefix"`
	// Key-value mapping of tags for the IAM Policy.
	Tags map[string]string `pulumi:"tags"`
}

// A collection of values returned by getPolicy.
type LookupPolicyResult struct {
	Arn string `pulumi:"arn"`
	// The description of the policy.
	Description string `pulumi:"description"`
	// The provider-assigned unique ID for this managed resource.
	Id   string `pulumi:"id"`
	Name string `pulumi:"name"`
	// The path to the policy.
	Path       string  `pulumi:"path"`
	PathPrefix *string `pulumi:"pathPrefix"`
	// The policy document of the policy.
	Policy string `pulumi:"policy"`
	// The policy's ID.
	PolicyId string `pulumi:"policyId"`
	// Key-value mapping of tags for the IAM Policy.
	Tags map[string]string `pulumi:"tags"`
}