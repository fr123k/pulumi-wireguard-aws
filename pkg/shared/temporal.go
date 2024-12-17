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
		"SECRET_OPERATOR_AUTHENTICATION_TOKEN":        "SECRET_OPERATOR_AUTHENTICATION_TOKEN",
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
