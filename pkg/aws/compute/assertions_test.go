package compute

import (
	"sync"
	"testing"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type assertion = func(t assert.TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) bool

func assertTags(t *testing.T, infra *Infrastructure, wg *sync.WaitGroup, expected interface{}, assertFnc assertion) {
	wg.Add(1)
	pulumi.All(infra.Server.URN(), infra.Server.Tags).ApplyT(func(all []interface{}) error {
		urn := all[0].(pulumi.URN)
		tags := all[1].(map[string]string)

		assertFnc(t, tags, expected, "missing a Name tag on server %v", urn)
		wg.Done()
		return nil
	})
}

func assertUserDataNotNil(t *testing.T, infra *Infrastructure, wg *sync.WaitGroup) {
	wg.Add(1)
	pulumi.All(infra.Server.URN(), infra.Server.UserData).ApplyT(func(all []interface{}) error {
		urn := all[0].(pulumi.URN)
		userData := all[1].(string)

		assert.NotNil(t, userData, "expect userData set on server on server %v", urn)
		assert.NotEmpty(t, userData, "expect userData set on server on server %v", urn)

		wg.Done()
		return nil
	})
}

func assertSSHPort(t *testing.T, infra *Infrastructure, wg *sync.WaitGroup, public bool) {
	wg.Add(1)
	pulumi.All(infra.Groups[0].URN(), infra.Groups[0].Ingress, infra.Groups[1].Ingress).ApplyT(func(all []interface{}) error {
		urn := all[0].(pulumi.URN)

		ingress := append(all[1].([]ec2.SecurityGroupIngress), all[2].([]ec2.SecurityGroupIngress)...)

		assert.Len(t, ingress, 4, "expect 4 ingress security rules set on server")

		for _, i := range ingress {
			openToInternet := false
			if i.ToPort == 22 {
				for _, b := range i.CidrBlocks {
					if b == "0.0.0.0/0" {
						openToInternet = true
						break
					}
				}
			}

			if i.FromPort != 22 {
				continue
			}
			assert.Equal(t, "tcp", i.Protocol, "Expect protocol 'tcp' for ssh rule on group %v", urn)
			assert.Falsef(t, openToInternet && !public, "illegal SSH port 22 open to the Internet (CIDR 0.0.0.0/0) on group %v", urn)

		}

		wg.Done()
		return nil
	})
}
