package actors

import (
	"fmt"
	"strings"
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/ssh"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
)

// SSHConnectorArgs this defines the ssh connection related arguments.
type SSHConnectorArgs struct {
	Port       int
	Username   string
	SSHKeyPair ssh.SSHKey
	Timeout    time.Duration
	Commands   []SSHCommand
}

// SSHConnector the ssh implementation of the Connector actor
type SSHConnector struct {
	args *SSHConnectorArgs
	connector
	log utility.Logger
}

type SSHCommand struct {
	Command     string
	Output      bool
	StdinContent string // If non-empty, pipe this content to the command's stdin
}

// NewSSHConnector initialize an ssh connector
func NewSSHConnector(args SSHConnectorArgs, log utility.Logger) SSHConnector {
	sshConnector := SSHConnector{
		args: &args,
		log:  log,
	}
	sshConnector.connector = newConnector()
	return sshConnector
}

// Connect this function is called when the virtual instance is created and can recevie connection.
// TODO implements retries ssh attemps because after an virtual machine is ready doesn't
func (c *SSHConnector) Connect(address string) string {
	resultChan := make(chan string)
	c.actions <- func() {
		sshClient := ssh.SSHClientConfig{
			Hostname:      address,
			Port:          c.args.Port,
			Username:      c.args.Username,
			SSHKeyPair:    c.args.SSHKeyPair,
			Timeout:       c.args.Timeout,
			IgnoreHostKey: true,
			Log:           c.log,
		}

		c.log.Info("Open SSH connection to %s with %s", address, sshClient.Username)

		var output strings.Builder
		for _, cmd := range c.args.Commands {
			var result *string
			var err error
			if cmd.StdinContent != "" {
				result, err = sshClient.SSHCommandWithStdin(cmd.Command, cmd.StdinContent)
			} else {
				result, err = sshClient.SSHCommand(cmd.Command)
			}
			if err != nil {
				c.log.Error("Failed to run cmd '%s' with error '%s'", cmd.Command, err)
				panic(fmt.Errorf("failed to run cmd '%s' with error '%s'", cmd.Command, err))
			}
			c.log.Info("Result: %s", *result)
			if cmd.Output {
				output.WriteString(*result)
			}
		}
		resultChan <- output.String()
	}
	return <-resultChan
}
