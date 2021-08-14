package compute

import (
	"sync"
	"testing"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type assertion = func(t assert.TestingT, s interface{}, contains interface{}, msg string, args ...interface{}) bool

func assertTags(t *testing.T, infra *infrastructure, wg *sync.WaitGroup, expected interface{}, assertFnc assertion) {
    wg.Add(1)
    pulumi.All(infra.server.URN(), infra.server.Tags).ApplyT(func(all []interface{}) error {
        urn := all[0].(pulumi.URN)
        tags := all[1].(map[string]string)

        // assert.Containsf(t, tags, "Name", "missing a Name tag on server %v", urn)
        assertFnc(t, tags, expected, "missing a Name tag on server %v", urn)
        wg.Done()
        return nil
    })
}

func assertUserDataNotNil(t *testing.T, infra *infrastructure, wg *sync.WaitGroup) {
    wg.Add(1)
    pulumi.All(infra.server.URN(), infra.server.UserData).ApplyT(func(all []interface{}) error {
        urn := all[0].(pulumi.URN)
        userData := all[1].(string)

        assert.NotNil(t, userData, "expect userData set on server on server %v", urn)
        wg.Done()
        return nil
    })
}

func assertSSHPort(t *testing.T, infra *infrastructure, wg *sync.WaitGroup, public bool) {
    wg.Add(1)
    pulumi.All(infra.groups[0].URN(), infra.groups[0].Ingress, infra.groups[1].Ingress).ApplyT(func(all []interface{}) error {
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

            assert.Falsef(t, i.FromPort == 22 && (openToInternet && !public), "illegal SSH port 22 open to the Internet (CIDR 0.0.0.0/0) on group %v", urn)
        }

        wg.Done()
        return nil
    })
}
