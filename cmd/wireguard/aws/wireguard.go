package main

import (
    "os"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/compute"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        return createInfraStructure(ctx)
    })
}

func exports(ctx *pulumi.Context, infra *compute.Infrastructure) {
    ctx.Export("publicIp", infra.Server.PublicIp)
    ctx.Export("publicDns", infra.Server.PublicDns)
}

func createInfraStructure(ctx *pulumi.Context) error {
    cfg := config.New(ctx, "")
    security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), model.VPCArgsDefault)
    security.Println()

    vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
    if err != nil {
        return err
    }

    keyPairName := "wireguard-"
    keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
    computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
    computeArgs.Name = "wireguard"
    computeArgs.Images = []*model.ImageArgs{
        {
            Name:   "wireguard-ami",
            Owners: []string{"self"},
            States: []string{"available"},
        },
        {
            Name:   "ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*",
            Owners: []string{"099720109477"},
        },
    }

    tags := map[string]string{
        "JobUrl":         os.Getenv("TRAVIS_JOB_WEB_URL"),
        "Project":        "wireguard",
        "pulumi-managed": "True",
    }

    externalSecurityGroup := model.SecurityGroup{
        Name:        "wireguard-external",
        Description: "Pulumi Managed. Allow Wireguard client traffic from internet.",
        Tags:        tags,
        IngressRules: []*model.SecurityRule{
            model.AllowOnePortRule("udp", 51820),
            model.AllowSSHRule(security),
        },
        EgressRules: []*model.SecurityRule{
            model.AllowAllRule(),
        },
    }
    //The order is important the referenced security groups has to be first.
    computeArgs.SecurityGroups = []*model.SecurityGroup{
        &externalSecurityGroup,
        {
            Name:        "wireguard-admin",
            Description: "Pulumi Managed. Allow admin traffic internal resources from VPN",
            Tags:        tags,
            IngressRules: []*model.SecurityRule{
                model.AllowAllRuleSecGroup(&externalSecurityGroup),
                model.AllowICMPRule(&externalSecurityGroup),
            },
            EgressRules: []*model.SecurityRule{
                model.AllowAllRule(),
            },
        },
    }

    vm, err := compute.CreateWireguardVM(ctx, computeArgs, exports)

    if err != nil {
        return err
    }

    sshConnector := shared.WireguardProvisioner(ctx, keyPair)

    compute.ProvisionVM(ctx, &model.ProvisionArgs{
        ExportName:    "wireguard.publicKey",
        SourceCompute: vm,
    }, &sshConnector)

    return err
}
