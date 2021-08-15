package compute

import (
	"fmt"
	"os"
	"strings"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const size = "t2.micro"

type Infrastructure struct {
    Groups   []*ec2.SecurityGroup
    Server   *ec2.Instance
    ImageID  *string
    UserData *model.UserData
}

type exportsFnc = func(ctx *pulumi.Context, infra *Infrastructure)

func CreateServer(ctx *pulumi.Context, computeArgs *model.ComputeArgs, exports exportsFnc) (*Infrastructure, error) {
    securityGroups, ec2SecurityGroups, err := CreateSecurityGroups(ctx, computeArgs)
    if err != nil {
        return nil, err
    }

    imageID, err := GetImage(ctx, computeArgs.Images)

    if err != nil {
        return nil, err
    }

    if computeArgs.UserData != nil {
        ctx.Export("cloud-init", pulumi.String(computeArgs.UserData.Content))
    }

    //TODO improve only create keypair if sshkey os specified
    if computeArgs.KeyPair.SSHKeyPair != nil {
        _, err := ec2.NewKeyPair(ctx, *computeArgs.KeyPair.Name, &ec2.KeyPairArgs{
            KeyName:   pulumi.String(*computeArgs.KeyPair.Name),
            PublicKey: pulumi.String(*computeArgs.KeyPair.SSHKeyPair.PublicKeyStr),
        })
        if err != nil {
            return nil, err
        }
    }

    ec2Args := &ec2.InstanceArgs{
        AssociatePublicIpAddress: pulumi.Bool(true),
        //TODO pass tags
        Tags: pulumi.StringMap{
            "Name":   pulumi.String(computeArgs.Name),
            "JobUrl": pulumi.String(os.Getenv("TRAVIS_JOB_WEB_URL")),
        },
        InstanceType:        pulumi.String(size),
        KeyName:             pulumi.String(*computeArgs.KeyPair.Name),
        Ami:                 pulumi.String(*imageID),
        VpcSecurityGroupIds: network.ToStringArray(securityGroups),
    }

    if computeArgs.UserData != nil {
        ec2Args.UserData = pulumi.String(computeArgs.UserData.Content)
    }

    if computeArgs.Vpc != nil {
        ec2Args.SubnetId = computeArgs.Vpc.SubnetResults[0].ID()
    }

    server, err := ec2.NewInstance(ctx, computeArgs.Name, ec2Args)

    if err != nil {
        return nil, err
    }

    infra := Infrastructure{
        Groups:   ec2SecurityGroups,
        Server:   server,
        ImageID:  imageID,
        UserData: computeArgs.UserData,
    }

    if exports != nil {
        exports(ctx, &infra)
    }

    return &infra, nil
}

//CreateWireguardVM creates a wireguard ec2 aws instance
func CreateWireguardVM(ctx *pulumi.Context, computeArgs *model.ComputeArgs, exports exportsFnc) (*model.ComputeResult, error) {
    userData, err := shared.WireguardUserData()
    if err != nil {
        return nil, err
    }

    computeArgs.UserData = userData

    infra, err := CreateServer(ctx, computeArgs, exports)

    if err != nil {
        return nil, err
    }

    return &model.ComputeResult{
        Compute: infra.Server.CustomResourceState,
    }, err
}

func ProvisionVM(ctx *pulumi.Context, provisionArgs *model.ProvisionArgs, actor actors.Connector) error {

    server, err := ec2.GetInstance(ctx, "wireguard2", provisionArgs.SourceCompute.ID(), &ec2.InstanceState{
        InstanceState: pulumi.String("running"),
    })

    if err != nil {
        return err
    }

    provision := server.PublicIp.ApplyT(func(hostip string) string {
        var result string
        if actor != nil {
            result = actor.Connect(hostip)
            defer actor.Stop()
        }
        return strings.TrimSuffix(result, "\r\n")
    })

    ctx.Export(provisionArgs.ExportName, provision)

    return nil
}

// CreateImage creates an virtual machine image from an running VM.
func CreateImage(ctx *pulumi.Context, imageArgs model.ImageArgs, actor actors.Connector) error {

    server, err := ec2.GetInstance(ctx, "wireguard2", imageArgs.SourceCompute.ID(), &ec2.InstanceState{
        InstanceState: pulumi.String("running"),
    })

    if err != nil {
        return err
    }

    provision := server.PublicIp.ApplyT(func(hostip string) string {
        var result string
        if actor != nil {
            result = actor.Connect(hostip)
            defer actor.Stop()
        }

        //TODO implement the NewAmiFromInstance logic as an actor as well

        _, err = ec2.NewAmiFromInstance(ctx, imageArgs.Name, &ec2.AmiFromInstanceArgs{
            SourceInstanceId:      imageArgs.SourceCompute.ID(),
            Name:                  pulumi.String(imageArgs.Name),
            SnapshotWithoutReboot: pulumi.Bool(false),
        }, pulumi.IgnoreChanges([]string{"sourceInstanceId"}))

        if err != nil {
            panic(fmt.Errorf("failed to create ami image: %s", err))
        }

        return result
    })

    ctx.Export("Provisioning", provision)

    return nil
}
