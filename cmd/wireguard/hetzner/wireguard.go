package main

import (
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const size = "t2.large"

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
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
                Name:   "ubuntu-20.04",
            },
        }

        vm, err := compute.CreateWireguardVM(ctx, computeArgs)

        if err != nil {
            return err
        }

        sshConnector := actors.NewSSHConnector(
            actors.SSHConnectorArgs{
                Port:       22,
                Username:   "root",
                Timeout:    2 * time.Minute,
                SSHKeyPair: *keyPair.SSHKeyPair,
                Commands: []actors.SSHCommand{
                    {
                        Command: "sudo cloud-init status --wait",
                        Output:  false,
                    },
                    {
                        Command: "sudo cat /tmp/server_publickey",
                        Output:  true,
                    },
                },
            },
            utility.Logger{
                Ctx: ctx,
            },
        )

        compute.ProvisionVM(ctx, &model.ProvisionArgs{
            ExportName:    "wireguard.publicKey",
            SourceCompute: vm,
        }, &sshConnector)

        return err
    })
}
