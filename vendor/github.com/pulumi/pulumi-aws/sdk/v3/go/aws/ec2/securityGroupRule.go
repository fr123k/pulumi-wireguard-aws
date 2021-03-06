// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

// Provides a security group rule resource. Represents a single `ingress` or
// `egress` group rule, which can be added to external Security Groups.
//
// > **NOTE on Security Groups and Security Group Rules:** This provider currently
// provides both a standalone Security Group Rule resource (a single `ingress` or
// `egress` rule), and a Security Group resource with `ingress` and `egress` rules
// defined in-line. At this time you cannot use a Security Group with in-line rules
// in conjunction with any Security Group Rule resources. Doing so will cause
// a conflict of rule settings and will overwrite rules.
//
// > **NOTE:** Setting `protocol = "all"` or `protocol = -1` with `fromPort` and `toPort` will result in the EC2 API creating a security group rule with all ports open. This API behavior cannot be controlled by this provider and may generate warnings in the future.
//
// > **NOTE:** Referencing Security Groups across VPC peering has certain restrictions. More information is available in the [VPC Peering User Guide](https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-security-groups.html).
//
// ## Usage with prefix list IDs
//
// Prefix list IDs are managed by AWS internally. Prefix list IDs
// are associated with a prefix list name, or service name, that is linked to a specific region.
// Prefix list IDs are exported on VPC Endpoints, so you can use this format:
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
// 	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		myEndpoint, err := ec2.NewVpcEndpoint(ctx, "myEndpoint", nil)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = ec2.NewSecurityGroupRule(ctx, "allowAll", &ec2.SecurityGroupRuleArgs{
// 			Type:     pulumi.String("egress"),
// 			ToPort:   pulumi.Int(0),
// 			Protocol: pulumi.String("-1"),
// 			PrefixListIds: pulumi.StringArray{
// 				myEndpoint.PrefixListId,
// 			},
// 			FromPort:        pulumi.Int(0),
// 			SecurityGroupId: pulumi.String("sg-123456"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
//
// ## Import
//
// ### Examples Import an ingress rule in security group `sg-6e616f6d69` for TCP port 8000 with an IPv4 destination CIDR of `10.0.3.0/24`console
//
// ```sh
//  $ pulumi import aws:ec2/securityGroupRule:SecurityGroupRule ingress sg-6e616f6d69_ingress_tcp_8000_8000_10.0.3.0/24
// ```
//
//  Import a rule with various IPv4 and IPv6 source CIDR blocksconsole
//
// ```sh
//  $ pulumi import aws:ec2/securityGroupRule:SecurityGroupRule ingress sg-4973616163_ingress_tcp_100_121_10.1.0.0/16_2001:db8::/48_10.2.0.0/16_2002:db8::/48
// ```
//
//  Import a rule, applicable to all ports, with a protocol other than TCP/UDP/ICMP/ALL, e.g., Multicast Transport Protocol (MTP), using the IANA protocol number, e.g., 92. console
//
// ```sh
//  $ pulumi import aws:ec2/securityGroupRule:SecurityGroupRule ingress sg-6777656e646f6c796e_ingress_92_0_65536_10.0.3.0/24_10.0.4.0/24
// ```
//
//  Import an egress rule with a prefix list ID destinationconsole
//
// ```sh
//  $ pulumi import aws:ec2/securityGroupRule:SecurityGroupRule egress sg-62726f6479_egress_tcp_8000_8000_pl-6469726b
// ```
//
//  Import a rule applicable to all protocols and ports with a security group sourceconsole
//
// ```sh
//  $ pulumi import aws:ec2/securityGroupRule:SecurityGroupRule ingress_rule sg-7472697374616e_ingress_all_0_65536_sg-6176657279
// ```
//
//  Import a rule that has itself and an IPv6 CIDR block as sourcesconsole
//
// ```sh
//  $ pulumi import aws:ec2/securityGroupRule:SecurityGroupRule rule_name sg-656c65616e6f72_ingress_tcp_80_80_self_2001:db8::/48
// ```
type SecurityGroupRule struct {
	pulumi.CustomResourceState

	// List of CIDR blocks. Cannot be specified with `sourceSecurityGroupId`.
	CidrBlocks pulumi.StringArrayOutput `pulumi:"cidrBlocks"`
	// Description of the rule.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// The start port (or ICMP type number if protocol is "icmp" or "icmpv6").
	FromPort pulumi.IntOutput `pulumi:"fromPort"`
	// List of IPv6 CIDR blocks.
	Ipv6CidrBlocks pulumi.StringArrayOutput `pulumi:"ipv6CidrBlocks"`
	// List of prefix list IDs (for allowing access to VPC endpoints).
	// Only valid with `egress`.
	PrefixListIds pulumi.StringArrayOutput `pulumi:"prefixListIds"`
	// The protocol. If not icmp, icmpv6, tcp, udp, or all use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)
	Protocol pulumi.StringOutput `pulumi:"protocol"`
	// The security group to apply this rule to.
	SecurityGroupId pulumi.StringOutput `pulumi:"securityGroupId"`
	// If true, the security group itself will be added as
	// a source to this ingress rule. Cannot be specified with `sourceSecurityGroupId`.
	Self pulumi.BoolPtrOutput `pulumi:"self"`
	// The security group id to allow access to/from,
	// depending on the `type`. Cannot be specified with `cidrBlocks` and `self`.
	SourceSecurityGroupId pulumi.StringOutput `pulumi:"sourceSecurityGroupId"`
	// The end port (or ICMP code if protocol is "icmp").
	ToPort pulumi.IntOutput `pulumi:"toPort"`
	// The type of rule being created. Valid options are `ingress` (inbound)
	// or `egress` (outbound).
	Type pulumi.StringOutput `pulumi:"type"`
}

// NewSecurityGroupRule registers a new resource with the given unique name, arguments, and options.
func NewSecurityGroupRule(ctx *pulumi.Context,
	name string, args *SecurityGroupRuleArgs, opts ...pulumi.ResourceOption) (*SecurityGroupRule, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.FromPort == nil {
		return nil, errors.New("invalid value for required argument 'FromPort'")
	}
	if args.Protocol == nil {
		return nil, errors.New("invalid value for required argument 'Protocol'")
	}
	if args.SecurityGroupId == nil {
		return nil, errors.New("invalid value for required argument 'SecurityGroupId'")
	}
	if args.ToPort == nil {
		return nil, errors.New("invalid value for required argument 'ToPort'")
	}
	if args.Type == nil {
		return nil, errors.New("invalid value for required argument 'Type'")
	}
	var resource SecurityGroupRule
	err := ctx.RegisterResource("aws:ec2/securityGroupRule:SecurityGroupRule", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetSecurityGroupRule gets an existing SecurityGroupRule resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetSecurityGroupRule(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *SecurityGroupRuleState, opts ...pulumi.ResourceOption) (*SecurityGroupRule, error) {
	var resource SecurityGroupRule
	err := ctx.ReadResource("aws:ec2/securityGroupRule:SecurityGroupRule", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering SecurityGroupRule resources.
type securityGroupRuleState struct {
	// List of CIDR blocks. Cannot be specified with `sourceSecurityGroupId`.
	CidrBlocks []string `pulumi:"cidrBlocks"`
	// Description of the rule.
	Description *string `pulumi:"description"`
	// The start port (or ICMP type number if protocol is "icmp" or "icmpv6").
	FromPort *int `pulumi:"fromPort"`
	// List of IPv6 CIDR blocks.
	Ipv6CidrBlocks []string `pulumi:"ipv6CidrBlocks"`
	// List of prefix list IDs (for allowing access to VPC endpoints).
	// Only valid with `egress`.
	PrefixListIds []string `pulumi:"prefixListIds"`
	// The protocol. If not icmp, icmpv6, tcp, udp, or all use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)
	Protocol *string `pulumi:"protocol"`
	// The security group to apply this rule to.
	SecurityGroupId *string `pulumi:"securityGroupId"`
	// If true, the security group itself will be added as
	// a source to this ingress rule. Cannot be specified with `sourceSecurityGroupId`.
	Self *bool `pulumi:"self"`
	// The security group id to allow access to/from,
	// depending on the `type`. Cannot be specified with `cidrBlocks` and `self`.
	SourceSecurityGroupId *string `pulumi:"sourceSecurityGroupId"`
	// The end port (or ICMP code if protocol is "icmp").
	ToPort *int `pulumi:"toPort"`
	// The type of rule being created. Valid options are `ingress` (inbound)
	// or `egress` (outbound).
	Type *string `pulumi:"type"`
}

type SecurityGroupRuleState struct {
	// List of CIDR blocks. Cannot be specified with `sourceSecurityGroupId`.
	CidrBlocks pulumi.StringArrayInput
	// Description of the rule.
	Description pulumi.StringPtrInput
	// The start port (or ICMP type number if protocol is "icmp" or "icmpv6").
	FromPort pulumi.IntPtrInput
	// List of IPv6 CIDR blocks.
	Ipv6CidrBlocks pulumi.StringArrayInput
	// List of prefix list IDs (for allowing access to VPC endpoints).
	// Only valid with `egress`.
	PrefixListIds pulumi.StringArrayInput
	// The protocol. If not icmp, icmpv6, tcp, udp, or all use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)
	Protocol pulumi.StringPtrInput
	// The security group to apply this rule to.
	SecurityGroupId pulumi.StringPtrInput
	// If true, the security group itself will be added as
	// a source to this ingress rule. Cannot be specified with `sourceSecurityGroupId`.
	Self pulumi.BoolPtrInput
	// The security group id to allow access to/from,
	// depending on the `type`. Cannot be specified with `cidrBlocks` and `self`.
	SourceSecurityGroupId pulumi.StringPtrInput
	// The end port (or ICMP code if protocol is "icmp").
	ToPort pulumi.IntPtrInput
	// The type of rule being created. Valid options are `ingress` (inbound)
	// or `egress` (outbound).
	Type pulumi.StringPtrInput
}

func (SecurityGroupRuleState) ElementType() reflect.Type {
	return reflect.TypeOf((*securityGroupRuleState)(nil)).Elem()
}

type securityGroupRuleArgs struct {
	// List of CIDR blocks. Cannot be specified with `sourceSecurityGroupId`.
	CidrBlocks []string `pulumi:"cidrBlocks"`
	// Description of the rule.
	Description *string `pulumi:"description"`
	// The start port (or ICMP type number if protocol is "icmp" or "icmpv6").
	FromPort int `pulumi:"fromPort"`
	// List of IPv6 CIDR blocks.
	Ipv6CidrBlocks []string `pulumi:"ipv6CidrBlocks"`
	// List of prefix list IDs (for allowing access to VPC endpoints).
	// Only valid with `egress`.
	PrefixListIds []string `pulumi:"prefixListIds"`
	// The protocol. If not icmp, icmpv6, tcp, udp, or all use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)
	Protocol string `pulumi:"protocol"`
	// The security group to apply this rule to.
	SecurityGroupId string `pulumi:"securityGroupId"`
	// If true, the security group itself will be added as
	// a source to this ingress rule. Cannot be specified with `sourceSecurityGroupId`.
	Self *bool `pulumi:"self"`
	// The security group id to allow access to/from,
	// depending on the `type`. Cannot be specified with `cidrBlocks` and `self`.
	SourceSecurityGroupId *string `pulumi:"sourceSecurityGroupId"`
	// The end port (or ICMP code if protocol is "icmp").
	ToPort int `pulumi:"toPort"`
	// The type of rule being created. Valid options are `ingress` (inbound)
	// or `egress` (outbound).
	Type string `pulumi:"type"`
}

// The set of arguments for constructing a SecurityGroupRule resource.
type SecurityGroupRuleArgs struct {
	// List of CIDR blocks. Cannot be specified with `sourceSecurityGroupId`.
	CidrBlocks pulumi.StringArrayInput
	// Description of the rule.
	Description pulumi.StringPtrInput
	// The start port (or ICMP type number if protocol is "icmp" or "icmpv6").
	FromPort pulumi.IntInput
	// List of IPv6 CIDR blocks.
	Ipv6CidrBlocks pulumi.StringArrayInput
	// List of prefix list IDs (for allowing access to VPC endpoints).
	// Only valid with `egress`.
	PrefixListIds pulumi.StringArrayInput
	// The protocol. If not icmp, icmpv6, tcp, udp, or all use the [protocol number](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)
	Protocol pulumi.StringInput
	// The security group to apply this rule to.
	SecurityGroupId pulumi.StringInput
	// If true, the security group itself will be added as
	// a source to this ingress rule. Cannot be specified with `sourceSecurityGroupId`.
	Self pulumi.BoolPtrInput
	// The security group id to allow access to/from,
	// depending on the `type`. Cannot be specified with `cidrBlocks` and `self`.
	SourceSecurityGroupId pulumi.StringPtrInput
	// The end port (or ICMP code if protocol is "icmp").
	ToPort pulumi.IntInput
	// The type of rule being created. Valid options are `ingress` (inbound)
	// or `egress` (outbound).
	Type pulumi.StringInput
}

func (SecurityGroupRuleArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*securityGroupRuleArgs)(nil)).Elem()
}

type SecurityGroupRuleInput interface {
	pulumi.Input

	ToSecurityGroupRuleOutput() SecurityGroupRuleOutput
	ToSecurityGroupRuleOutputWithContext(ctx context.Context) SecurityGroupRuleOutput
}

func (SecurityGroupRule) ElementType() reflect.Type {
	return reflect.TypeOf((*SecurityGroupRule)(nil)).Elem()
}

func (i SecurityGroupRule) ToSecurityGroupRuleOutput() SecurityGroupRuleOutput {
	return i.ToSecurityGroupRuleOutputWithContext(context.Background())
}

func (i SecurityGroupRule) ToSecurityGroupRuleOutputWithContext(ctx context.Context) SecurityGroupRuleOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SecurityGroupRuleOutput)
}

type SecurityGroupRuleOutput struct {
	*pulumi.OutputState
}

func (SecurityGroupRuleOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*SecurityGroupRuleOutput)(nil)).Elem()
}

func (o SecurityGroupRuleOutput) ToSecurityGroupRuleOutput() SecurityGroupRuleOutput {
	return o
}

func (o SecurityGroupRuleOutput) ToSecurityGroupRuleOutputWithContext(ctx context.Context) SecurityGroupRuleOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(SecurityGroupRuleOutput{})
}
