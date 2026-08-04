package shared

import (
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TemporalUserData() (*model.UserData, error) {
	userDataVariables := map[string]string{
		"SECRET_OPERATOR_TOKEN": "SECRET_OPERATOR_TOKEN",
		// "CLIENT_IP_ADDRESS":       "CLIENT_IP_ADDRESS",
		// "MAILJET_API_CREDENTIALS": "MAILJET_API_CREDENTIALS",
		// "METADATA_URL":            "METADATA_URL",
	}

	userData, err := model.NewUserData("cloud-init/temporal.txt", model.TemplateVariablesEnvironment(userDataVariables))
	if err != nil {
		return nil, err
	}
	return userData, nil
}

// TemporalPrebakedUserData returns the minimal cloud-init script for pre-baked Temporal images.
// This script only handles runtime-specific configuration (SSH, secrets, SSL, service startup).
func TemporalPrebakedUserData() (*model.UserData, error) {
	return prebakedUserData("cloud-init/temporal-prebaked.txt", map[string]string{
		"SECRET_OPERATOR_TOKEN": "SECRET_OPERATOR_TOKEN",
		"TEMPORAL_DOMAIN":       "TEMPORAL_DOMAIN",
		"DUNEBOT_DOMAIN":        "DUNEBOT_DOMAIN",
	})
}

func TemporalProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs) actors.SSHConnector {
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
				// {
				// 	Command: "sudo cat /tmp/server_publickey",
				// 	Output:  true,
				// },
			},
		},
		utility.Logger{
			Ctx: ctx,
		},
	)
}

// func TemporalPasswordProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs) actors.SSHConnector {
// 	return actors.NewSSHConnector(
// 		actors.SSHConnectorArgs{
// 			Port:       22,
// 			Username:   keyPair.Username,
// 			Timeout:    2 * time.Minute,
// 			SSHKeyPair: *keyPair.SSHKeyPair,
// 			Commands: []actors.SSHCommand{
// 				{
// 					Command: "sudo cloud-init status --wait",
// 					Output:  false,
// 				},
// 				{
// 					Command: "sudo cat /tmp/user_password",
// 					Output:  true,
// 				},
// 			},
// 		},
// 		utility.Logger{
// 			Ctx: ctx,
// 		},
// 	)
// }
