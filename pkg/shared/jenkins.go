package shared

import (
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func JenkinsSecGroup(tags map[string]string, security *model.SecurityArgs) []*model.SecurityGroup {
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
    return []*model.SecurityGroup{
        &externalSecurityGroup,
    }
}

func JenkinsUserData(filename string) (*model.UserData, error) {
    userDataEnvVariables := map[string]string{
        "SEED_BRANCH_JOBS": "SEED_BRANCH_JOBS",
    }

    userDataPropertyVariables := map[string]string{
        "ADMIN_PASSWORD": utility.RandomSecret(32),
        "AWS_KEY_ID":     "undefined",
        "AWS_KEY_SECRET": "undefined",
    }

    userData, err := model.NewUserData(filename, append(model.TemplateVariablesEnvironment(userDataEnvVariables), model.TemplateVariablesStringProperty(userDataPropertyVariables)...))
    if err != nil {
        return nil, err
    }
    return userData, nil
}

func JenkinsProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs) actors.SSHConnector {
    return actors.NewSSHConnector(
        actors.SSHConnectorArgs{
            Port:       22,
            Username:   keyPair.Username,
            Timeout:    2 * time.Minute,
            SSHKeyPair: *keyPair.SSHKeyPair,
            Commands: []actors.SSHCommand{
                {
                    Command: "sudo cloud-init status --wait",
                    Output:  false,
                },
            },
        },
        utility.Logger{
            Ctx: ctx,
        },
    )
}
