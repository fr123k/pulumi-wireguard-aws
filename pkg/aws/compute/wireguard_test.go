package compute

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "testing"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/aws/network"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
    "github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
    "github.com/pulumi/pulumi/sdk/v3/go/common/resource"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/stretchr/testify/assert"
)

//TODO reduce code lines
//TODO reduce complexity for testing

type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
    outputs := args.Inputs.Mappable()
    fmt.Printf("Mock Called %s\n", args.TypeToken)
    if args.TypeToken == "aws:ec2/instance:Instance" {
        outputs["publicIp"] = "203.0.113.12"
        outputs["publicDns"] = "ec2-203-0-113-12.compute-1.amazonaws.com"
    }
    return args.Name + "_id", resource.NewPropertyMapFromMap(outputs), nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
    outputs := map[string]interface{}{}
    fmt.Printf("Mock Called %s\n", args.Token)
    if args.Token == "aws:ec2/getAmiIds:getAmiIds" {
        outputs["architecture"] = "x86_64"
        outputs["ids"] = []string{"ami-0eb1f3cdeeb8eed2a"}
    }
    return resource.NewPropertyMapFromMap(outputs), nil
}

func DefaultComputeArgs(ctx *pulumi.Context) (*model.ComputeArgs, error) {
    // cfg := config.New(ctx, "")
    security := model.NewSecurityArgsForVPC(true, model.VPCArgsDefault)
    security.Println()

    vpc, err := network.CreateVPC(ctx, model.VPCArgsDefault)
    if err != nil {
        return nil, err
    }

    keyPairName := "wireguard-"
    keyPair := model.NewKeyPairArgsWithRandomNameAndKey(&keyPairName)
    computeArgs := model.NewComputeArgsWithKeyPair(vpc, security, keyPair)

    computeArgs.Images = []*model.ImageArgs{
        {
            Name:   "wireguard-ami",
            Owners: []string{"self"},
            States: []string{"available"},
        },
        {
            Name:   "ubuntu/images/hvm-ssd/ubuntu-*-18.04-amd64-server-*",
            Owners: []string{"099720109477"},
        },
    }

    tags := map[string]string{
        "JobUrl":         "travis_job_url",
        "Project":        "wireguard",
        "pulumi-managed": "True",
    }

    externalSecurityGroup := model.SecurityGroup{
        Name:        "wireguard-external",
        Description: "Pulumi Managed. Allow Wireguard client traffic from internet.",
        Tags:        tags,
        IngressRules: []*model.SecurityRule{
            model.AllowOnePortRule("udp", 51820),
            model.AllowSSHRule(security),
        },
        EgressRules: []*model.SecurityRule{
            model.AllowAllRule(),
        },
    }
    //The order is important the referenced security groups has to be first.
    computeArgs.SecurityGroups = []*model.SecurityGroup{
        &externalSecurityGroup,
        {
            Name:        "wireguard-admin",
            Description: "Pulumi Managed. Allow admin traffic internal resources from VPN",
            Tags:        tags,
            IngressRules: []*model.SecurityRule{
                model.AllowAllRuleSecGroup(&externalSecurityGroup),
                model.AllowICMPRule(&externalSecurityGroup),
            },
            EgressRules: []*model.SecurityRule{
                model.AllowAllRule(),
            },
        },
    }
    return computeArgs, nil
}

// InMemoryFileReader define type for reading files from memory instead of the filesystem
type ProjectRootFileReader struct {
}

// ReadFile read the file content from a string in the memory instead of the filesystem
func (fileReader ProjectRootFileReader) ReadFile(filename string) ([]byte, error) {
    wd, _ := os.Getwd()
    for !strings.HasSuffix(wd, "pulumi-wireguard-aws") {
        wd = filepath.Dir(wd)
    }
    return ioutil.ReadFile(fmt.Sprintf("%s/%s", wd, filename))
}

func ProjectFileContent() {
    fake := ProjectRootFileReader{}
    model.Util = utility.Util{
        OsReadFile: fake.ReadFile,
    }
}

func TestUserData(t *testing.T) {
    ProjectFileContent()
    _, err := model.Util.ReadFile("cloud-init/user-data.txt")
    assert.NoError(t, err)
}

func TestInfrastructure(t *testing.T) {
    err := pulumi.RunErr(func(ctx *pulumi.Context) error {
        ProjectFileContent()
        computeArgs, err := DefaultComputeArgs(ctx)
        if err != nil {
            return err
        }

        infra, err := createWireguardVM(ctx, computeArgs)
        assert.NoError(t, err)

        var wg sync.WaitGroup
        wg.Add(3)

        pulumi.All(infra.server.URN(), infra.server.Tags).ApplyT(func(all []interface{}) error {
            urn := all[0].(pulumi.URN)
            tags := all[1].(map[string]string)

            assert.Containsf(t, tags, "Name", "missing a Name tag on server %v", urn)
            wg.Done()
            return nil
        })

        // Test if the instance is configured with user_data.
        pulumi.All(infra.server.URN(), infra.server.UserData).ApplyT(func(all []interface{}) error {
            urn := all[0].(pulumi.URN)
            userData := all[1].(string)

            assert.NotNil(t, userData, "expect userData set on server on server %v", urn)
            wg.Done()
            return nil
        })

        // Test if port 22 for ssh is exposed.
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

                assert.Falsef(t, i.FromPort == 22 && openToInternet, "illegal SSH port 22 open to the Internet (CIDR 0.0.0.0/0) on group %v", urn)
            }

            wg.Done()
            return nil
        })

        wg.Wait()
        return nil
    }, pulumi.WithMocks("project", "stack", mocks(0)))
    assert.NoError(t, err)
}
