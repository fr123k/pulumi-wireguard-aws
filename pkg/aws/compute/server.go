package compute

import (
    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"

    "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateSecurityGroups(ctx *pulumi.Context, computeArgs *model.ComputeArgs) ([]*model.SecurityGroup, error) {
    for _, securityGroup := range computeArgs.SecurityGroups {
        securityGroupArgs := &ec2.SecurityGroupArgs{
            Description: pulumi.String(securityGroup.Description),
            Ingress:     network.IngressRules(computeArgs.Security, securityGroup.IngressRules),
            Egress:      network.EgressRules(computeArgs.Security, securityGroup.EgressRules),
            Tags:        pulumi.ToStringMap(securityGroup.Tags),
        }
        if computeArgs.Vpc != nil {
            securityGroupArgs.VpcId = computeArgs.Vpc.ID()
        }
        sgExternal, err := ec2.NewSecurityGroup(ctx, securityGroup.Name, securityGroupArgs)
        if err != nil {
            return nil, err
        }

        securityGroup.State = sgExternal.CustomResourceState
    }
    return computeArgs.SecurityGroups, nil
}
