package shared

import (
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func FrankyUserData() (*model.UserData, error) {
	// sbx.txt has no template variables currently, but keep the function
	// signature consistent for future use
	return model.NewUserDataNoVariables("cloud-init/franky.txt")
}

func FrankyProvisioner(ctx *pulumi.Context, keyPair *model.KeyPairArgs) actors.SSHConnector {
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
