package compute

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/model"
    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
    "github.com/pulumi/pulumi/sdk/v3/go/common/resource"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/stretchr/testify/assert"
)

type mocks int

//TODO reduce code lines
//TODO reduce complexity for testing
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
    //TODO implement mock based on invocation count and args
    if args.Token == "aws:ec2/getAmiIds:getAmiIds" {
        outputs["architecture"] = "x86_64"
        outputs["ids"] = []string{"ami-0eb1f3cdeeb8eed2a"}
    }
    return resource.NewPropertyMapFromMap(outputs), nil
}

type ProjectRootFileReader struct {
}

// ReadFile read the file content from a string in the memory instead of the filesystem
func (fileReader ProjectRootFileReader) ReadFile(filename string) ([]byte, error) {
    wd, _ := os.Getwd()
    for !Exists(fmt.Sprintf("%s/%s", wd, filename)) {
        wd = filepath.Dir(wd)
        if len(wd) <= 1 {
            fmt.Println(wd)
            break
        }
    }
    return ioutil.ReadFile(fmt.Sprintf("%s/%s", wd, filename))
}

func Exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func ProjectFileContent() {
    fake := ProjectRootFileReader{}
    model.Util = utility.Util{
        OsReadFile: fake.ReadFile,
    }
}

func TestUserData(t *testing.T) {
    ProjectFileContent()
    _, err := model.Util.ReadFile("cloud-init/wireguard.txt")
    assert.NoError(t, err)
}
