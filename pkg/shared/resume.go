package shared

import (
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ResumeUserData returns the cloud-init script for the resume (GDPR API) server.
// Template variables are populated from environment variables at deploy time:
//   - GDPR_DOMAIN: the domain for the GDPR API endpoint (e.g. api.resume.fr123k.uk)
//   - GDPR_AUTH_TOKEN: the Bearer token for authenticating GDPR notification requests
func ResumeUserData() (*model.UserData, error) {
	userDataVariables := map[string]string{
		"GDPR_DOMAIN":     "GDPR_DOMAIN",
		"GDPR_AUTH_TOKEN": "GDPR_AUTH_TOKEN",
	}

	userData, err := model.NewUserData("cloud-init/resume.txt", model.TemplateVariablesEnvironment(userDataVariables))
	if err != nil {
		return nil, err
	}
	return userData, nil
}

// ResumeProvisioner returns an SSH connector that waits for cloud-init to complete
// on the resume server. SSH is behind wireguard, so the connector uses the
// wireguard peer IP to reach the server.
func ResumeProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs) actors.SSHConnector {
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