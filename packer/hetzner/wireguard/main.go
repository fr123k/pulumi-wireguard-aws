// This code only exists to make go.mod happy and go mod tidy to not remove the dependency.
// The dependencies are defined so that DependaBot can detect new versions of
// wireguard-ui and secret-operator.
package main

import (
	_ "github.com/ngoduykhanh/wireguard-ui/store"   // wireguard-ui
	_ "github.com/containifyci/secret-operator/pkg/model" // secret operator
)

func main() {}