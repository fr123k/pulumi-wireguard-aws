package main

import (
    "github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
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

        vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
        if err != nil {
            return err
        }
        keyPairName := "wireguard-"
        
        keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
        //Uncomment to enable ssh access for debugging
        // keyPair := model.NewKeyPairArgsWithPrivateKeyFile(&keyPairName, "./development.pem")
        // keyPair.Name = &keyPairName

        keyPair.Username = "frank.ittermann"
        computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
        computeArgs.Name = "wireguard"
        computeArgs.Images = []*model.ImageArgs{
            {
                Name: "ubuntu-20.04",
            },
        }

        vm, err := compute.CreateWireguardVM(ctx, computeArgs)

        if err != nil {
            return err
        }

        sshConnector := shared.WireguardProvisioner(ctx, keyPair)

        //TODO implement exporting of mutliptl ssh output with one session
        compute.ProvisionVM(ctx, &model.ProvisionArgs{
            ExportName:    "wireguard.publicKey",
            SourceCompute: vm,
        }, &sshConnector)

        sshConnectorPassword := shared.WireguardPasswordProvisioner(ctx, keyPair)
        compute.ProvisionVM(ctx, &model.ProvisionArgs{
            ExportName:    "wireguard.password",
            SourceCompute: vm,
        }, &sshConnectorPassword)

        return err
    })
}
