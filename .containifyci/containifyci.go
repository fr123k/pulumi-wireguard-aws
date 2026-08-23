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
	wireguard := build.NewGoServiceBuild("pulumi-wireguard")
	wireguard.SourcePackages = []string{}
	wireguard.SourceFiles = []string{}
	wireguard.Verbose = false
	wireguard.Image = ""
	wireguard.File = "cmd/wireguard/hetzner/wireguard.go"
	wireguard.Properties = map[string]*build.ListValue{
		//TODO add a good documentation of possible values (best would build from code)
		// "pulumi": build.NewList("true"), //disable pulumi for now
		"stack": build.NewList("wireguard-hetzner"),
		// "cmd":    build.NewList("up --yes"),
	}
	//TODO: adjust the registry to your own container registry
	wireguard.Registry = "containifyci"

	temporal := build.NewGoServiceBuild("pulumi-temporal")
	temporal.SourcePackages = []string{}
	temporal.SourceFiles = []string{}
	temporal.Verbose = false
	temporal.Image = ""
	temporal.File = "cmd/temporal/hetzner/temporal.go"
	temporal.Properties = map[string]*build.ListValue{
		//TODO add a good documentation of possible values (best would build from code)
		// "pulumi": build.NewList("true"), //disable pulumi for now
		"stack": build.NewList("wireguard-temporal"),
		// "cmd":    build.NewList("up --yes"),
	}
	//TODO: adjust the registry to your own container registry
	temporal.Registry = "containifyci"

	minipc := build.NewGoServiceBuild("pulumi-minipc")
	minipc.SourcePackages = []string{}
	minipc.SourceFiles = []string{}
	minipc.Verbose = false
	minipc.Image = ""
	minipc.File = "cmd/minipc/local/minipc.go"
	minipc.Properties = map[string]*build.ListValue{
		//TODO add a good documentation of possible values (best would build from code)
		// "pulumi": build.NewList("true"), //disable pulumi for now
		"stack": build.NewList("wireguard-minipc"),
		// "cmd":    build.NewList("up --yes"),
	}
	//TODO: adjust the registry to your own container registry
	minipc.Registry = "containifyci"

	build.BuildAsync(wireguard, temporal, minipc)
}
