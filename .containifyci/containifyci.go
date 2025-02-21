//go:generate bash -c "if [ ! -f go.mod ]; then echo 'Initializing go.mod...'; go mod init .containifyci; else echo 'go.mod already exists. Skipping initialization.'; fi"
//go:generate go get github.com/containifyci/engine-ci/protos2
//go:generate go get github.com/containifyci/engine-ci/client
//go:generate go mod tidy

package main

import (
	"os"

	"github.com/containifyci/engine-ci/client/pkg/build"
)

func main() {
	os.Chdir("../")
	opts := build.NewGoServiceBuild("pulumi-wireguard")
	opts.SourcePackages = []string{}
	opts.SourceFiles = []string{}
	opts.Verbose = false
	opts.File = "cmd/wireguard/hetzner/wireguard.go"
	opts.Properties = map[string]*build.ListValue{
		//TODO add a good documentation of possible values (best would build from code)
		"pulumi": build.NewList("true"),
		"stack":  build.NewList("wireguard-hetzner"),
		// "cmd":    build.NewList("up --yes"),
	}
	//TODO: adjust the registry to your own container registry
	opts.Registry = "containifyci"
	build.Serve(opts)
}
