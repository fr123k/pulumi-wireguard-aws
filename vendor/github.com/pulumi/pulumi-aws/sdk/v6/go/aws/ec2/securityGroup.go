// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Provides a security group resource.
//
// > **NOTE:** Avoid using the `ingress` and `egress` arguments of the `ec2.SecurityGroup` resource to configure in-line rules, as they struggle with managing multiple CIDR blocks, and, due to the historical lack of unique IDs, tags and descriptions. To avoid these problems, use the current best practice of the `vpc.SecurityGroupEgressRule` and `vpc.SecurityGroupIngressRule` resources with one CIDR block per rule.
//
// !> **WARNING:** You should not use the `ec2.SecurityGroup` resource with _in-line rules_ (using the `ingress` and `egress` arguments of `ec2.SecurityGroup`) in conjunction with the `vpc.SecurityGroupEgressRule` and `vpc.SecurityGroupIngressRule` resources or the `ec2.SecurityGroupRule` resource. Doing so may cause rule conflicts, perpetual differences, and result in rules being overwritten.
//
// > **NOTE:** Referencing Security Groups across VPC peering has certain restrictions. More information is available in the [VPC Peering User Guide](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-security-groups.html).
//
// > **NOTE:** Due to [AWS Lambda improved VPC networking changes that began deploying in September 2019](https://aws.amazon.com/blogs/compute/announcing-improved-vpc-networking-for-aws-lambda-functions/), security groups associated with Lambda Functions can take up to 45 minutes to successfully delete. To allow for successful deletion, the provider will wait for at least 45 minutes even if a shorter delete timeout is specified.
//
// > **NOTE:** The `cidrBlocks` and `ipv6CidrBlocks` parameters are optional in the `ingress` and `egress` blocks. If nothing is specified, traffic will be blocked as described in _NOTE on Egress rules_ later.
//
// ## Example Usage
//
// ### Basic Usage
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
//	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/vpc"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			allowTls, err := ec2.NewSecurityGroup(ctx, "allow_tls", &ec2.SecurityGroupArgs{
//				Name:        pulumi.String("allow_tls"),
//				Description: pulumi.String("Allow TLS inbound traffic and all outbound traffic"),
//				VpcId:       pulumi.Any(main.Id),
//				Tags: pulumi.StringMap{
//					"Name": pulumi.String("allow_tls"),
//				},
//			})
//			if err != nil {
//				return err
//			}
//			_, err = vpc.NewSecurityGroupIngressRule(ctx, "allow_tls_ipv4", &vpc.SecurityGroupIngressRuleArgs{
//				SecurityGroupId: allowTls.ID(),
//				CidrIpv4:        pulumi.Any(main.CidrBlock),
//				FromPort:        pulumi.Int(443),
//				IpProtocol:      pulumi.String("tcp"),
//				ToPort:          pulumi.Int(443),
//			})
//			if err != nil {
//				return err
//			}
//			_, err = vpc.NewSecurityGroupIngressRule(ctx, "allow_tls_ipv6", &vpc.SecurityGroupIngressRuleArgs{
//				SecurityGroupId: allowTls.ID(),
//				CidrIpv6:        pulumi.Any(main.Ipv6CidrBlock),
//				FromPort:        pulumi.Int(443),
//				IpProtocol:      pulumi.String("tcp"),
//				ToPort:          pulumi.Int(443),
//			})
//			if err != nil {
//				return err
//			}
//			_, err = vpc.NewSecurityGroupEgressRule(ctx, "allow_all_traffic_ipv4", &vpc.SecurityGroupEgressRuleArgs{
//				SecurityGroupId: allowTls.ID(),
//				CidrIpv4:        pulumi.String("0.0.0.0/0"),
//				IpProtocol:      pulumi.String("-1"),
//			})
//			if err != nil {
//				return err
//			}
//			_, err = vpc.NewSecurityGroupEgressRule(ctx, "allow_all_traffic_ipv6", &vpc.SecurityGroupEgressRuleArgs{
//				SecurityGroupId: allowTls.ID(),
//				CidrIpv6:        pulumi.String("::/0"),
//				IpProtocol:      pulumi.String("-1"),
//			})
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// > **NOTE on Egress rules:** By default, AWS creates an `ALLOW ALL` egress rule when creating a new Security Group inside of a VPC. When creating a new Security Group inside a VPC, **this provider will remove this default rule**, and require you specifically re-create it if you desire that rule. We feel this leads to fewer surprises in terms of controlling your egress rules. If you desire this rule to be in place, you can use this `egress` block:
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
//			_, err := ec2.NewSecurityGroup(ctx, "example", &ec2.SecurityGroupArgs{
//				Egress: ec2.SecurityGroupEgressArray{
//					&ec2.SecurityGroupEgressArgs{
//						FromPort: pulumi.Int(0),
//						ToPort:   pulumi.Int(0),
//						Protocol: pulumi.String("-1"),
//						CidrBlocks: pulumi.StringArray{
//							pulumi.String("0.0.0.0/0"),
//						},
//						Ipv6CidrBlocks: pulumi.StringArray{
//							pulumi.String("::/0"),
//						},
//					},
//				},
//			})
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// ### Usage With Prefix List IDs
//
// Prefix Lists are either managed by AWS internally, or created by the customer using a
// Prefix List resource. Prefix Lists provided by
// AWS are associated with a prefix list name, or service name, that is linked to a specific region.
// Prefix list IDs are exported on VPC Endpoints, so you can use this format:
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
//			myEndpoint, err := ec2.NewVpcEndpoint(ctx, "my_endpoint", nil)
//			if err != nil {
//				return err
//			}
//			_, err = ec2.NewSecurityGroup(ctx, "example", &ec2.SecurityGroupArgs{
//				Egress: ec2.SecurityGroupEgressArray{
//					&ec2.SecurityGroupEgressArgs{
//						FromPort: pulumi.Int(0),
//						ToPort:   pulumi.Int(0),
//						Protocol: pulumi.String("-1"),
//						PrefixListIds: pulumi.StringArray{
//							myEndpoint.PrefixListId,
//						},
//					},
//				},
//			})
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// You can also find a specific Prefix List using the `ec2.getPrefixList` data source.
//
// ### Removing All Ingress and Egress Rules
//
// The `ingress` and `egress` arguments are processed in attributes-as-blocks mode. Due to this, removing these arguments from the configuration will **not** cause the provider to destroy the managed rules. To subsequently remove all managed ingress and egress rules:
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
//			_, err := ec2.NewSecurityGroup(ctx, "example", &ec2.SecurityGroupArgs{
//				Name:    pulumi.String("sg"),
//				VpcId:   pulumi.Any(exampleAwsVpc.Id),
//				Ingress: ec2.SecurityGroupIngressArray{},
//				Egress:  ec2.SecurityGroupEgressArray{},
//			})
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// ### Recreating a Security Group
//
// A simple security group `name` change "forces new" the security group--the provider destroys the security group and creates a new one. (Likewise, `description`, `namePrefix`, or `vpcId` [cannot be changed](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/working-with-security-groups.html#creating-security-group).) Attempting to recreate the security group leads to a variety of complications depending on how it is used.
//
// Security groups are generally associated with other resources--**more than 100** AWS Provider resources reference security groups. Referencing a resource from another resource creates a one-way dependency. For example, if you create an EC2 `ec2.Instance` that has a `vpcSecurityGroupIds` argument that refers to an `ec2.SecurityGroup` resource, the `ec2.SecurityGroup` is a dependent of the `ec2.Instance`. Because of this, the provider will create the security group first so that it can then be associated with the EC2 instance.
//
// However, the dependency relationship actually goes both directions causing the _Security Group Deletion Problem_. AWS does not allow you to delete the security group associated with another resource (_e.g._, the `ec2.Instance`).
//
// The provider does not model bi-directional dependencies like this, but, even if it did, simply knowing the dependency situation would not be enough to solve it. For example, some resources must always have an associated security group while others don't need to. In addition, when the `ec2.SecurityGroup` resource attempts to recreate, it receives a dependent object error, which does not provide information on whether the dependent object is a security group rule or, for example, an associated EC2 instance. Within the provider, the associated resource (_e.g._, `ec2.Instance`) does not receive an error when the `ec2.SecurityGroup` is trying to recreate even though that is where changes to the associated resource would need to take place (_e.g._, removing the security group association).
//
// Despite these sticky problems, below are some ways to improve your experience when you find it necessary to recreate a security group.
//
// ### Shorter timeout
//
// (This example is one approach to recreating security groups. For more information on the challenges and the _Security Group Deletion Problem_, see the section above.)
//
// If destroying a security group takes a long time, it may be because the provider cannot distinguish between a dependent object (_e.g._, a security group rule or EC2 instance) that is _in the process of being deleted_ and one that is not. In other words, it may be waiting for a train that isn't scheduled to arrive. To fail faster, shorten the `delete` timeout from the default timeout:
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
//			_, err := ec2.NewSecurityGroup(ctx, "example", &ec2.SecurityGroupArgs{
//				Name: pulumi.String("izizavle"),
//			})
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// ### Provisioners
//
// (This example is one approach to recreating security groups. For more information on the challenges and the _Security Group Deletion Problem_, see the section above.)
//
// **DISCLAIMER:** We **_HIGHLY_** recommend using one of the above approaches and _NOT_ using local provisioners. Provisioners, like the one shown below, should be considered a **last resort** since they are _not readable_, _require skills outside standard configuration_, are _error prone_ and _difficult to maintain_, are not compatible with cloud environments and upgrade tools, require AWS CLI installation, and are subject to changes outside the AWS Provider.
//
// ```go
// package main
//
// import (
//
//	"fmt"
//
//	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
//	"github.com/pulumi/pulumi-command/sdk/go/command/local"
//	"github.com/pulumi/pulumi-null/sdk/go/null"
//	"github.com/pulumi/pulumi-std/sdk/go/std"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			_default, err := ec2.LookupSecurityGroup(ctx, &ec2.LookupSecurityGroupArgs{
//				Name: pulumi.StringRef("default"),
//			}, nil)
//			if err != nil {
//				return err
//			}
//			example, err := ec2.NewSecurityGroup(ctx, "example", &ec2.SecurityGroupArgs{
//				Name: pulumi.String("sg"),
//				Tags: pulumi.StringMap{
//					"workaround1": pulumi.String("tagged-name"),
//					"workaround2": pulumi.String(_default.Id),
//				},
//			})
//			if err != nil {
//				return err
//			}
//			_, err = local.NewCommand(ctx, "exampleProvisioner0", &local.CommandArgs{
//				Create: "true",
//				Update: "true",
//				Delete: fmt.Sprintf("            ENDPOINT_ID=`aws ec2 describe-vpc-endpoints --filters \"Name=tag:Name,Values=%v\" --query \"VpcEndpoints[0].VpcEndpointId\" --output text` &&\n            aws ec2 modify-vpc-endpoint --vpc-endpoint-id ${ENDPOINT_ID} --add-security-group-ids %v --remove-security-group-ids %v\n", tags.Workaround1, tags.Workaround2, id),
//			}, pulumi.DependsOn([]pulumi.Resource{
//				example,
//			}))
//			if err != nil {
//				return err
//			}
//			invokeJoin, err := std.Join(ctx, &std.JoinArgs{
//				Separator: ",",
//				Input:     exampleAwsVpcEndpoint.SecurityGroupIds,
//			}, nil)
//			if err != nil {
//				return err
//			}
//			exampleResource, err := null.NewResource(ctx, "example", &null.ResourceArgs{
//				Triggers: pulumi.StringMap{
//					"rerun_upon_change_of": pulumi.String(invokeJoin.Result),
//				},
//			})
//			if err != nil {
//				return err
//			}
//			_, err = local.NewCommand(ctx, "exampleResourceProvisioner0", &local.CommandArgs{
//				Create: fmt.Sprintf("            aws ec2 modify-vpc-endpoint --vpc-endpoint-id %v --remove-security-group-ids %v\n", exampleAwsVpcEndpoint.Id, _default.Id),
//			}, pulumi.DependsOn([]pulumi.Resource{
//				exampleResource,
//			}))
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
//
// ## Import
//
// Using `pulumi import`, import Security Groups using the security group `id`. For example:
//
// ```sh
// $ pulumi import aws:ec2/securityGroup:SecurityGroup elb_sg sg-903004f8
// ```
type SecurityGroup struct {
	pulumi.CustomResourceState

	// ARN of the security group.
	Arn pulumi.StringOutput `pulumi:"arn"`
	// Security group description. Defaults to `Managed by Pulumi`. Cannot be `""`. **NOTE**: This field maps to the AWS `GroupDescription` attribute, for which there is no Update API. If you'd like to classify your security groups in a way that can be updated, use `tags`.
	Description pulumi.StringOutput `pulumi:"description"`
	// Configuration block for egress rules. Can be specified multiple times for each egress rule. Each egress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Egress SecurityGroupEgressArrayOutput `pulumi:"egress"`
	// Configuration block for ingress rules. Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Ingress SecurityGroupIngressArrayOutput `pulumi:"ingress"`
	// Name of the security group. If omitted, the provider will assign a random, unique name.
	Name pulumi.StringOutput `pulumi:"name"`
	// Creates a unique name beginning with the specified prefix. Conflicts with `name`.
	NamePrefix pulumi.StringOutput `pulumi:"namePrefix"`
	// Owner ID.
	OwnerId pulumi.StringOutput `pulumi:"ownerId"`
	// Instruct the provider to revoke all of the Security Groups attached ingress and egress rules before deleting the rule itself. This is normally not needed, however certain AWS services such as Elastic Map Reduce may automatically add required rules to security groups used with the service, and those rules may contain a cyclic dependency that prevent the security groups from being destroyed without removing the dependency first. Default `false`.
	RevokeRulesOnDelete pulumi.BoolPtrOutput `pulumi:"revokeRulesOnDelete"`
	// Map of tags to assign to the resource. If configured with a provider `defaultTags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
	Tags pulumi.StringMapOutput `pulumi:"tags"`
	// A map of tags assigned to the resource, including those inherited from the provider `defaultTags` configuration block.
	//
	// Deprecated: Please use `tags` instead.
	TagsAll pulumi.StringMapOutput `pulumi:"tagsAll"`
	// VPC ID. Defaults to the region's default VPC.
	VpcId pulumi.StringOutput `pulumi:"vpcId"`
}

// NewSecurityGroup registers a new resource with the given unique name, arguments, and options.
func NewSecurityGroup(ctx *pulumi.Context,
	name string, args *SecurityGroupArgs, opts ...pulumi.ResourceOption) (*SecurityGroup, error) {
	if args == nil {
		args = &SecurityGroupArgs{}
	}

	if args.Description == nil {
		args.Description = pulumi.StringPtr("Managed by Pulumi")
	}
	opts = internal.PkgResourceDefaultOpts(opts)
	var resource SecurityGroup
	err := ctx.RegisterResource("aws:ec2/securityGroup:SecurityGroup", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetSecurityGroup gets an existing SecurityGroup resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetSecurityGroup(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *SecurityGroupState, opts ...pulumi.ResourceOption) (*SecurityGroup, error) {
	var resource SecurityGroup
	err := ctx.ReadResource("aws:ec2/securityGroup:SecurityGroup", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering SecurityGroup resources.
type securityGroupState struct {
	// ARN of the security group.
	Arn *string `pulumi:"arn"`
	// Security group description. Defaults to `Managed by Pulumi`. Cannot be `""`. **NOTE**: This field maps to the AWS `GroupDescription` attribute, for which there is no Update API. If you'd like to classify your security groups in a way that can be updated, use `tags`.
	Description *string `pulumi:"description"`
	// Configuration block for egress rules. Can be specified multiple times for each egress rule. Each egress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Egress []SecurityGroupEgress `pulumi:"egress"`
	// Configuration block for ingress rules. Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Ingress []SecurityGroupIngress `pulumi:"ingress"`
	// Name of the security group. If omitted, the provider will assign a random, unique name.
	Name *string `pulumi:"name"`
	// Creates a unique name beginning with the specified prefix. Conflicts with `name`.
	NamePrefix *string `pulumi:"namePrefix"`
	// Owner ID.
	OwnerId *string `pulumi:"ownerId"`
	// Instruct the provider to revoke all of the Security Groups attached ingress and egress rules before deleting the rule itself. This is normally not needed, however certain AWS services such as Elastic Map Reduce may automatically add required rules to security groups used with the service, and those rules may contain a cyclic dependency that prevent the security groups from being destroyed without removing the dependency first. Default `false`.
	RevokeRulesOnDelete *bool `pulumi:"revokeRulesOnDelete"`
	// Map of tags to assign to the resource. If configured with a provider `defaultTags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
	Tags map[string]string `pulumi:"tags"`
	// A map of tags assigned to the resource, including those inherited from the provider `defaultTags` configuration block.
	//
	// Deprecated: Please use `tags` instead.
	TagsAll map[string]string `pulumi:"tagsAll"`
	// VPC ID. Defaults to the region's default VPC.
	VpcId *string `pulumi:"vpcId"`
}

type SecurityGroupState struct {
	// ARN of the security group.
	Arn pulumi.StringPtrInput
	// Security group description. Defaults to `Managed by Pulumi`. Cannot be `""`. **NOTE**: This field maps to the AWS `GroupDescription` attribute, for which there is no Update API. If you'd like to classify your security groups in a way that can be updated, use `tags`.
	Description pulumi.StringPtrInput
	// Configuration block for egress rules. Can be specified multiple times for each egress rule. Each egress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Egress SecurityGroupEgressArrayInput
	// Configuration block for ingress rules. Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Ingress SecurityGroupIngressArrayInput
	// Name of the security group. If omitted, the provider will assign a random, unique name.
	Name pulumi.StringPtrInput
	// Creates a unique name beginning with the specified prefix. Conflicts with `name`.
	NamePrefix pulumi.StringPtrInput
	// Owner ID.
	OwnerId pulumi.StringPtrInput
	// Instruct the provider to revoke all of the Security Groups attached ingress and egress rules before deleting the rule itself. This is normally not needed, however certain AWS services such as Elastic Map Reduce may automatically add required rules to security groups used with the service, and those rules may contain a cyclic dependency that prevent the security groups from being destroyed without removing the dependency first. Default `false`.
	RevokeRulesOnDelete pulumi.BoolPtrInput
	// Map of tags to assign to the resource. If configured with a provider `defaultTags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
	Tags pulumi.StringMapInput
	// A map of tags assigned to the resource, including those inherited from the provider `defaultTags` configuration block.
	//
	// Deprecated: Please use `tags` instead.
	TagsAll pulumi.StringMapInput
	// VPC ID. Defaults to the region's default VPC.
	VpcId pulumi.StringPtrInput
}

func (SecurityGroupState) ElementType() reflect.Type {
	return reflect.TypeOf((*securityGroupState)(nil)).Elem()
}

type securityGroupArgs struct {
	// Security group description. Defaults to `Managed by Pulumi`. Cannot be `""`. **NOTE**: This field maps to the AWS `GroupDescription` attribute, for which there is no Update API. If you'd like to classify your security groups in a way that can be updated, use `tags`.
	Description *string `pulumi:"description"`
	// Configuration block for egress rules. Can be specified multiple times for each egress rule. Each egress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Egress []SecurityGroupEgress `pulumi:"egress"`
	// Configuration block for ingress rules. Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Ingress []SecurityGroupIngress `pulumi:"ingress"`
	// Name of the security group. If omitted, the provider will assign a random, unique name.
	Name *string `pulumi:"name"`
	// Creates a unique name beginning with the specified prefix. Conflicts with `name`.
	NamePrefix *string `pulumi:"namePrefix"`
	// Instruct the provider to revoke all of the Security Groups attached ingress and egress rules before deleting the rule itself. This is normally not needed, however certain AWS services such as Elastic Map Reduce may automatically add required rules to security groups used with the service, and those rules may contain a cyclic dependency that prevent the security groups from being destroyed without removing the dependency first. Default `false`.
	RevokeRulesOnDelete *bool `pulumi:"revokeRulesOnDelete"`
	// Map of tags to assign to the resource. If configured with a provider `defaultTags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
	Tags map[string]string `pulumi:"tags"`
	// VPC ID. Defaults to the region's default VPC.
	VpcId *string `pulumi:"vpcId"`
}

// The set of arguments for constructing a SecurityGroup resource.
type SecurityGroupArgs struct {
	// Security group description. Defaults to `Managed by Pulumi`. Cannot be `""`. **NOTE**: This field maps to the AWS `GroupDescription` attribute, for which there is no Update API. If you'd like to classify your security groups in a way that can be updated, use `tags`.
	Description pulumi.StringPtrInput
	// Configuration block for egress rules. Can be specified multiple times for each egress rule. Each egress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Egress SecurityGroupEgressArrayInput
	// Configuration block for ingress rules. Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
	Ingress SecurityGroupIngressArrayInput
	// Name of the security group. If omitted, the provider will assign a random, unique name.
	Name pulumi.StringPtrInput
	// Creates a unique name beginning with the specified prefix. Conflicts with `name`.
	NamePrefix pulumi.StringPtrInput
	// Instruct the provider to revoke all of the Security Groups attached ingress and egress rules before deleting the rule itself. This is normally not needed, however certain AWS services such as Elastic Map Reduce may automatically add required rules to security groups used with the service, and those rules may contain a cyclic dependency that prevent the security groups from being destroyed without removing the dependency first. Default `false`.
	RevokeRulesOnDelete pulumi.BoolPtrInput
	// Map of tags to assign to the resource. If configured with a provider `defaultTags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
	Tags pulumi.StringMapInput
	// VPC ID. Defaults to the region's default VPC.
	VpcId pulumi.StringPtrInput
}

func (SecurityGroupArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*securityGroupArgs)(nil)).Elem()
}

type SecurityGroupInput interface {
	pulumi.Input

	ToSecurityGroupOutput() SecurityGroupOutput
	ToSecurityGroupOutputWithContext(ctx context.Context) SecurityGroupOutput
}

func (*SecurityGroup) ElementType() reflect.Type {
	return reflect.TypeOf((**SecurityGroup)(nil)).Elem()
}

func (i *SecurityGroup) ToSecurityGroupOutput() SecurityGroupOutput {
	return i.ToSecurityGroupOutputWithContext(context.Background())
}

func (i *SecurityGroup) ToSecurityGroupOutputWithContext(ctx context.Context) SecurityGroupOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SecurityGroupOutput)
}

// SecurityGroupArrayInput is an input type that accepts SecurityGroupArray and SecurityGroupArrayOutput values.
// You can construct a concrete instance of `SecurityGroupArrayInput` via:
//
//	SecurityGroupArray{ SecurityGroupArgs{...} }
type SecurityGroupArrayInput interface {
	pulumi.Input

	ToSecurityGroupArrayOutput() SecurityGroupArrayOutput
	ToSecurityGroupArrayOutputWithContext(context.Context) SecurityGroupArrayOutput
}

type SecurityGroupArray []SecurityGroupInput

func (SecurityGroupArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*SecurityGroup)(nil)).Elem()
}

func (i SecurityGroupArray) ToSecurityGroupArrayOutput() SecurityGroupArrayOutput {
	return i.ToSecurityGroupArrayOutputWithContext(context.Background())
}

func (i SecurityGroupArray) ToSecurityGroupArrayOutputWithContext(ctx context.Context) SecurityGroupArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SecurityGroupArrayOutput)
}

// SecurityGroupMapInput is an input type that accepts SecurityGroupMap and SecurityGroupMapOutput values.
// You can construct a concrete instance of `SecurityGroupMapInput` via:
//
//	SecurityGroupMap{ "key": SecurityGroupArgs{...} }
type SecurityGroupMapInput interface {
	pulumi.Input

	ToSecurityGroupMapOutput() SecurityGroupMapOutput
	ToSecurityGroupMapOutputWithContext(context.Context) SecurityGroupMapOutput
}

type SecurityGroupMap map[string]SecurityGroupInput

func (SecurityGroupMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*SecurityGroup)(nil)).Elem()
}

func (i SecurityGroupMap) ToSecurityGroupMapOutput() SecurityGroupMapOutput {
	return i.ToSecurityGroupMapOutputWithContext(context.Background())
}

func (i SecurityGroupMap) ToSecurityGroupMapOutputWithContext(ctx context.Context) SecurityGroupMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SecurityGroupMapOutput)
}

type SecurityGroupOutput struct{ *pulumi.OutputState }

func (SecurityGroupOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**SecurityGroup)(nil)).Elem()
}

func (o SecurityGroupOutput) ToSecurityGroupOutput() SecurityGroupOutput {
	return o
}

func (o SecurityGroupOutput) ToSecurityGroupOutputWithContext(ctx context.Context) SecurityGroupOutput {
	return o
}

// ARN of the security group.
func (o SecurityGroupOutput) Arn() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringOutput { return v.Arn }).(pulumi.StringOutput)
}

// Security group description. Defaults to `Managed by Pulumi`. Cannot be `""`. **NOTE**: This field maps to the AWS `GroupDescription` attribute, for which there is no Update API. If you'd like to classify your security groups in a way that can be updated, use `tags`.
func (o SecurityGroupOutput) Description() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringOutput { return v.Description }).(pulumi.StringOutput)
}

// Configuration block for egress rules. Can be specified multiple times for each egress rule. Each egress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
func (o SecurityGroupOutput) Egress() SecurityGroupEgressArrayOutput {
	return o.ApplyT(func(v *SecurityGroup) SecurityGroupEgressArrayOutput { return v.Egress }).(SecurityGroupEgressArrayOutput)
}

// Configuration block for ingress rules. Can be specified multiple times for each ingress rule. Each ingress block supports fields documented below. This argument is processed in attribute-as-blocks mode.
func (o SecurityGroupOutput) Ingress() SecurityGroupIngressArrayOutput {
	return o.ApplyT(func(v *SecurityGroup) SecurityGroupIngressArrayOutput { return v.Ingress }).(SecurityGroupIngressArrayOutput)
}

// Name of the security group. If omitted, the provider will assign a random, unique name.
func (o SecurityGroupOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// Creates a unique name beginning with the specified prefix. Conflicts with `name`.
func (o SecurityGroupOutput) NamePrefix() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringOutput { return v.NamePrefix }).(pulumi.StringOutput)
}

// Owner ID.
func (o SecurityGroupOutput) OwnerId() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringOutput { return v.OwnerId }).(pulumi.StringOutput)
}

// Instruct the provider to revoke all of the Security Groups attached ingress and egress rules before deleting the rule itself. This is normally not needed, however certain AWS services such as Elastic Map Reduce may automatically add required rules to security groups used with the service, and those rules may contain a cyclic dependency that prevent the security groups from being destroyed without removing the dependency first. Default `false`.
func (o SecurityGroupOutput) RevokeRulesOnDelete() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.BoolPtrOutput { return v.RevokeRulesOnDelete }).(pulumi.BoolPtrOutput)
}

// Map of tags to assign to the resource. If configured with a provider `defaultTags` configuration block present, tags with matching keys will overwrite those defined at the provider-level.
func (o SecurityGroupOutput) Tags() pulumi.StringMapOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringMapOutput { return v.Tags }).(pulumi.StringMapOutput)
}

// A map of tags assigned to the resource, including those inherited from the provider `defaultTags` configuration block.
//
// Deprecated: Please use `tags` instead.
func (o SecurityGroupOutput) TagsAll() pulumi.StringMapOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringMapOutput { return v.TagsAll }).(pulumi.StringMapOutput)
}

// VPC ID. Defaults to the region's default VPC.
func (o SecurityGroupOutput) VpcId() pulumi.StringOutput {
	return o.ApplyT(func(v *SecurityGroup) pulumi.StringOutput { return v.VpcId }).(pulumi.StringOutput)
}

type SecurityGroupArrayOutput struct{ *pulumi.OutputState }

func (SecurityGroupArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*SecurityGroup)(nil)).Elem()
}

func (o SecurityGroupArrayOutput) ToSecurityGroupArrayOutput() SecurityGroupArrayOutput {
	return o
}

func (o SecurityGroupArrayOutput) ToSecurityGroupArrayOutputWithContext(ctx context.Context) SecurityGroupArrayOutput {
	return o
}

func (o SecurityGroupArrayOutput) Index(i pulumi.IntInput) SecurityGroupOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *SecurityGroup {
		return vs[0].([]*SecurityGroup)[vs[1].(int)]
	}).(SecurityGroupOutput)
}

type SecurityGroupMapOutput struct{ *pulumi.OutputState }

func (SecurityGroupMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*SecurityGroup)(nil)).Elem()
}

func (o SecurityGroupMapOutput) ToSecurityGroupMapOutput() SecurityGroupMapOutput {
	return o
}

func (o SecurityGroupMapOutput) ToSecurityGroupMapOutputWithContext(ctx context.Context) SecurityGroupMapOutput {
	return o
}

func (o SecurityGroupMapOutput) MapIndex(k pulumi.StringInput) SecurityGroupOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *SecurityGroup {
		return vs[0].(map[string]*SecurityGroup)[vs[1].(string)]
	}).(SecurityGroupOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*SecurityGroupInput)(nil)).Elem(), &SecurityGroup{})
	pulumi.RegisterInputType(reflect.TypeOf((*SecurityGroupArrayInput)(nil)).Elem(), SecurityGroupArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*SecurityGroupMapInput)(nil)).Elem(), SecurityGroupMap{})
	pulumi.RegisterOutputType(SecurityGroupOutput{})
	pulumi.RegisterOutputType(SecurityGroupArrayOutput{})
	pulumi.RegisterOutputType(SecurityGroupMapOutput{})
}