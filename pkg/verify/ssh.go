package verify

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConfig holds the configuration for SSH connections
type SSHConfig struct {
	Host       string
	Port       int
	User       string
	KeyPath    string
	Timeout    time.Duration
}

// Connect establishes an SSH connection using the provided configuration
func Connect(config SSHConfig) (*ssh.Client, error) {
	keyData, err := os.ReadFile(config.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         config.Timeout,
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	return client, nil
}
