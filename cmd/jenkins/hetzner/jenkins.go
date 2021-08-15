package main

import (
    "fmt"
    "os"
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/compute"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/hetzner/network"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/ssh"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
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
    userDataEnvVariables := map[string]string{
        "{{ SEED_BRANCH_JOBS }}": "SEED_BRANCH_JOBS",
    }

    userDataSetVariables := map[string]string{
        "{{ ADMIN_PASSWORD }}": fmt.Sprintf("ADMIN_PASSWORD=%s", utility.RandomSecret(32)),
        "{{ AWS_KEY_ID }}":     "AWS_KEY_ID=undefined",
        "{{ AWS_KEY_SECRET }}": fmt.Sprintf("AWS_KEY_SECRET=%s", "undefined"),
    }

    userData, err := model.NewUserData("cloud-init/jenkins.yaml", append(model.TemplateVariablesEnvironment(userDataEnvVariables), model.TemplateVariablesString(userDataSetVariables)...))
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
    keyPair := model.NewKeyPairArgs(&keyPairName)
    computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
    computeArgs.Name = "jenkins-master"
    computeArgs.UserData = userData
    computeArgs.Images = []*model.ImageArgs{
        {
            Name: "45467990", //jenkins-master snapshots image id
        },
        {
            Name: "ubuntu-20.04",
        },
    }

    tags := map[string]string{
        "JobUrl":         os.Getenv("TRAVIS_JOB_WEB_URL"),
        "Project":        "jenkins",
        "pulumi-managed": "True",
    }

    externalSecurityGroup := model.SecurityGroup{
        Name:        "jenkins-security-group",
        Description: "Pulumi Managed.",
        Tags:        tags,
        IngressRules: []*model.SecurityRule{
            model.AllowOnePortRule("tcp", 80),
            model.AllowOnePortRule("tcp", 22).CidrBlock("95.90.244.46/32"),
            model.AllowSSHRule(security),
        },
        EgressRules: []*model.SecurityRule{
            model.AllowOnePortRule("tcp", 80),
            model.AllowOnePortRule("tcp", 443),
            model.AllowOnePortRule("tcp", 22),
            model.AllowOnePortRule("tcp", 22).CidrBlock("140.82.118.0/24"),
            model.AllowOnePortRule("tcp", 22).CidrBlock("140.82.121.4/32"),
            model.AllowOnePortRule("tcp", 22).CidrBlock("204.232.175.90/32"),
            model.AllowOnePortRule("tcp", 22).CidrBlock("207.97.227.239/32"),
        },
    }
    //The order is important the referenced security groups has to be first.
    computeArgs.SecurityGroups = []*model.SecurityGroup{
        &externalSecurityGroup,
    }

    vm, err := compute.CreateServer(ctx, computeArgs, exports)

    if err != nil {
        return err
    }

    sshKeys := ssh.ReadPrivateKey("/home/vagrant/.ssh/development.pem")
    sshConnector := actors.NewSSHConnector(
        actors.SSHConnectorArgs{
            Port:       22,
            Username:   "root",
            Timeout:    2 * time.Minute,
            SSHKeyPair: *sshKeys,
            Commands: []actors.SSHCommand{
                {
                    Command: "sudo cloud-init status --wait",
                    Output: false,
                },
            },
        },
        utility.Logger{
            Ctx: ctx,
        },
    )

    compute.ProvisionVM(ctx, &model.ProvisionArgs{
        ExportName: "jenkins.publicKey",
        SourceCompute: &model.ComputeResult{
            Compute: vm.Server.CustomResourceState,
        },
    }, &sshConnector)

    return err
}
