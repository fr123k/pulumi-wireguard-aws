package ssh

import (
    "context"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "fmt"
    "io/ioutil"
    "os"
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
    "github.com/pulumi/pulumi/sdk/v3/go/common/util/retry"
    "golang.org/x/crypto/ssh"
)

// SSHClientConfig define type to pass ssh client configuration parameters.
type SSHClientConfig struct {
    Hostname      string
    Port          int `default:22`
    Username      string
    SSHKeyPair    SSHKey
    IgnoreHostKey bool          `default:true`
    Timeout       time.Duration `default:30000`
    Log           utility.Logger
}

type SSHKey struct {
    PublicKeyStr *string
    PrivateKey   *rsa.PrivateKey
}

func publicKeyFile(keyPair SSHKey) (ssh.AuthMethod, error) {
    key, err := ssh.NewSignerFromKey(keyPair.PrivateKey)
    if err != nil {
        return nil, err
    }
    return ssh.PublicKeys(key), nil
}

// SSHSession open ann ssh client session
func (sshClientConfig SSHClientConfig) SSHSession() (*ssh.Session, error) {
    publicKey, err := publicKeyFile(sshClientConfig.SSHKeyPair)

    if err != nil {
        return nil, fmt.Errorf("failed to read key file: %s", err)
    }

    sshConfig := &ssh.ClientConfig{
        User: sshClientConfig.Username,
        Auth: []ssh.AuthMethod{
            publicKey,
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         10 * time.Second,
    }

    _, connection, err := retry.Until(context.Background(), retry.Acceptor{
        Accept: func(retries int, nextRetryTime time.Duration) (bool, interface{}, error) {
            con, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshClientConfig.Hostname, sshClientConfig.Port), sshConfig)
            if err != nil {
                if retries > 10 {
                    sshClientConfig.Log.Error("failed to dial: %s", err)
                    return true, nil, fmt.Errorf("failed to dial: %s", err)
                }
                sshClientConfig.Log.Debug("Retry ssh connection retries: %d, error: %s", retries, err)
                return false, nil, nil
            }
            return true, con, nil
        },
    })

    if err != nil || connection == nil {
        sshClientConfig.Log.Error("failed to dial: %s", err)
        return nil, fmt.Errorf("failed to dial: %s", err)
    }

    // connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshClientConfig.Hostname, sshClientConfig.Port), sshConfig)
    // if err != nil {
    // 	return nil, fmt.Errorf("Failed to dial: %s", err)
    // }

    session, err := connection.(*ssh.Client).NewSession()
    if err != nil {
        sshClientConfig.Log.Error("failed to create session: %s", err)
        return nil, fmt.Errorf("failed to create session: %s", err)
    }

    modes := ssh.TerminalModes{
        ssh.ECHO:  0, // disable echoing
        ssh.IGNCR: 1, // Ignore CR on input.
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
        return nil, fmt.Errorf("failed to run ssh command: %s", err)
    }
    defer session.Close()

    sshClientConfig.Log.Info("Run SSH command to %s", command)

    result, err := session.Output(command)
    if err != nil {
        return nil, fmt.Errorf("failed to run ssh command: %s", err)
    }
    sshClientConfig.Log.Info("Result: %s", result)

    fmt.Printf("Result: %s", result)

    str := string(result)
    return &str, nil
}

func GenerateKeyPair() *SSHKey {
    reader := rand.Reader
    bitSize := 4096

    key, err := rsa.GenerateKey(reader, bitSize)
    checkError(err)

    publicKey := key.PublicKey

    savePEMKey("private.pem", key)

    publicKeyStr := exportRsaPublicKeyAsPemStr(&publicKey)

    return &SSHKey{
        PrivateKey:   key,
        PublicKeyStr: &publicKeyStr,
    }
}

func parsePrivateKeyFile(file string) (*rsa.PrivateKey, error) {
    buffer, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }

    key, err := ssh.ParseRawPrivateKey(buffer)
    if err != nil {
        return nil, err
    }
    rsaKey := key.(*rsa.PrivateKey)
    return rsaKey, nil
}

func ReadPrivateKey(privateKeyFile string) *SSHKey {
    privateKey, err := parsePrivateKeyFile(privateKeyFile)
    checkError(err)

    publicKey := exportRsaPublicKeyAsPemStr(&privateKey.PublicKey)

    return &SSHKey{
        PrivateKey:   privateKey,
        PublicKeyStr: &publicKey,
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

func exportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
    // Generate the ssh public key
    pub, err := ssh.NewPublicKey(pubkey)
    checkError(err)

    // Encode to store to file
    sshPubKey := base64.StdEncoding.EncodeToString(pub.Marshal())

    return fmt.Sprintf("ssh-rsa %s", sshPubKey)
}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
