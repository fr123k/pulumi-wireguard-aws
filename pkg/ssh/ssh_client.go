package ssh

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHClientConfig define type to pass ssh client configuration parameters.
type SSHClientConfig struct {
	Hostname string
	Port int `default:22`
	Username string
	PrivateKeyFileName string
	IgnoreHostKey bool `default:true`
	Timeout time.Duration `default:30000`
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

// SSHSession open ann ssh client session
func (sshClientConfig SSHClientConfig) SSHSession() (*ssh.Session, error) {
	publicKey, err := publicKeyFile(sshClientConfig.PrivateKeyFileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read key file: %s", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: sshClientConfig.Username,
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 10 * time.Second,
	}

	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshClientConfig.Hostname, sshClientConfig.Port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:           0,     // disable echoing
		ssh.IGNCR:          1, // Ignore CR on input.
	}

	// stdin, err := session.StdinPipe()
	// if err != nil {
	// 	panic(fmt.Errorf("Unable to setup stdin for session: %v", err))
	// }
	// go io.Copy(stdin, os.Stdin)

	// stdout, err := session.StdoutPipe()
	// if err != nil {
	// 	panic(fmt.Errorf("Unable to setup stdout for session: %v", err))
	// }
	// go io.Copy(os.Stdout, stdout)

	// stderr, err := session.StderrPipe()
	// if err != nil {
	// 	panic(fmt.Errorf("Unable to setup stderr for session: %v", err))
	// }
    // go io.Copy(os.Stderr, stderr)

	
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}
	return session, err
}

func (sshClientConfig SSHClientConfig) SSHCommand(command string) (*string, error) {

	session, err := sshClientConfig.SSHSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to run ssh command: %s", err)
	}
	defer session.Close()

	result, err := session.Output(command)
	if err != nil {
		return nil, fmt.Errorf("Failed to run ssh command: %s", err)
	}
    fmt.Printf("Result: %s", result)
    
	str := string(result)
	return &str, nil
}
