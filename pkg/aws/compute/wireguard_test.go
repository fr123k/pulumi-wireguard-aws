package compute

import (
    "sync"
    "testing"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/stretchr/testify/assert"
)

type testFnc = func(t *testing.T, wg *sync.WaitGroup, infra *infrastructure)

func Setup(t *testing.T, infraArgs InfrastructureArgsFnc, test testFnc) {
    err := pulumi.RunErr(func(ctx *pulumi.Context) error {
        var wg sync.WaitGroup
        ProjectFileContent()
        computeArgs, err := infraArgs(ctx)
        if err != nil {
            return err
        }

        infra, err := CreateServer(ctx, computeArgs)
        assert.NoError(t, err)
        test(t, &wg, infra)

        wg.Wait()
        return nil
    }, pulumi.WithMocks("project", "stack", mocks(0)))
    assert.NoError(t, err)
}

func TestWireguardWithVPNEnabled(t *testing.T) {
    Setup(t, DefaultComputeArgs ,func(t *testing.T, wg *sync.WaitGroup, infra *infrastructure) {
        assert.Equal(t, "ami-0eb1f3cdeeb8eed2a", *infra.imageID)

        assertTags(t, infra, wg, "Name", assert.Containsf)
        assertTags(t, infra, wg, "JobUrl", assert.Containsf)
        // Test if the instance is configured with user_data.
        assertUserDataNotNil(t, infra, wg)

        // Test if port 22 for ssh is not exposed to public.
        assertSSHPort(t, infra, wg, false)
    })
}

func TestWireguardWithVPNDisabled(t *testing.T) {
    Setup(t, DefaultComputeArgs2, func(t *testing.T, wg *sync.WaitGroup, infra *infrastructure) {
        assert.Equal(t, "ami-0eb1f3cdeeb8eed2a", *infra.imageID)

        assertTags(t, infra, wg, "Name", assert.Containsf)
        assertTags(t, infra, wg, "JobUrl", assert.Containsf)
        // Test if the instance is configured with user_data.
        assertUserDataNotNil(t, infra, wg)

        // Test if port 22 for ssh is exposed to public.
        assertSSHPort(t, infra, wg, true)
    })
}
