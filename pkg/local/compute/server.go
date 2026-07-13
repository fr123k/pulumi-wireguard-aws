package compute

import (
	"fmt"
	"strings"
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Infrastructure represents the result of provisioning a physical server
type Infrastructure struct {
	ServerIP   string
	ServerName string
	UserData   *model.UserData
}

// LocalComputeArgs defines the input parameters for provisioning a physical server
type LocalComputeArgs struct {
	// Host is the IP address or hostname of the physical server
	Host string
	// Port is the SSH port (default: 22)
	Port int
	// Username is the SSH username
	Username string
	// KeyPair contains the SSH key pair for authentication
	KeyPair *model.KeyPairArgs
	// UserData is the rendered cloud-init script to apply
	UserData *model.UserData
	// ServerName is a descriptive name for the server
	ServerName string
	// SSHTimeout is the timeout for SSH connections
	SSHTimeout time.Duration
	// ProvisionCommands are additional commands to run after cloud-init
	ProvisionCommands []actors.SSHCommand
}

// ExportsFn is the callback type for exporting infrastructure details
type ExportsFn = func(ctx *pulumi.Context, infra *Infrastructure)

// ProvisionServer provisions a physical server via SSH using the actor pattern.
// It connects to the server, waits for it to be ready, and runs provisioning commands.
// This follows the same pattern used by the Hetzner cloud provisioners but connects
// directly to a known physical server IP instead of waiting for a cloud VM to be created.
func ProvisionServer(ctx *pulumi.Context, args *LocalComputeArgs, exportsFn ExportsFn) (*Infrastructure, error) {
	if args.Port == 0 {
		args.Port = 22
	}
	if args.SSHTimeout == 0 {
		args.SSHTimeout = 5 * time.Minute
	}

	if args.UserData != nil {
		ctx.Export("cloud-init", pulumi.String(args.UserData.Content))
	}

	infra := &Infrastructure{
		ServerName: args.ServerName,
		UserData:   args.UserData,
		ServerIP:   args.Host,
	}

	if exportsFn != nil {
		exportsFn(ctx, infra)
	}

	return infra, nil
}

// RunProvisioner connects to the physical server via SSH and executes provisioning commands.
// This should be called after ProvisionServer, typically within a pulumi.Run callback.
// The server IP must be reachable before calling this function.
func RunProvisioner(ctx *pulumi.Context, args *LocalComputeArgs, exportName string) error {
	connector := actors.NewSSHConnector(
		actors.SSHConnectorArgs{
			Port:       args.Port,
			Username:   args.Username,
			Timeout:    args.SSHTimeout,
			SSHKeyPair: *args.KeyPair.SSHKeyPair,
			Commands:   args.ProvisionCommands,
		},
		utility.Logger{Ctx: ctx},
	)

	result := connector.Connect(args.Host)
	defer connector.Stop()

	ctx.Export(exportName, pulumi.String(strings.TrimSuffix(result, "\r\n")))
	return nil
}

// WaitForServer waits for a physical server to become reachable via SSH.
// Uses the same retry logic as the cloud VM provisioners.
// Returns once the server responds to an SSH connection.
func WaitForServer(ctx *pulumi.Context, args *LocalComputeArgs) error {
	connector := actors.NewSSHConnector(
		actors.SSHConnectorArgs{
			Port:       args.Port,
			Username:   args.Username,
			Timeout:    args.SSHTimeout,
			SSHKeyPair: *args.KeyPair.SSHKeyPair,
			Commands: []actors.SSHCommand{
				{Command: "echo 'server ready'", Output: false},
			},
		},
		utility.Logger{Ctx: ctx},
	)

	result := connector.Connect(args.Host)
	defer connector.Stop()

	_ = ctx.Log.Info(fmt.Sprintf("Server %s is ready: %s", args.Host, strings.TrimSuffix(result, "\r\n")), nil)
	return nil
}