package main

import (
    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/compute"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        cfg := config.New(ctx, "")
        security := model.NewSecurityArgsForVPC(cfg.GetBool("vpn_enabled_ssh"), model.VPCArgsDefault)
        security.Println()

        keyPairName := "wireguard-packer"
        keyPair := model.NewKeyPairArgsWithKey(&keyPairName)
        vm, err := compute.CreateWireguardVM(ctx, model.NewComputeArgsWithSecurityAndKeyPair(security, keyPair), nil)

        if err != nil {
            return err
        }

        sshConnector := shared.WireguardProvisioner(ctx, keyPair)

        err = compute.CreateImage(ctx, model.ImageArgs{
            Name:          "wireguard-ami-new",
            SourceCompute: vm,
        }, &sshConnector)
        return err
    })
}
