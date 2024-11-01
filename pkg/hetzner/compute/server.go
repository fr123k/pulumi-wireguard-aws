package compute

// import (
//     "fmt"

//     "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
//     "github.com/fr123k/pulumi-wireguard-aws/pkg/model"

//     "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
//     "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )

// func CreateSecurityGroups(ctx *pulumi.Context, computeArgs *model.ComputeArgs) ([]*model.SecurityGroup, []*ec2.SecurityGroup, error) {
//     ec2SecGroups := make([]*ec2.SecurityGroup, 0)
//     for _, securityGroup := range computeArgs.SecurityGroups {
//         securityGroupArgs := &ec2.SecurityGroupArgs{
//             Description: pulumi.String(securityGroup.Description),
//             Ingress:     network.IngressRules(computeArgs.Security, securityGroup.IngressRules),
//             Egress:      network.EgressRules(computeArgs.Security, securityGroup.EgressRules),
//             Tags:        pulumi.ToStringMap(securityGroup.Tags),
//         }
//         if computeArgs.Vpc != nil {
//             securityGroupArgs.VpcId = computeArgs.Vpc.ID()
//         }
//         sgExternal, err := ec2.NewSecurityGroup(ctx, securityGroup.Name, securityGroupArgs)
//         if err != nil {
//             return nil, nil, err
//         }
//         ec2SecGroups = append(ec2SecGroups, sgExternal)
//         securityGroup.State = sgExternal.CustomResourceState
//     }
//     return computeArgs.SecurityGroups, ec2SecGroups, nil
// }

// func GetImage(ctx *pulumi.Context, imageArgs []*model.ImageArgs) (*string, error) {
//     for _, image := range imageArgs {
//         // mostRecent := true
//         filters := make([]ec2.GetAmiIdsFilter, 0)
//         filters = append(filters, ec2.GetAmiIdsFilter{
//             Name:   "name",
//             Values: []string{image.Name},
//         })
//         if len(image.States) > 0 {
//             filters = append(filters, ec2.GetAmiIdsFilter{
//                 Name:   "state",
//                 Values: image.States,
//             })
//         }
//         amiIds, err := ec2.GetAmiIds(ctx, &ec2.GetAmiIdsArgs{
//             Filters: filters,
//             Owners: image.Owners,
//             // MostRecent: &mostRecent,
//         })

//         if err != nil {
//             return nil, err
//         }

//         if amiIds.Ids != nil && len(amiIds.Ids) > 0 {
//             return &amiIds.Ids[0], nil
//         }
//     }

//     return nil, fmt.Errorf("no AMI found")
// }
