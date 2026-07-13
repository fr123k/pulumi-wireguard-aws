package ssh

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
	"golang.org/x/crypto/ssh"
)

// SSHClientConfig define type to pass ssh client configuration parameters.
type SSHClientConfig struct {
	Hostname      string
	Port          int `default:"22"`
	Username      string
	SSHKeyPair    SSHKey
	IgnoreHostKey bool          `default:"true"`
	Timeout       time.Duration `default:"30000"`
	Log           utility.Logger
}

type SSHKey struct {
	PublicKeyStr *string
	PrivateKey   crypto.PrivateKey
}

func publicKeyFile(keyPair SSHKey) (ssh.AuthMethod, error) {
	key, err := ssh.NewSignerFromKey(keyPair.PrivateKey)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

// SSHSession opens an SSH connection with retry and returns a client and session.
// It retries until the configured Timeout is reached when dialing fails
// (e.g. due to temporary connection issues or rate limiting by fail2ban).
func (sshClientConfig SSHClientConfig) SSHSession() (*ssh.Client, *ssh.Session, error) {
	publicKey, err := publicKeyFile(sshClientConfig.SSHKeyPair)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read key file: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: sshClientConfig.Username,
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Use the configured Timeout as the overall deadline for retries.
	// Default to 5 minutes if not set.
	timeout := sshClientConfig.Timeout
	if timeout == 0 {
		timeout = 5 * time.Minute
	}
	deadline := time.Now().Add(timeout)

	var connection *ssh.Client
	attempt := 0
	for {
		connection, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshClientConfig.Hostname, sshClientConfig.Port), sshConfig)
		if err == nil {
			break
		}

		remaining := time.Until(deadline)
		if remaining <= 0 {
			sshClientConfig.Log.Error("failed to dial %s:%d within %v: %s", sshClientConfig.Hostname, sshClientConfig.Port, timeout, err)
			return nil, nil, fmt.Errorf("failed to dial %s:%d within %v: %w", sshClientConfig.Hostname, sshClientConfig.Port, timeout, err)
		}

		// Exponential backoff: 2s, 4s, 8s, ... up to 30s max, but never exceed the deadline.
		backoff := time.Duration(attempt+1) * 2 * time.Second
		if backoff > 30*time.Second {
			backoff = 30 * time.Second
		}
		if backoff > remaining {
			backoff = remaining
		}

		sshClientConfig.Log.Info("SSH connection attempt %d failed (%v), retrying in %v (deadline in %v): %s", attempt+1, timeout, backoff, remaining.Round(time.Second), err)
		time.Sleep(backoff)
		attempt++
	}

	session, err := connection.NewSession()
	if err != nil {
		_ = connection.Close()
		sshClientConfig.Log.Error("failed to create session: %s", err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return connection, session, nil
}

// runSession is a helper that runs a command on the remote session, capturing both
// stdout and stderr. It returns the combined output and any error with details.
func (sshClientConfig SSHClientConfig) runSession(command string, setupSession func(s *ssh.Session) error) (string, error) {
	client, session, err := sshClientConfig.SSHSession()
	if err != nil {
		return "", err
	}
	defer func() { _ = client.Close() }()
	defer func() { _ = session.Close() }()

	if err := setupSession(session); err != nil {
		return "", err
	}

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	sshClientConfig.Log.Info("Run SSH command: %s", command)

	err = session.Run(command)
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return stdoutStr, fmt.Errorf("SSH command '%s' failed (exit code: %v): %s\nstdout: %s\nstderr: %s", command, err, err, truncate(stdoutStr, 500), truncate(stderrStr, 1000))
	}

	sshClientConfig.Log.Info("SSH command '%s' completed successfully", command)
	if stdoutStr != "" {
		sshClientConfig.Log.Debug("stdout: %s", truncate(stdoutStr, 200))
	}
	if stderrStr != "" {
		sshClientConfig.Log.Debug("stderr: %s", truncate(stderrStr, 200))
	}

	return stdoutStr, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (sshClientConfig SSHClientConfig) SSHCommand(command string) (*string, error) {
	result, err := sshClientConfig.runSession(command, func(s *ssh.Session) error {
		// Request PTY for interactive commands
		modes := ssh.TerminalModes{
			ssh.ECHO:  0,    // disable echoing
			ssh.IGNCR: 1,    // Ignore CR on input.
		}
		return s.RequestPty("xterm", 80, 40, modes)
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SSHCommandWithStdin runs a command on the remote host, piping the provided
// stdinContent to the command's stdin. This is useful for executing scripts
// without writing them to a temporary file on the remote host first.
// No PTY is requested for stdin-piped commands to avoid interference.
func (sshClientConfig SSHClientConfig) SSHCommandWithStdin(command string, stdinContent string) (*string, error) {
	result, err := sshClientConfig.runSession(command, func(s *ssh.Session) error {
		// No PTY for stdin commands — PTY interferes with piped input
		stdinPipe, err := s.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to create stdin pipe: %w", err)
		}
		go func() {
			_, _ = io.WriteString(stdinPipe, stdinContent)
			_ = stdinPipe.Close()
		}()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GenerateKeyPair() *SSHKey {
	reader := rand.Reader
	bitSize := 4096

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	savePEMKey("private.pem", key)

	publicKeyStr := exportPublicKeyAsPemStr(key)

	return &SSHKey{
		PrivateKey:   key,
		PublicKeyStr: &publicKeyStr,
	}
}

func parsePrivateKeyFile(file string) (crypto.PrivateKey, error) {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ReadPrivateKey(privateKeyFile string) *SSHKey {
	privateKey, err := parsePrivateKeyFile(privateKeyFile)
	checkError(err)

	publicKeyStr := exportPublicKeyAsPemStr(privateKey)

	return &SSHKey{
		PrivateKey:   privateKey,
		PublicKeyStr: &publicKeyStr,
	}
}

func ReadKeyPair(privateKeyFile string, publicKeyFile string) *SSHKey {
	publicKeyStr, err := utility.ReadFile(publicKeyFile)
	checkError(err)
	privateKey, err := parsePrivateKeyFile(privateKeyFile)
	checkError(err)

	return &SSHKey{
		PrivateKey:   privateKey,
		PublicKeyStr: publicKeyStr,
	}
}

func savePEMKey(fileName string, key *rsa.PrivateKey) *pem.Block {
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return privateKey
}

func exportPublicKeyAsPemStr(privateKey crypto.PrivateKey) string {
	var cryptoPublicKey crypto.PublicKey

	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		cryptoPublicKey = &k.PublicKey
	case *ed25519.PrivateKey:
		cryptoPublicKey = k.Public()
	default:
		checkError(fmt.Errorf("unsupported private key type: %T", privateKey))
	}

	pub, err := ssh.NewPublicKey(cryptoPublicKey)
	checkError(err)

	sshPubKey := base64.StdEncoding.EncodeToString(pub.Marshal())

	keyType := pub.Type()
	return fmt.Sprintf("%s %s", keyType, sshPubKey)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
