package model

import (
    "fmt"

    "github.com/creasty/defaults"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// VpcArgs define vpc resource
type VpcArgs struct {
    Subnets         []SubnetArgs
    InstanceTenancy string `default:"default"`
    Name            string
    Cidr            string
}

// SubnetArgs define subnet resource
type SubnetArgs struct {
    Cidr string
}

// VpcResult define the generated properties
type VpcResult struct {
    SubnetResults []SubnetResult
    // ID pulumi.IDOutput
    Vpc pulumi.CustomResourceState
}

var VPCArgsDefault = VpcArg("wireguard", "10.8.0.0")

// SubnetResult define the generated properties
type SubnetResult struct {
    Subnet pulumi.CustomResourceState
}

//ID return resource id
func (vpc VpcResult) ID() pulumi.IDOutput {
    return vpc.Vpc.ID()
}

//IDtoInt return ID as int
func (vpc VpcResult) IDtoInt() pulumi.IntOutput {
    return utility.IDtoInt(vpc.Vpc)
}

//ID return resource id
func (subnet SubnetResult) ID() pulumi.IDOutput {
    return subnet.Subnet.ID()
}

func VpcArg(name string, cidr string) *VpcArgs {
    vpcArgs := &VpcArgs{
        Name: name,
        Cidr: fmt.Sprintf("%s/16", cidr),
        Subnets: []SubnetArgs{{
            Cidr: fmt.Sprintf("%s/24", cidr),
        },
        },
    }
    defaults.MustSet(vpcArgs)
    return vpcArgs
}
