// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package iam

import (
	"context"
	"reflect"

	"errors"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ## Import
//
// Using `pulumi import`, import exclusive management of customer managed policy assignments using the `user_name`. For example:
//
// ```sh
// $ pulumi import aws:iam/userPolicyAttachmentsExclusive:UserPolicyAttachmentsExclusive example MyUser
// ```
type UserPolicyAttachmentsExclusive struct {
	pulumi.CustomResourceState

	// A list of customer managed policy ARNs to be attached to the user. Policies attached to this user but not configured in this argument will be removed.
	PolicyArns pulumi.StringArrayOutput `pulumi:"policyArns"`
	// IAM user name.
	UserName pulumi.StringOutput `pulumi:"userName"`
}

// NewUserPolicyAttachmentsExclusive registers a new resource with the given unique name, arguments, and options.
func NewUserPolicyAttachmentsExclusive(ctx *pulumi.Context,
	name string, args *UserPolicyAttachmentsExclusiveArgs, opts ...pulumi.ResourceOption) (*UserPolicyAttachmentsExclusive, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.PolicyArns == nil {
		return nil, errors.New("invalid value for required argument 'PolicyArns'")
	}
	if args.UserName == nil {
		return nil, errors.New("invalid value for required argument 'UserName'")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource UserPolicyAttachmentsExclusive
	err := ctx.RegisterResource("aws:iam/userPolicyAttachmentsExclusive:UserPolicyAttachmentsExclusive", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetUserPolicyAttachmentsExclusive gets an existing UserPolicyAttachmentsExclusive resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetUserPolicyAttachmentsExclusive(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *UserPolicyAttachmentsExclusiveState, opts ...pulumi.ResourceOption) (*UserPolicyAttachmentsExclusive, error) {
	var resource UserPolicyAttachmentsExclusive
	err := ctx.ReadResource("aws:iam/userPolicyAttachmentsExclusive:UserPolicyAttachmentsExclusive", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering UserPolicyAttachmentsExclusive resources.
type userPolicyAttachmentsExclusiveState struct {
	// A list of customer managed policy ARNs to be attached to the user. Policies attached to this user but not configured in this argument will be removed.
	PolicyArns []string `pulumi:"policyArns"`
	// IAM user name.
	UserName *string `pulumi:"userName"`
}

type UserPolicyAttachmentsExclusiveState struct {
	// A list of customer managed policy ARNs to be attached to the user. Policies attached to this user but not configured in this argument will be removed.
	PolicyArns pulumi.StringArrayInput
	// IAM user name.
	UserName pulumi.StringPtrInput
}

func (UserPolicyAttachmentsExclusiveState) ElementType() reflect.Type {
	return reflect.TypeOf((*userPolicyAttachmentsExclusiveState)(nil)).Elem()
}

type userPolicyAttachmentsExclusiveArgs struct {
	// A list of customer managed policy ARNs to be attached to the user. Policies attached to this user but not configured in this argument will be removed.
	PolicyArns []string `pulumi:"policyArns"`
	// IAM user name.
	UserName string `pulumi:"userName"`
}

// The set of arguments for constructing a UserPolicyAttachmentsExclusive resource.
type UserPolicyAttachmentsExclusiveArgs struct {
	// A list of customer managed policy ARNs to be attached to the user. Policies attached to this user but not configured in this argument will be removed.
	PolicyArns pulumi.StringArrayInput
	// IAM user name.
	UserName pulumi.StringInput
}

func (UserPolicyAttachmentsExclusiveArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*userPolicyAttachmentsExclusiveArgs)(nil)).Elem()
}

type UserPolicyAttachmentsExclusiveInput interface {
	pulumi.Input

	ToUserPolicyAttachmentsExclusiveOutput() UserPolicyAttachmentsExclusiveOutput
	ToUserPolicyAttachmentsExclusiveOutputWithContext(ctx context.Context) UserPolicyAttachmentsExclusiveOutput
}

func (*UserPolicyAttachmentsExclusive) ElementType() reflect.Type {
	return reflect.TypeOf((**UserPolicyAttachmentsExclusive)(nil)).Elem()
}

func (i *UserPolicyAttachmentsExclusive) ToUserPolicyAttachmentsExclusiveOutput() UserPolicyAttachmentsExclusiveOutput {
	return i.ToUserPolicyAttachmentsExclusiveOutputWithContext(context.Background())
}

func (i *UserPolicyAttachmentsExclusive) ToUserPolicyAttachmentsExclusiveOutputWithContext(ctx context.Context) UserPolicyAttachmentsExclusiveOutput {
	return pulumi.ToOutputWithContext(ctx, i).(UserPolicyAttachmentsExclusiveOutput)
}

// UserPolicyAttachmentsExclusiveArrayInput is an input type that accepts UserPolicyAttachmentsExclusiveArray and UserPolicyAttachmentsExclusiveArrayOutput values.
// You can construct a concrete instance of `UserPolicyAttachmentsExclusiveArrayInput` via:
//
//	UserPolicyAttachmentsExclusiveArray{ UserPolicyAttachmentsExclusiveArgs{...} }
type UserPolicyAttachmentsExclusiveArrayInput interface {
	pulumi.Input

	ToUserPolicyAttachmentsExclusiveArrayOutput() UserPolicyAttachmentsExclusiveArrayOutput
	ToUserPolicyAttachmentsExclusiveArrayOutputWithContext(context.Context) UserPolicyAttachmentsExclusiveArrayOutput
}

type UserPolicyAttachmentsExclusiveArray []UserPolicyAttachmentsExclusiveInput

func (UserPolicyAttachmentsExclusiveArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*UserPolicyAttachmentsExclusive)(nil)).Elem()
}

func (i UserPolicyAttachmentsExclusiveArray) ToUserPolicyAttachmentsExclusiveArrayOutput() UserPolicyAttachmentsExclusiveArrayOutput {
	return i.ToUserPolicyAttachmentsExclusiveArrayOutputWithContext(context.Background())
}

func (i UserPolicyAttachmentsExclusiveArray) ToUserPolicyAttachmentsExclusiveArrayOutputWithContext(ctx context.Context) UserPolicyAttachmentsExclusiveArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(UserPolicyAttachmentsExclusiveArrayOutput)
}

// UserPolicyAttachmentsExclusiveMapInput is an input type that accepts UserPolicyAttachmentsExclusiveMap and UserPolicyAttachmentsExclusiveMapOutput values.
// You can construct a concrete instance of `UserPolicyAttachmentsExclusiveMapInput` via:
//
//	UserPolicyAttachmentsExclusiveMap{ "key": UserPolicyAttachmentsExclusiveArgs{...} }
type UserPolicyAttachmentsExclusiveMapInput interface {
	pulumi.Input

	ToUserPolicyAttachmentsExclusiveMapOutput() UserPolicyAttachmentsExclusiveMapOutput
	ToUserPolicyAttachmentsExclusiveMapOutputWithContext(context.Context) UserPolicyAttachmentsExclusiveMapOutput
}

type UserPolicyAttachmentsExclusiveMap map[string]UserPolicyAttachmentsExclusiveInput

func (UserPolicyAttachmentsExclusiveMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*UserPolicyAttachmentsExclusive)(nil)).Elem()
}

func (i UserPolicyAttachmentsExclusiveMap) ToUserPolicyAttachmentsExclusiveMapOutput() UserPolicyAttachmentsExclusiveMapOutput {
	return i.ToUserPolicyAttachmentsExclusiveMapOutputWithContext(context.Background())
}

func (i UserPolicyAttachmentsExclusiveMap) ToUserPolicyAttachmentsExclusiveMapOutputWithContext(ctx context.Context) UserPolicyAttachmentsExclusiveMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(UserPolicyAttachmentsExclusiveMapOutput)
}

type UserPolicyAttachmentsExclusiveOutput struct{ *pulumi.OutputState }

func (UserPolicyAttachmentsExclusiveOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**UserPolicyAttachmentsExclusive)(nil)).Elem()
}

func (o UserPolicyAttachmentsExclusiveOutput) ToUserPolicyAttachmentsExclusiveOutput() UserPolicyAttachmentsExclusiveOutput {
	return o
}

func (o UserPolicyAttachmentsExclusiveOutput) ToUserPolicyAttachmentsExclusiveOutputWithContext(ctx context.Context) UserPolicyAttachmentsExclusiveOutput {
	return o
}

// A list of customer managed policy ARNs to be attached to the user. Policies attached to this user but not configured in this argument will be removed.
func (o UserPolicyAttachmentsExclusiveOutput) PolicyArns() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *UserPolicyAttachmentsExclusive) pulumi.StringArrayOutput { return v.PolicyArns }).(pulumi.StringArrayOutput)
}

// IAM user name.
func (o UserPolicyAttachmentsExclusiveOutput) UserName() pulumi.StringOutput {
	return o.ApplyT(func(v *UserPolicyAttachmentsExclusive) pulumi.StringOutput { return v.UserName }).(pulumi.StringOutput)
}

type UserPolicyAttachmentsExclusiveArrayOutput struct{ *pulumi.OutputState }

func (UserPolicyAttachmentsExclusiveArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*UserPolicyAttachmentsExclusive)(nil)).Elem()
}

func (o UserPolicyAttachmentsExclusiveArrayOutput) ToUserPolicyAttachmentsExclusiveArrayOutput() UserPolicyAttachmentsExclusiveArrayOutput {
	return o
}

func (o UserPolicyAttachmentsExclusiveArrayOutput) ToUserPolicyAttachmentsExclusiveArrayOutputWithContext(ctx context.Context) UserPolicyAttachmentsExclusiveArrayOutput {
	return o
}

func (o UserPolicyAttachmentsExclusiveArrayOutput) Index(i pulumi.IntInput) UserPolicyAttachmentsExclusiveOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *UserPolicyAttachmentsExclusive {
		return vs[0].([]*UserPolicyAttachmentsExclusive)[vs[1].(int)]
	}).(UserPolicyAttachmentsExclusiveOutput)
}

type UserPolicyAttachmentsExclusiveMapOutput struct{ *pulumi.OutputState }

func (UserPolicyAttachmentsExclusiveMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*UserPolicyAttachmentsExclusive)(nil)).Elem()
}

func (o UserPolicyAttachmentsExclusiveMapOutput) ToUserPolicyAttachmentsExclusiveMapOutput() UserPolicyAttachmentsExclusiveMapOutput {
	return o
}

func (o UserPolicyAttachmentsExclusiveMapOutput) ToUserPolicyAttachmentsExclusiveMapOutputWithContext(ctx context.Context) UserPolicyAttachmentsExclusiveMapOutput {
	return o
}

func (o UserPolicyAttachmentsExclusiveMapOutput) MapIndex(k pulumi.StringInput) UserPolicyAttachmentsExclusiveOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *UserPolicyAttachmentsExclusive {
		return vs[0].(map[string]*UserPolicyAttachmentsExclusive)[vs[1].(string)]
	}).(UserPolicyAttachmentsExclusiveOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*UserPolicyAttachmentsExclusiveInput)(nil)).Elem(), &UserPolicyAttachmentsExclusive{})
	pulumi.RegisterInputType(reflect.TypeOf((*UserPolicyAttachmentsExclusiveArrayInput)(nil)).Elem(), UserPolicyAttachmentsExclusiveArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*UserPolicyAttachmentsExclusiveMapInput)(nil)).Elem(), UserPolicyAttachmentsExclusiveMap{})
	pulumi.RegisterOutputType(UserPolicyAttachmentsExclusiveOutput{})
	pulumi.RegisterOutputType(UserPolicyAttachmentsExclusiveArrayOutput{})
	pulumi.RegisterOutputType(UserPolicyAttachmentsExclusiveMapOutput{})
}