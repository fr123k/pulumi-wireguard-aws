package main

import (
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/compute"
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

        keyPairName := "wireguard-packer"
        keyPair := model.NewKeyPairArgsWithKey(&keyPairName)
        vm, err := compute.CreateWireguardVM(ctx, model.NewComputeArgsWithSecurityAndKeyPair(security, keyPair))

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
            utility.Logger{
                Ctx: ctx,
            },
        )

        err = compute.CreateImage(ctx, model.ImageArgs{
            Name:          "wireguard-ami-new",
            SourceCompute: vm,
        }, &sshConnector)
        return err
    })
}
