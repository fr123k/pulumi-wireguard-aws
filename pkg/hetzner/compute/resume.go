package compute

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateResumeVM creates a resume (GDPR API) VM on Hetzner Cloud.
// The server is a public HTTP/HTTPS server (nginx + certbot) with SSH
// restricted to the wireguard peer IP via iptables.
func CreateResumeVM(ctx *pulumi.Context, computeArgs *model.ComputeArgs, vmIP string) (*model.ComputeResult, error) {
	userData, err := shared.ResumeUserData()
	if err != nil {
		return nil, err
	}

	computeArgs.UserData = userData

	infra, err := CreateServer(ctx, computeArgs, vmIP, exports)
	if err != nil {
		return nil, err
	}

	return &model.ComputeResult{
		Compute: &infra.Server.CustomResourceState,
	}, nil
}