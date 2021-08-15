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
    Command string
    Output  bool
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
    resultChan := make(chan string, 0)
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

        c.log.Info("Open SSH connection to %s", address)
        var output strings.Builder
        for _, cmd := range c.args.Commands {
            result, err := sshClient.SSHCommand(cmd.Command)
            if err != nil {
                c.log.Error("Failed to run cmd '%s' with error '%s'", cmd.Command, err)
                panic(fmt.Errorf("Failed to run cmd '%s' with error '%s'", cmd.Command, err))
            }
            c.log.Info("Result: %s", *result)
            if cmd.Output {
                output.WriteString(*result)
            }
        }
        resultChan <- output.String()
        // result, err := sshClient.SSHCommand("sudo cloud-init status --wait")
        // if err != nil {
        //     c.log.Error("Failed to run cmd : %s", err)
        //     panic(fmt.Errorf("Failed to run cmd : %s", err))
        // }
        // c.log.Info("Result: %s", *result)

        // result, err = sshClient.SSHCommand("sudo cat /tmp/server_publickey")
        // if err != nil {
        //     c.log.Error("Failed to run cmd : %s", err)
        //     panic(fmt.Errorf("Failed to run cmd : %s", err))
        // }
        // resultChan <- *result
    }
    return <-resultChan
}
