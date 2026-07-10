package shared

import (
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// MiniPCUserData returns the rendered cloud-init userdata for a mini PC server.
// Template variables can be provided via environment variables.
func MiniPCUserData() (*model.UserData, error) {
	userDataVariables := map[string]string{
		"MINIPC_USER":      "MINIPC_USER",
		"MINIPC_SSH_PORT":  "MINIPC_SSH_PORT",
		"MINIPC_DOCKER":    "MINIPC_DOCKER",
		"MINIPC_NIC":       "MINIPC_NIC",
		"FRANKY_VERSION":   "FRANKY_VERSION",
		"SSH_USER":         "SSH_USER",
	}

	userData, err := model.NewUserData("cloud-init/minipc.txt", model.TemplateVariablesEnvironment(userDataVariables))
	if err != nil {
		return nil, err
	}
	return userData, nil
}

// MiniPCProvisioner returns an SSH connector that runs the cloud-init
// script directly on the server and then checks the franky service status.
// On bare metal servers, cloud-init doesn't automatically process user data,
// so we pipe the script via stdin to sudo bash instead.
func MiniPCProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs, scriptContent string) actors.SSHConnector {
	return actors.NewSSHConnector(
		actors.SSHConnectorArgs{
			Port:       22,
			Username:   keyPair.Username,
			Timeout:    10 * time.Minute,
			SSHKeyPair: *keyPair.SSHKeyPair,
			Commands: []actors.SSHCommand{
				{
					Command:     "sudo bash",
					StdinContent: scriptContent,
					Output:      false,
				},
				{
					Command: "systemctl is-active franky",
					Output:  true,
				},
			},
		},
		utility.Logger{
			Ctx: ctx,
		},
	)
}