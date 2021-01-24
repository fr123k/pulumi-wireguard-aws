package main

import (
	"time"

	wireguardCfg "github.com/fr123k/pulumi-wireguard-aws/cmd/wireguard/config"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
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
        vm, err := compute.CreateWireguardVM(ctx, model.NewComputeArgsWithKeyPair(vpc, security, keyPair))

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
