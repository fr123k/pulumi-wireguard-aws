package model

import (
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ComputeResult define the generated properties from compute package
type ComputeResult struct {
    // ID pulumi.IDOutput
    Compute pulumi.CustomResourceState
}

type ProvisionArgs struct {
    ExportName    string
    SourceCompute *ComputeResult
}

type ImageArgs struct {
    Name          string
    Owners        []string
    States        []string
    SourceCompute *ComputeResult
}

// ID return resource id
func (compute ComputeResult) ID() pulumi.IDOutput {
    return compute.Compute.ID()
}

func SelfImage(name string) *ImageArgs {
    return &ImageArgs{
        Name:   name,
        Owners: []string{"self"},
        States: []string{"available"},
    }
}

// ComputeArgs defines the input parameter for the compute resource functions.
type ComputeArgs struct {
    Vpc            *VpcResult
    Security       *SecurityArgs
    SecurityGroups []*SecurityGroup
    IngressRules   []*SecurityRule
    EgressRules    []*SecurityRule
    KeyPair        *KeyPairArgs
    Images         []*ImageArgs
}

func NewComputeArgs(vpc *VpcResult, security *SecurityArgs) *ComputeArgs {
    return NewComputeArgsWithKeyPair(vpc, security, nil)
}

func NewComputeArgsWithSecurityAndKeyPair(security *SecurityArgs, keyPair *KeyPairArgs) *ComputeArgs {
    return NewComputeArgsWithKeyPair(nil, security, keyPair)
}

func NewComputeArgsWithKeyPair(vpc *VpcResult, security *SecurityArgs, keyPair *KeyPairArgs) *ComputeArgs {
    return &ComputeArgs{
        Vpc:      vpc,
        Security: security,
        KeyPair:  keyPair,
    }
}
