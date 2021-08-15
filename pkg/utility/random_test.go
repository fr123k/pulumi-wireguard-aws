package utility

import (
    "fmt"
    "testing"
)

func TestRandomSecretLength(t *testing.T) {
    secret := RandomSecret(32)

    fmt.Printf("Random secret: %s\n", secret)
    if len(secret) != 32 {
        t.Errorf("Expected 32 random characters")
    }
}

func TestRandomSecretUseingSeed(t *testing.T) {
    secret := RandomSecret(32)

    fmt.Printf("Random secret: %s\n", secret)
    if secret == "(wy)yRNA+a#fL&jK_PV8evlK_yV!T1)y" {
        t.Errorf("Random secrets are not unique (wy)yRNA+a#fL&jK_PV8evlK_yV!T1)y==%s", secret)
    }
}
