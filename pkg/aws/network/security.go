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
