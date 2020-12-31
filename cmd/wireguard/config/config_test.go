package config

import (
	"testing"
)

// TestVpcArg with default args
func TestVpcArg(t *testing.T) {
	vpcArg := vpcArg("awesome name", "192.168.0.0")
	
	if vpcArg.Name != "awesome name" {
		t.Errorf("The Name is wrong, got: %s, want: %s.", vpcArg.Name, "awesome name")
	}

	if vpcArg.InstanceTenancy != "default" {
		t.Errorf("The InstanceTenancy is wrong, got: %s, want: %s.", vpcArg.InstanceTenancy, "default")
	}

	if vpcArg.Cidr != "192.168.0.0/16" {
		t.Errorf("The Cidr is wrong, got: %s, want: %s.", vpcArg.Cidr, "192.168.0.0/16")
	}

	for _, subnet := range vpcArg.Subnets {
		if subnet.Cidr != "192.168.0.0/24" {
			t.Errorf("The Subnet cidr is wrong, got: %s, want: %s.", subnet.Cidr, "192.168.0.0/24")
		}
	}
}
