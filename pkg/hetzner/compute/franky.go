package compute

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/model"
	"github.com/fr123k/pulumi-wireguard-aws/pkg/shared"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateFrankyVM creates an SBX VM on Hetzner Cloud.
// If the image is a pre-baked snapshot (numeric ID or "prebaked" prefix), it uses the minimal
// cloud-init script. Otherwise, it uses the full cloud-init script for base Ubuntu images.
func CreateFrankyVM(ctx *pulumi.Context, computeArgs *model.ComputeArgs, vmIP string) (*model.ComputeResult, error) {
	userData, err := shared.FrankyUserData()
	if err != nil {
		return nil, err
	}

	computeArgs.UserData = userData

	infra, err := CreateServer(ctx, computeArgs, vmIP, exports)
	if err != nil {
		return nil, err
	}

	return &model.ComputeResult{
		Compute: infra.Server.CustomResourceState,
	}, nil
}
