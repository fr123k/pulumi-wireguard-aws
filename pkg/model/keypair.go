package model

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/ssh"
)

type KeyPairArgs struct {
    Name         *string
    RandomNumber *int
    SSHKeyPair   *ssh.SSHKey
}

func NewKeyPairArgs(name *string) *KeyPairArgs {
    return &KeyPairArgs{Name: name}
}

func NewKeyPairArgsWithRandomName(name *string) *KeyPairArgs {
    randSrc := rand.NewSource(time.Now().UnixNano())
    randomNumber := rand.New(randSrc).Intn(100000)
    keyPairName := fmt.Sprintf("%s%d", *name, randomNumber)
    return &KeyPairArgs{
        Name:         &keyPairName,
        RandomNumber: &randomNumber,
    }
}

func NewKeyPairArgsWithRandomNameWithKeyFile(name *string, publicKeyFile *string) *KeyPairArgs {
    randSrc := rand.NewSource(time.Now().UnixNano())
    randomNumber := rand.New(randSrc).Intn(100000)
    keyPairName := fmt.Sprintf("%s%d", *name, randomNumber)
    return &KeyPairArgs{
        Name:         &keyPairName,
        RandomNumber: &randomNumber,
        //TODO pass private key file as well ?
        SSHKeyPair: ssh.ReadKeyPair("", *publicKeyFile),
    }
}

func NewKeyPairArgsWithKeyFile(name *string, privateKeyFile string, publicKeyFile string) *KeyPairArgs {
    return &KeyPairArgs{
        Name:       name,
        SSHKeyPair: ssh.ReadKeyPair(privateKeyFile, publicKeyFile),
    }
}

func NewKeyPairArgsWithKey(name *string) *KeyPairArgs {
    return &KeyPairArgs{
        Name:       name,
        SSHKeyPair: ssh.GenerateKeyPair(),
    }
}

func NewKeyPairArgsWithRandomNameAndKey(name *string) *KeyPairArgs {
    randSrc := rand.NewSource(time.Now().UnixNano())
    randomNumber := rand.New(randSrc).Intn(100000)
    keyPairName := fmt.Sprintf("%s%d", *name, randomNumber)

    return &KeyPairArgs{
        Name:         &keyPairName,
        RandomNumber: &randomNumber,
        SSHKeyPair:   ssh.GenerateKeyPair(),
    }
}
