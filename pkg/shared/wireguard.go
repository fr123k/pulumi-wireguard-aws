package shared

import (
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func WireguardUserData() (*model.UserData, error) {
    userDataVariables := map[string]string{
        "CLIENT_PUBLICKEY":                "CLIENT_PUBLICKEY",
        "CLIENT_IP_ADDRESS":              "CLIENT_IP_ADDRESS",
        "MAILJET_API_CREDENTIALS":        "MAILJET_API_CREDENTIALS",
        "METADATA_URL":                   "METADATA_URL",
        "WIREGUARD_DOMAIN":               "WIREGUARD_DOMAIN",
        "SECRET_OPERATOR_AUTHENTICATION_TOKEN": "SECRET_OPERATOR_AUTHENTICATION_TOKEN",
    }

    userData, err := model.NewUserData("cloud-init/wireguard.txt", model.TemplateVariablesEnvironment(userDataVariables))
    if err != nil {
        return nil, err
    }
    return userData, nil
}

// WireguardPrebakedUserData returns the minimal cloud-init script for pre-baked WireGuard images.
// This script only handles runtime-specific configuration (SSH, secrets, SSL, WireGuard keys,
// service startup). All binaries, packages, and systemd units are already installed in the image.
func WireguardPrebakedUserData() (*model.UserData, error) {
    userDataVariables := map[string]string{
        "CLIENT_PUBLICKEY":                      "CLIENT_PUBLICKEY",
        "CLIENT_IP_ADDRESS":                     "CLIENT_IP_ADDRESS",
        "MAILJET_API_CREDENTIALS":                "MAILJET_API_CREDENTIALS",
        "WIREGUARD_DOMAIN":                       "WIREGUARD_DOMAIN",
        "SECRET_OPERATOR_AUTHENTICATION_TOKEN":   "SECRET_OPERATOR_AUTHENTICATION_TOKEN",
    }

    userData, err := model.NewUserData("cloud-init/wireguard-prebaked.txt", model.TemplateVariablesEnvironment(userDataVariables))
    if err != nil {
        return nil, err
    }
    return userData, nil
}

func WireguardProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs) actors.SSHConnector {
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
}

func WireguardPasswordProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs) actors.SSHConnector {
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
                {
                    Command: "sudo cat /tmp/user_password",
                    Output:  true,
                },
            },
        },
        utility.Logger{
            Ctx: ctx,
        },
    )
}
