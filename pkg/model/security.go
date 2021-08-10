package model

import (
	"github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
)

// SecurityArgs type that define security attributes
type SecurityArgs struct {
	VPNEnabledSSH bool
	VPNCidr	   string
}

type SecurityRule struct {
	Protocol		string
	SourcePort		int
	DestinationPort	int
	CidrBlocks		[]string
}

// NewSecurityArgs initialize a SecurityArgs type
func NewSecurityArgs(vpnEnabledSSH bool, vpnCidr string) *SecurityArgs {
	return &SecurityArgs{
		VPNEnabledSSH:	vpnEnabledSSH,
		VPNCidr:		vpnCidr,
	}
}

// NewSecurityArgsForVPC initialize a SecurityArgs type
func NewSecurityArgsForVPC(vpnEnabledSSH bool, vpc *VpcArgs) *SecurityArgs {
	//TODO fix the hard coded Subned Index support multiple subnets or even an subset of subnets.
	return NewSecurityArgs(vpnEnabledSSH, vpc.Subnets[0].Cidr)
}

// Println prints the struct as json to stdout
func (security SecurityArgs) Println() { utility.Println(security) }
