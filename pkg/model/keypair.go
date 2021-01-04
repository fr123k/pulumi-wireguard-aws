package model

import (
	"fmt"
	"math/rand"
	"time"
)

type KeyPairArgs struct {
	Name *string
	RandomNumber *int
}

func NewKeyPairArgs(name *string) *KeyPairArgs {
	return &KeyPairArgs{Name: name}
}

func NewKeyPairArgsWithRandomName(name *string) *KeyPairArgs {
	randSrc := rand.NewSource(time.Now().UnixNano())
	randomNumber := rand.New(randSrc).Intn(100000)
	keyPairName := fmt.Sprintf("%s%d", *name, randomNumber)
	return &KeyPairArgs{
		Name: &keyPairName,
		RandomNumber: &randomNumber,
	}
}
