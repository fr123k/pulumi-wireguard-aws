// This code only exists to make go.mod happy and go mod tidy to not remove the dependency
// The dependency is defined so that DependaBot can detect new versions of golangci-lint
package main

import (
	_ "github.com/containifyci/dunebot/pkg/version"                // dunebot
	_ "github.com/containifyci/oauth2-storage/pkg/storage"         // oauth2 storage
	_ "github.com/containifyci/secret-operator/pkg/model"          // secret operator
	_ "github.com/containifyci/temporal-worker/pkg/activities/git" // temporal worker
	_ "github.com/temporalio/cli/cmd/temporal"
	_ "github.com/temporalio/ui-server/v2/server/version" // temporal UI
	_ "go.temporal.io/server/common/primitives"           // temporal server
)

func main() {}
