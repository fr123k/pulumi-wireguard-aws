package main

import (
	"os"
	"time"

	wireguardCfg "github.com/fr123k/pulumi-wireguard-aws/cmd/wireguard/config"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const size = "t2.large"

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        cfg := config.New(ctx, "")
        security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), wireguardCfg.VPCArgsDefault)
        security.Println()

        vpc, err := network.CreateVPC(ctx, wireguardCfg.VPCArgsDefault)
        if err != nil {
            return err
        }

        keyPairName := "wireguard-"
        keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
        computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)

        tags := map[string]string{
            "JobUrl":         os.Getenv("TRAVIS_JOB_WEB_URL"),
            "Project":        "wireguard",
            "pulumi-managed": "True",
        }

        externalSecurityGroup := model.SecurityGroup{
            Name: "wireguard-external",
            Description: "Pulumi Managed. Allow Wireguard client traffic from internet.",
            Tags: tags,
            IngressRules: []*model.SecurityRule{{
                Protocol: "udp",
                SourcePort: 51820,
                DestinationPort: 51820,
                CidrBlocks: []string{"0.0.0.0/0"},
            },},
            EgressRules: []*model.SecurityRule{{
                Protocol: "-1",
                SourcePort: 0,
                DestinationPort: 0,
                CidrBlocks: []string{"0.0.0.0/0"},
            },},
        }
        //The order is important the referenced security groups has to be first.
        computeArgs.SecurityGroups = []*model.SecurityGroup{
            &externalSecurityGroup,
            {
                Name: "wireguard-admin",
                Description: "Terraform Managed. Allow admin traffic internal resources from VPN",
                Tags: tags,
                IngressRules: []*model.SecurityRule{{
                    Protocol: "-1",
                    SourcePort: 0,
                    DestinationPort: 0,
                    SecurityGroups: []*model.SecurityGroup{&externalSecurityGroup,},
                },{
                    Protocol: "icmp",
                    SourcePort: 8,
                    DestinationPort: 0,
                    SecurityGroups: []*model.SecurityGroup{&externalSecurityGroup,},
                },},
                EgressRules: []*model.SecurityRule{{
                    Protocol: "-1",
                    SourcePort: 0,
                    DestinationPort: 0,
                    CidrBlocks: []string{"0.0.0.0/0"},
                },},
            },
        }


        vm, err := compute.CreateWireguardVM(ctx, computeArgs)

        if err != nil {
            return err
        }

        sshConnector := actors.NewSSHConnector(
            actors.SSHConnectorArgs{
                Port:       22,
                Username:   "ubuntu",
                Timeout:    2 * time.Minute,
                SSHKeyPair: *keyPair.SSHKeyPair,
            },
        )

        compute.ProvisionVM(ctx,  &model.ProvisionArgs{
            ExportName:     "wireguard.publicKey",
            SourceCompute:  vm,
        }, &sshConnector)

        return err
    })
}
