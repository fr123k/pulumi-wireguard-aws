package actors

import (
    "fmt"
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/ssh"
)

// SSHConnectorArgs this defines the ssh connection related arguments.
type SSHConnectorArgs struct {
    Port       int
    Username   string
    SSHKeyPair ssh.SSHKey
    Timeout    time.Duration
}

// SSHConnector the ssh implementation of the Connector actor
type SSHConnector struct {
    args *SSHConnectorArgs
    connector
}

// NewSSHConnector initialize an ssh connector
func NewSSHConnector(args SSHConnectorArgs) SSHConnector {
    sshConnector := SSHConnector{
        args: &args,
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
            Hostname:   address,
            Port:       c.args.Port,
            Username:   c.args.Username,
            SSHKeyPair: c.args.SSHKeyPair,
            Timeout:    c.args.Timeout,
        }

        fmt.Printf("Open SSH connection to %s", address)

        result, err := sshClient.SSHCommand("sudo cloud-init status --wait")
        if err != nil {
            panic(fmt.Errorf("Failed to run cmd : %s", err))
        }
        fmt.Printf("Result: %s", *result)

        result, err = sshClient.SSHCommand("sudo cat /tmp/server_publickey")
        if err != nil {
            panic(fmt.Errorf("Failed to run cmd : %s", err))
        }
        resultChan <- *result
    }
    return <-resultChan
}
