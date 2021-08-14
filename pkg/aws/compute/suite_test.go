package compute

import (
    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type InfrastructureArgsFnc = func (ctx *pulumi.Context) (*model.ComputeArgs, error)

func DefaultComputeArgs(ctx *pulumi.Context) (*model.ComputeArgs, error) {
    security := model.NewSecurityArgsForVPC(true, model.VPCArgsDefault)
    security.Println()

    vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
    if err != nil {
        return nil, err
    }

    userDataVariables := map[string]string{
        "{{ CLIENT_PUBLICKEY }}":        "CLIENT_PUBLICKEY",
        "{{ CLIENT_IP_ADDRESS }}":       "CLIENT_IP_ADDRESS",
        "{{ MAILJET_API_CREDENTIALS }}": "MAILJET_API_CREDENTIALS",
        "{{ METADATA_URL }}":            "METADATA_URL",
    }

    userData, err := model.NewUserData("cloud-init/user-data.txt", model.TemplateVariablesEnvironment(userDataVariables))
    if err != nil {
        return nil, err
    }

    keyPairName := "wireguard-"
    keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
    computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
    computeArgs.UserData = userData
    computeArgs.Name = "wireguard"
    computeArgs.Images = []*model.ImageArgs{
        model.SelfImage("wireguard-ami"),
        {
            Name:   "ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*",
            Owners: []string{"099720109477"},
        },
    }

    tags := map[string]string{
        "JobUrl":         "travis_job_url",
        "Project":        "wireguard",
        "pulumi-managed": "True",
    }

    externalSecurityGroup := model.SecurityGroup{
        Name:        "wireguard-external",
        Description: "Pulumi Managed. Allow Wireguard client traffic from internet.",
        Tags:        tags,
        IngressRules: []*model.SecurityRule{
            model.AllowOnePortRule("udp", 51820),
            model.AllowSSHRule(security),
        },
        EgressRules: []*model.SecurityRule{
            model.AllowAllRule(),
        },
    }
    //The order is important the referenced security groups has to be first.
    computeArgs.SecurityGroups = []*model.SecurityGroup{
        &externalSecurityGroup,
        {
            Name:        "wireguard-admin",
            Description: "Pulumi Managed. Allow admin traffic internal resources from VPN",
            Tags:        tags,
            IngressRules: []*model.SecurityRule{
                model.AllowAllRuleSecGroup(&externalSecurityGroup),
                model.AllowICMPRule(&externalSecurityGroup),
            },
            EgressRules: []*model.SecurityRule{
                model.AllowAllRule(),
            },
        },
    }
    return computeArgs, nil
}

func DefaultComputeArgs2(ctx *pulumi.Context) (*model.ComputeArgs, error) {
    security := model.NewSecurityArgsForVPC(false, model.VPCArgsDefault)
    security.Println()

    vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
    if err != nil {
        return nil, err
    }

    userDataVariables := map[string]string{
        "{{ CLIENT_PUBLICKEY }}":        "CLIENT_PUBLICKEY",
        "{{ CLIENT_IP_ADDRESS }}":       "CLIENT_IP_ADDRESS",
        "{{ MAILJET_API_CREDENTIALS }}": "MAILJET_API_CREDENTIALS",
        "{{ METADATA_URL }}":            "METADATA_URL",
    }

    userData, err := model.NewUserData("cloud-init/user-data.txt", model.TemplateVariablesEnvironment(userDataVariables))
    if err != nil {
        return nil, err
    }

    keyPairName := "wireguard-"
    keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
    computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)
    computeArgs.UserData = userData
    computeArgs.Name = "wireguard"
    computeArgs.Images = []*model.ImageArgs{
        model.SelfImage("wireguard-ami"),
        {
            Name:   "ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*",
            Owners: []string{"099720109477"},
        },
    }

    tags := map[string]string{
        "JobUrl":         "travis_job_url",
        "Project":        "wireguard",
        "pulumi-managed": "True",
    }

    externalSecurityGroup := model.SecurityGroup{
        Name:        "wireguard-external",
        Description: "Pulumi Managed. Allow Wireguard client traffic from internet.",
        Tags:        tags,
        IngressRules: []*model.SecurityRule{
            model.AllowOnePortRule("udp", 51820),
            model.AllowSSHRule(security),
        },
        EgressRules: []*model.SecurityRule{
            model.AllowAllRule(),
        },
    }
    //The order is important the referenced security groups has to be first.
    computeArgs.SecurityGroups = []*model.SecurityGroup{
        &externalSecurityGroup,
        {
            Name:        "wireguard-admin",
            Description: "Pulumi Managed. Allow admin traffic internal resources from VPN",
            Tags:        tags,
            IngressRules: []*model.SecurityRule{
                model.AllowAllRuleSecGroup(&externalSecurityGroup),
                model.AllowICMPRule(&externalSecurityGroup),
            },
            EgressRules: []*model.SecurityRule{
                model.AllowAllRule(),
            },
        },
    }
    return computeArgs, nil
}
