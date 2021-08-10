package main

import (
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
        computeArgs.IngressRules = []*model.SecurityRule{{
                Protocol: "udp",
                SourcePort: 51820,
                DestinationPort: 51820,
                CidrBlocks: []string{"0.0.0.0/0"},
            },
        }
        computeArgs.EgressRules = []*model.SecurityRule{{
                Protocol: "-1",
                SourcePort: 0,
                DestinationPort: 0,
                CidrBlocks: []string{"0.0.0.0/0"},
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
