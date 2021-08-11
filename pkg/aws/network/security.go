package network

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// SSHIngressRule return a ingress security group rule for the ssh port that based on the passed SecurityArgs
//                restrict the access to the VPN cidr or is publicly open.
func SSHIngressRule(security *model.SecurityArgs) *ec2.SecurityGroupIngressArgs {
	if security.VPNEnabledSSH {
		return &ec2.SecurityGroupIngressArgs{
			Protocol:   pulumi.String("tcp"),
			FromPort:   pulumi.Int(22),
			ToPort:     pulumi.Int(22),
			CidrBlocks: pulumi.StringArray{pulumi.String(security.VPNCidr)},
		}
	}
	return &ec2.SecurityGroupIngressArgs{
		Protocol:   pulumi.String("tcp"),
		FromPort:   pulumi.Int(22),
		ToPort:     pulumi.Int(22),
		CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
	}
}

func MapIngress(security *model.SecurityArgs, vs []*model.SecurityRule, f func(*model.SecurityRule) *ec2.SecurityGroupIngressArgs) ec2.SecurityGroupIngressArray{
    var vsm ec2.SecurityGroupIngressArray
	vsm = append(vsm, SSHIngressRule(security))
    for _, v := range vs {
        vsm = append(vsm, f(v))
    }
    return vsm
}


func IngressRules(security *model.SecurityArgs, securityRules []*model.SecurityRule) ec2.SecurityGroupIngressArray {
	transform := func(securityRule *model.SecurityRule) *ec2.SecurityGroupIngressArgs {
		return &ec2.SecurityGroupIngressArgs{
			Protocol:   pulumi.String(securityRule.Protocol),
			FromPort:   pulumi.Int(securityRule.SourcePort),
			ToPort:     pulumi.Int(securityRule.DestinationPort),
			CidrBlocks: pulumi.ToStringArray(securityRule.CidrBlocks),
			SecurityGroups: ToStringArray(securityRule.SecurityGroups),
		}
	};
	return MapIngress(security, securityRules, transform);
}

func ToStringArray(securityGroups []*model.SecurityGroup) pulumi.StringArray {
	a := make(pulumi.StringArray, len(securityGroups))
	for i, securityGroup := range securityGroups {
		a[i] = securityGroup.ID()
	}
	return a
}

func MapEgress(security *model.SecurityArgs, vs []*model.SecurityRule, f func(*model.SecurityRule) *ec2.SecurityGroupEgressArgs) ec2.SecurityGroupEgressArray{
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
	};
	return MapEgress(security, securityRules, transform);
}
