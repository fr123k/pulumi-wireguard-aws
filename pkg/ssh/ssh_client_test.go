package ssh

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "strings"
    "testing"
)

const expectedPublicKeyFromFile = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDZBYRgaR5XFEKS3P5+Wd02jhrHO0LYsqsB0No06Q6anBbC5QDrMMCZoy9Fixoww051mraWQ/vePqyePwd2JpN1CyYIG1nMH2MB3IjGQHy5efsRKH2SY/gjeWaJCIp8DSSpDOmds3ccc7GCGkM608Hg8lUDslhf6VxpNkvpC0/DVVpEzgr0fv6JSK+htdTOrVR6ttqBsu1HKMBmOlkfG9ivf4Sdj/uxFOZhIPnXKQiBVzwouavYS9j9R7EOlax8VZxFrn7a3pj9VhhYpUh+CJs+HNjaPYtLCPGpnwi/94csGJbQzwgMupG/FD5lZ4tco1wcxcPfCUqIdNWVPfXVFARNZEoSfkYJn+ez+iOjzn9a4Iwe+SG5cA1dc5hltBjzSIgWwruKPj9mwJEgluA3owVseInXi4DR1B2IrTK6TyGKKWBEI0YKGjVPKzCF9z+TzIWfxStMVPB16Qx2lVzBkgaqpaljFY+NWM83/T6xFNiDV7kS4a215wLUpJ23qQO5RcmOsNtgp0vka3Sb2qdIqOPI1Z+1BEDyS3sibGKYViJyTy5bnV24BgArXoVE6UfXzl1hPOcD0eWZE3vyEPjlq4f1WO+hbUUrqaEMXjBYJqtpEx3Q5f2iUiBc8dCtbiyLrTmj8mEgmpVhQzgm/pmzSEaHyNwGsm2OY72qidnY46pQxw== franki@MacBook-Pro.local\n"

const expectedPrivateKeyFromFile = `-----BEGIN RSA PRIVATE KEY-----
MIIJKAIBAAKCAgEA2QWEYGkeVxRCktz+flndNo4axztC2LKrAdDaNOkOmpwWwuUA
6zDAmaMvRYsaMMNOdZq2lkP73j6snj8HdiaTdQsmCBtZzB9jAdyIxkB8uXn7ESh9
kmP4I3lmiQiKfA0kqQzpnbN3HHOxghpDOtPB4PJVA7JYX+lcaTZL6QtPw1VaRM4K
9H7+iUivobXUzq1UerbagbLtRyjAZjpZHxvYr3+EnY/7sRTmYSD51ykIgVc8KLmr
2EvY/UexDpWsfFWcRa5+2t6Y/VYYWKVIfgibPhzY2j2LSwjxqZ8Iv/eHLBiW0M8I`

func convertPrivateKeyTostr(privateKey *rsa.PrivateKey) string {
    privkeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
    privkeyPem := pem.EncodeToMemory(
        &pem.Block{
            Type:  "RSA PRIVATE KEY",
            Bytes: privkeyBytes,
        },
    )
    return string(privkeyPem)
}

func TestGenerateKeyPair(t *testing.T) {
    sshKey := GenerateKeyPair()
    if !strings.HasPrefix(*sshKey.PublicKeyStr, "ssh-rsa ") {
        t.Errorf("The publicKeyStr variables is wrong, got: %s, want: %s.", *sshKey.PublicKeyStr, "ssh-rsa ")
    }

    fmt.Printf("private key %s", convertPrivateKeyTostr(sshKey.PrivateKey))
}

func TestReadKeyPair(t *testing.T) {

    sshKey := ReadKeyPair("ssh_testdata/keys/wireguard.pem", "ssh_testdata/keys/wireguard.pem.pub")
    if *sshKey.PublicKeyStr != expectedPublicKeyFromFile {
        t.Errorf("The publicKeyStr variables is wrong, got: %s, want: %s.", *sshKey.PublicKeyStr, expectedPublicKeyFromFile)
    }

    privateKey := convertPrivateKeyTostr(sshKey.PrivateKey)

    if !strings.HasPrefix(privateKey, expectedPrivateKeyFromFile) {
        t.Errorf("The privateKeyStr variables is wrong, got: %s, has to start with: %s.", privateKey, expectedPrivateKeyFromFile)
    }
}
