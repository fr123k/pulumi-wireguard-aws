package compute

import (
    "fmt"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"

    "github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
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

func GetImage(ctx *pulumi.Context, imageArgs []*model.ImageArgs) (*string, error) {
    for _, image := range imageArgs {
        // mostRecent := true

        amiIds, err := aws.GetAmiIds(ctx, &aws.GetAmiIdsArgs{
            Filters: []aws.GetAmiIdsFilter{
                {
                    Name:   "name",
                    Values: []string{image.Name},
                },
                {
                    Name:   "state",
                    Values: image.States,
                },
            },
            Owners: image.Owners,
            // MostRecent: &mostRecent,
        })

        if err != nil {
            return nil, err
        }

        if amiIds.Ids != nil && len(amiIds.Ids) > 0 {
            return &amiIds.Ids[0], nil
        }
    }

    return nil, fmt.Errorf("no AMI found.")
}
