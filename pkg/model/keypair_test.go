package model

import (
	"fmt"
	"testing"
)

// TestNewKeyPair test the NewKeyPair method
func TestNewKeyPair(t *testing.T) {

	keyPairName := "keypairname"
	keyPair := NewKeyPairArgs(&keyPairName)

	if *keyPair.Name != keyPairName {
		t.Errorf("The Name variables of the KeyPair is wrong, got: %s, want: %s.", *keyPair.Name, keyPairName)
	}
}

// TestKeyPairArgsWithRandomName test the KeyPairArgsWithRandomName method
func TestKeyPairArgsWithRandomName(t *testing.T) {

	keyPairName := "keypairname-"
	keyPair := NewKeyPairArgsWithRandomName(&keyPairName)

	expectedKeyPairName := fmt.Sprintf("%s%d", keyPairName, *keyPair.RandomNumber)
	if *keyPair.Name != expectedKeyPairName {
		t.Errorf("The Name variables of the KeyPair is wrong, got: %s, want: %s.", *keyPair.Name, expectedKeyPairName)
	}
}
