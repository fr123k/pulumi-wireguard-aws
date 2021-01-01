package model

import "github.com/pulumi/pulumi/sdk/v2/go/pulumi"

// ComputeResult define the generated properties from compute package
type ComputeResult struct {
	// ID pulumi.IDOutput
	Compute pulumi.CustomResourceState
}

type ImageArgs struct {
	Name string
	SourceCompute *ComputeResult
}

//ID return resource id
func (compute ComputeResult) ID() pulumi.IDOutput {
	return compute.Compute.ID()
}
