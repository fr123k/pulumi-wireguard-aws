package main

import (
	"fmt"
	"os"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
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
	userData, err := shared.JenkinsUserData("cloud-init/eucaptcha.yaml")
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

	//Uncomment to enable ssha access for debugging
	// keyPair := model.NewKeyPairArgsWithPrivateKeyFile(&keyPairName, "./development.pem")
	// keyPair.Name = &keyPairName

	keyPair.Username = "root"
	computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
	computeArgs.Name = "eucaptcha"
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
		"Project":        "eucaptcha",
		"pulumi-managed": "True",
	}

	computeArgs.SecurityGroups = shared.JenkinsSecGroup(tags, security)

	vm, err := compute.CreateServer(ctx, computeArgs, "10.8.0.145", exports)

	if err != nil {
		return err
	}

	sshConnector := shared.JenkinsProvisioner(ctx, keyPair)

	compute.ProvisionVM(ctx, "eucaptcha", &model.ProvisionArgs{
		ExportName: "eucaptcha.result",
		SourceCompute: &model.ComputeResult{
			Compute: vm.Server.CustomResourceState,
		},
	}, &sshConnector)

	opt0 := "fr123k.uk."
	opt1 := false
	selected, err := route53.LookupZone(ctx, &route53.LookupZoneArgs{
		Name:        &opt0,
		PrivateZone: &opt1,
	}, nil)
	if err != nil {
		return err
	}

	_, err = route53.NewRecord(ctx, "captcha", &route53.RecordArgs{
		ZoneId: pulumi.String(selected.ZoneId),
		Name:   pulumi.String(fmt.Sprintf("%v%v", "captcha.", selected.Name)),
		Type:   pulumi.String("A"),
		Ttl:    pulumi.Int(300),
		Records: pulumi.StringArray{
			vm.Server.Ipv4Address,
		},
	})

	if err != nil {
		return err
	}

	return err
}
