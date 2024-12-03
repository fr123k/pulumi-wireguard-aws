package main

import (
	"os"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func exports(ctx *pulumi.Context, infra *compute.Infrastructure) {
	ctx.Export("publicIp", infra.Server.Ipv4Address)
	ctx.Export("publicDns", infra.Server.Ipv4Address)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		return createInfraStructure(ctx)
	})
}

func createInfraStructure(ctx *pulumi.Context) error {
	config := config.New(ctx, "")

	//TODO fetch new created aws key and secret
	// secret := config.Require("secret")
	userData, err := shared.JenkinsUserData("cloud-init/jenkins.yaml")
	if err != nil {
		return err
	}

	security := model.NewSecurityArgsForVPC(config.GetBool("vpn_enabled_ssh"), model.VPCArgsDefault)
	security.Println()

	vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
	if err != nil {
		return err
	}

	keyPairName := "development"
	keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
	keyPair.Username = "root"
	computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
	computeArgs.Name = "jenkins-master"
	computeArgs.UserData = userData
	computeArgs.Images = []*model.ImageArgs{
		{
			Name: "ubuntu-20.04",
		},
		{
			Name: "45467990", //jenkins-master snapshots image id
		},
	}

	tags := map[string]string{
		"JobUrl":         os.Getenv("TRAVIS_JOB_WEB_URL"),
		"Project":        "jenkins",
		"pulumi-managed": "True",
	}

	computeArgs.SecurityGroups = shared.JenkinsSecGroup(tags, security)

	vm, err := compute.CreateServer(ctx, computeArgs, "10.8.0.145", exports)

	if err != nil {
		return err
	}

	sshConnector := shared.JenkinsProvisioner(ctx, keyPair)

	compute.ProvisionVM(ctx, "jenkins", &model.ProvisionArgs{
		ExportName: "jenkins.publicKey",
		SourceCompute: &model.ComputeResult{
			Compute: vm.Server.CustomResourceState,
		},
	}, &sshConnector)

	return err
}
