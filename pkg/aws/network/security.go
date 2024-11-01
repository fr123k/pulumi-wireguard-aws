package network

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func MapIngress(security *model.SecurityArgs, vs []*model.SecurityRule, f func(*model.SecurityRule) *ec2.SecurityGroupIngressArgs) ec2.SecurityGroupIngressArray {
	var vsm ec2.SecurityGroupIngressArray
	for _, v := range vs {
		vsm = append(vsm, f(v))
	}
	return vsm
}

func IngressRules(security *model.SecurityArgs, securityRules []*model.SecurityRule) ec2.SecurityGroupIngressArray {
	transform := func(securityRule *model.SecurityRule) *ec2.SecurityGroupIngressArgs {
		return &ec2.SecurityGroupIngressArgs{
			Protocol:       pulumi.String(securityRule.Protocol),
			FromPort:       pulumi.Int(securityRule.SourcePort),
			ToPort:         pulumi.Int(securityRule.DestinationPort),
			CidrBlocks:     pulumi.ToStringArray(securityRule.CidrBlocks),
			SecurityGroups: ToStringArray(securityRule.SecurityGroups),
		}
	}
	return MapIngress(security, securityRules, transform)
}

func ToStringArray(securityGroups []*model.SecurityGroup) pulumi.StringArray {
	a := make(pulumi.StringArray, len(securityGroups))
	for i, securityGroup := range securityGroups {
		a[i] = securityGroup.ID()
	}
	return a
}

func MapEgress(security *model.SecurityArgs, vs []*model.SecurityRule, f func(*model.SecurityRule) *ec2.SecurityGroupEgressArgs) ec2.SecurityGroupEgressArray {
	var vsm ec2.SecurityGroupEgressArray
	for _, v := range vs {
		vsm = append(vsm, f(v))
	}
	return vsm
}

func EgressRules(security *model.SecurityArgs, securityRules []*model.SecurityRule) ec2.SecurityGroupEgressArray {
	transform := func(securityRule *model.SecurityRule) *ec2.SecurityGroupEgressArgs {
		return &ec2.SecurityGroupEgressArgs{
			Protocol:   pulumi.String(securityRule.Protocol),
			FromPort:   pulumi.Int(securityRule.SourcePort),
			ToPort:     pulumi.Int(securityRule.DestinationPort),
			CidrBlocks: pulumi.StringArray{pulumi.String(securityRule.CidrBlocks[0])},
		}
	}
	return MapEgress(security, securityRules, transform)
}
