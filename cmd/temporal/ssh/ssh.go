package main

import (
    "fmt"
    _ "io"
    _ "io/ioutil"
    _ "os"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/ssh"
)

//TODO pass private key file
//TODO pass full qualified hostname or ip
//TODO pass command to execute and parse result
func main() {

    sshClient := ssh.SSHClientConfig{
        Hostname:   "ec2-34-242-223-163.eu-west-1.compute.amazonaws.com",
        Port:       22,
        Username:   "ubuntu",
        SSHKeyPair: *ssh.ReadPrivateKey("/Users/franki/private/github/pulumi-wireguard-aws/keys/wireguard.pem"),
    }

    result, err := sshClient.SSHCommand("sudo cloud-init status --wait")
    if err != nil {
        panic(fmt.Errorf("Failed to create session: %s", err))
    }
    fmt.Printf("Result: %s", *result)
}
