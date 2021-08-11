package model

import (
    "os"
    "strings"
    "testing"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
)

const userDataContent = `
MYV4IP=$(curl {{ TEST_METADATA_URL }})

PublicKey = {{ TEST_CLIENT_PUBLICKEY }}
PersistentKeepalive = 25
EOF`

const userDataContentExpected = `
MYV4IP=$(curl TEST_METADATA_URL)

PublicKey = TEST_CLIENT_PUBLICKEY
PersistentKeepalive = 25
EOF`

const userDataContentExpectedEmptyVariables = `
MYV4IP=$(curl )

PublicKey = 
PersistentKeepalive = 25
EOF`

const userDataContentEnvironmentVariablesExpected = `
MYV4IP=$(curl ENV_TEST_METADATA_URL)

PublicKey = ENV_CLIENT_PUBLICKEY
PersistentKeepalive = 25
EOF`

func userDataVariables() map[string]string {
    return map[string]string{
        "{{ TEST_CLIENT_PUBLICKEY }}": "TEST_CLIENT_PUBLICKEY",
        "{{ TEST_METADATA_URL }}":     "TEST_METADATA_URL",
    }
}

func assert(t *testing.T, err error) {
    if err != nil {
        t.Error(err.Error())
    }
}

func assertError(t *testing.T, err error, want string) {
    if err == nil {
        t.Errorf("Expected error was nil, got: nil, want: %s", want)
    }
    if !strings.Contains(err.Error(), want) {
        t.Errorf("The error message is wrong, got: %s, want: %s.", err.Error(), want)
    }
}

//TestTemplateVariablesString test the TemplateVariablesString method
func TestTemplateVariablesString(t *testing.T) {

    variables := TemplateVariablesString(userDataVariables())

    if len(variables) != 2 {
        t.Errorf("The userData variables size is wrong, got: %d, want: %d.", len(variables), 2)
    }

    for _, variable := range variables {
        if variable.Type != STRING {
            t.Errorf("The variable '%s' type is wrong, got: %d, want: %d.", variable.Key, variable.Type, STRING)
        }
    }
}

//TestTemplateVariablesEnvironment test the TemplateVariablesEnvironment method
func TestTemplateVariablesEnvironment(t *testing.T) {
    variables := TemplateVariablesEnvironment(userDataVariables())

    if len(variables) != 2 {
        t.Errorf("The userData variables size is wrong, got: %d, want: %d.", len(variables), 2)
    }

    for _, variable := range variables {
        if variable.Type != ENVIRONMENT {
            t.Errorf("The variable '%s' type is wrong, got: %d, want: %d.", variable.Key, variable.Type, ENVIRONMENT)
        }
    }
}

//TestNewUserDataWithContent test the NewUserDataWithContent method
func TestNewUserDataWithContent(t *testing.T) {
    userData, err := NewUserDataWithContent(userDataContent, TemplateVariablesString(userDataVariables()))

    assert(t, err)

    if userData.OriginalContent != userDataContent {
        t.Errorf("The userData original content is wrong, got: %s, want: %s.", userData.OriginalContent, userDataContent)
    }

    if userData.Content != userDataContentExpected {
        t.Errorf("The userData content is wrong, got: %s, want: %s.", userData.Content, userDataContentExpected)
    }
}

//TestNewUserDataWithContentAndUndefineEnvironmentVariables test the NewUserDataWithContent method
func TestNewUserDataWithContentAndUndefineEnvironmentVariables(t *testing.T) {
    userData, err := NewUserDataWithContent(userDataContent, TemplateVariablesEnvironment(userDataVariables()))

    assert(t, err)

    if userData.OriginalContent != userDataContent {
        t.Errorf("The userData original content is wrong, got: %s, want: %s.", userData.OriginalContent, userDataContent)
    }

    if userData.Content != userDataContentExpectedEmptyVariables {
        t.Errorf("The userData content is wrong, got: %s, want: %s.", userData.Content, userDataContentExpectedEmptyVariables)
    }
}

func TestNewUserDataWithContentAndEnvironmentVariables(t *testing.T) {
    os.Setenv("TEST_CLIENT_PUBLICKEY", "ENV_CLIENT_PUBLICKEY")
    os.Setenv("TEST_METADATA_URL", "ENV_TEST_METADATA_URL")
    userData, err := NewUserDataWithContent(userDataContent, TemplateVariablesEnvironment(userDataVariables()))

    assert(t, err)

    if userData.OriginalContent != userDataContent {
        t.Errorf("The userData original content is wrong, got: %s, want: %s.", userData.OriginalContent, userDataContent)
    }

    if userData.Content != userDataContentEnvironmentVariablesExpected {
        t.Errorf("The userData content is wrong, got: %s, want: %s.", userData.Content, userDataContentEnvironmentVariablesExpected)
    }
}

//TestNewUserDataWithContentNoVariables test the NewUserDataWithContent method with no variables
func TestNewUserDataWithContentNoVariables(t *testing.T) {
    userData, err := NewUserDataWithContentNoVariables(userDataContent)

    assert(t, err)

    if userData.OriginalContent != userDataContent {
        t.Errorf("The userData original content is wrong, got: %s, want: %s.", userData.OriginalContent, userDataContent)
    }

    if userData.Content != userDataContent {
        t.Errorf("The userData content is wrong, got: %s, want: %s.", userData.Content, userDataContent)
    }
}

func MockFileContent(content string) {
    fake := utility.InMemoryFileReader{Str: content}
    Util = utility.Util{
        OsReadFile: fake.ReadFile,
    }
}

//TestNewUserDataNoVariables test the NewUserDataNoVariables method with no variables
func TestNewUserDataNoVariables(t *testing.T) {

    MockFileContent(userDataContent)

    userData, err := NewUserDataNoVariables("nonexistingfile")

    assert(t, err)

    if userData.OriginalContent != userDataContent {
        t.Errorf("The userData original content is wrong, got: %s, want: %s.", userData.OriginalContent, userDataContent)
    }

    if userData.Content != userDataContent {
        t.Errorf("The userData content is wrong, got: %s, want: %s.", userData.Content, userDataContent)
    }
}

//TestNewUserDataWithEnvironmentVariables test the method NewUserData with environment variables
func TestNewUserDataWithEnvironmentVariables(t *testing.T) {
    os.Setenv("TEST_CLIENT_PUBLICKEY", "ENV_CLIENT_PUBLICKEY")
    os.Setenv("TEST_METADATA_URL", "ENV_TEST_METADATA_URL")

    MockFileContent(userDataContent)

    userData, err := NewUserData("nonexistingfile", TemplateVariablesEnvironment(userDataVariables()))

    assert(t, err)

    if userData.OriginalContent != userDataContent {
        t.Errorf("The userData original content is wrong, got: %s, want: %s.", userData.OriginalContent, userDataContent)
    }

    if userData.Content != userDataContentEnvironmentVariablesExpected {
        t.Errorf("The userData content is wrong, got: %s, want: %s.", userData.Content, userDataContent)
    }
}

//test error handling

//TestNewUserDataWithEnvironmentVariables test the method NewUserData with environment variables
func TestNewUserDataWithNonExistingFile(t *testing.T) {
    Util = utility.NewUtil()
    _, err := NewUserData("nonexistingfile.txt", TemplateVariablesEnvironment(userDataVariables()))

    assertError(t, err, "open nonexistingfile.txt: no such file or directory")
}
