package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/actors"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/local/compute"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		return createInfraStructure(ctx)
	})
}

// setEnvVars sets multiple environment variables and returns the first error encountered.
func setEnvVars(keyValues ...string) error {
	for i := 0; i < len(keyValues); i += 2 {
		if err := os.Setenv(keyValues[i], keyValues[i+1]); err != nil {
			return err
		}
	}
	return nil
}

func exports(ctx *pulumi.Context, infra *compute.Infrastructure) {
	ctx.Export("serverIp", pulumi.String(infra.ServerIP))
	ctx.Export("serverName", pulumi.String(infra.ServerName))
}

func createInfraStructure(ctx *pulumi.Context) error {
	cfg := config.New(ctx, "")

	// Required: the IP address or hostname of the physical server
	serverIP := cfg.Get("server_ip")
	if serverIP == "" {
		return fmt.Errorf("configuration key 'server_ip' is required (set with: pulumi config set server_ip <ip>)")
	}

	// Optional: SSH port (default: 22)
	sshPort := 22
	if port := cfg.GetInt("ssh_port"); port != 0 {
		sshPort = port
	}

	// Optional: SSH username (default: root)
	username := cfg.Get("username")
	if username == "" {
		username = "root"
	}

	// Optional: SSH private key file
	keyFile := cfg.Get("ssh_key_file")

	// Optional: Docker installation (default: true)
	installDocker := "true"
	if d := cfg.Get("install_docker"); d != "" {
		installDocker = d
	}

	// Optional: Network interface
	nic := cfg.Get("nic")

	// Optional: Server name
	serverName := cfg.Get("server_name")
	if serverName == "" {
		serverName = "minipc"
	}

	// Optional: SSH username for the created user
	minipcUser := cfg.Get("minipc_user")
	if minipcUser == "" {
		minipcUser = "frank.ittermann"
	}

	// Optional: Franky version
	frankyVersion := cfg.Get("franky_version")
	if frankyVersion == "" {
		frankyVersion = "0.32.1"
	}

	// Set environment variables for template rendering
	if err := setEnvVars(
		"MINIPC_USER", minipcUser,
		"MINIPC_SSH_PORT", fmt.Sprintf("%d", sshPort),
		"MINIPC_DOCKER", installDocker,
		"MINIPC_NIC", nic,
		"FRANKY_VERSION", frankyVersion,
		"SSH_USER", username,
	); err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}

	// Setup SSH key pair
	keyPairName := serverName + "-"
	var keyPair *model.KeyPairArgs

	if _, err := os.Stat(keyFile); err == nil {
		keyPair = model.NewKeyPairArgsWithPrivateKeyFile(&keyPairName, keyFile)
		fmt.Printf("Using local SSH key file %s\n", keyFile)
	} else {
		keyPair = model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
		fmt.Println("Using generated SSH key")
	}

	keyPair.Username = username

	// Render cloud-init template for the mini PC
	userData, err := shared.MiniPCUserData()
	if err != nil {
		return fmt.Errorf("failed to render cloud-init template: %w", err)
	}

	// Create local compute args
	computeArgs := &compute.LocalComputeArgs{
		Host:       serverIP,
		Port:       sshPort,
		Username:   username,
		KeyPair:    keyPair,
		UserData:   userData,
		ServerName: serverName,
		SSHTimeout: 5 * time.Minute,
	}

	// Provision the server (export cloud-init and server info)
	infra, err := compute.ProvisionServer(ctx, computeArgs, exports)
	if err != nil {
		return err
	}

	// Connect via SSH and run provisioning commands:
	// Pipe the cloud-init script directly to the server (bare metal doesn't auto-run cloud-init).
	// The script will restart SSH and configure the firewall, so we expect the connection to drop.
	err = compute.RunProvisioner(ctx, &compute.LocalComputeArgs{
		Host:     infra.ServerIP,
		Port:     sshPort,
		Username: username,
		KeyPair:  keyPair,
		ProvisionCommands: []actors.SSHCommand{
			{Command: "sudo bash", StdinContent: userData.Content, Output: false},
		},
		SSHTimeout: 10 * time.Minute,
	}, "minipc.provision")
	if err != nil {
		return err
	}

	// After provisioning, SSH back in (the server may have restarted SSH/firewall)
	// and verify that the services are running.
	err = compute.RunProvisioner(ctx, &compute.LocalComputeArgs{
		Host:     infra.ServerIP,
		Port:     sshPort,
		Username: username,
		KeyPair:  keyPair,
		ProvisionCommands: []actors.SSHCommand{
			{Command: "systemctl is-active franky", Output: true},
			{Command: "systemctl is-active nginx", Output: true},
		},
		// Generous timeout because fail2ban may penalize the IP after the
		// script reconfigured SSH and the firewall.
		SSHTimeout: 5 * time.Minute,
	}, "minipc.serviceStatus")
	if err != nil {
		return err
	}

	return nil
}