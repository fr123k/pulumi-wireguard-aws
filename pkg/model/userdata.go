package model

import (
    "fmt"
    "os"
    "strings"

    "github.com/fr123k/pulumi-wireguard-aws/pkg/utility"
)

//Util define the Util type to use can be overwritten as happen in the unit tests.
var Util = utility.NewUtil()

//UserData define userdata input for virtual machine resources
type UserData struct {
    OriginalContent string
    Content         string
    FileName        string
    Variables       []TemplateVariable
}

//TemplateVariable define a variable and its type
type TemplateVariable struct {
    Type  TemplateVariableType
    Key   string
    Value string
}

//TemplateVariableType define the type of an variable. Possible values are STRING and ENVIRONMENT
type TemplateVariableType int32

const (
    //STRING the variable is just a string and its key and value are used without modification.
    STRING TemplateVariableType = 0
    //ENVIRONMENT the variable value reference a environment variable value.
    ENVIRONMENT TemplateVariableType = 1
    //ENV_PROPERTY the variable value reference a environment variable value.
    ENV_PROPERTY TemplateVariableType = 2
    //STRING_PROPERTY the variable is just a string and its key and value are used without modification.
    STRING_PROPERTY TemplateVariableType = 3
)

//TemplateVariables converts a list map of variables to a TemplateVariable array with the passed variablesType
func TemplateVariables(variables map[string]string, variablesType TemplateVariableType) []TemplateVariable {
    templateVariablesIdx := 0
    templateVariables := make([]TemplateVariable, len(variables))
    for key, value := range variables {
        // fmt.Println("Key:", key, "Value:", value)
        templateVariables[templateVariablesIdx] = TemplateVariable{
            Type:  variablesType,
            Key:   key,
            Value: value,
        }
        templateVariablesIdx++
    }
    return templateVariables
}

//TemplateVariablesString converts a list map of variables to a TemplateVariable array with the variablesType STRING
func TemplateVariablesString(variables map[string]string) []TemplateVariable {
    return TemplateVariables(variables, STRING)
}

//TemplateVariablesStringProperty converts a list map of variables to a TemplateVariable array with the variablesType STRING_PROPERTY
func TemplateVariablesStringProperty(variables map[string]string) []TemplateVariable {
    return TemplateVariables(variables, STRING_PROPERTY)
}

//TemplateVariablesEnvironment converts a list map of variables to a TemplateVariable array with the variablesType ENVIRONMENT
func TemplateVariablesEnvironment(variables map[string]string) []TemplateVariable {
    return TemplateVariables(variables, ENVIRONMENT)
}

//TemplateVariablesEnvProperty converts a list map of variables to a TemplateVariable array with the variablesType ENV_PROPERTY
func TemplateVariablesEnvProperty(variables map[string]string) []TemplateVariable {
    return TemplateVariables(variables, ENV_PROPERTY)
}

//NewUserDataWithContentNoVariables return a UserData type fully initialized
func NewUserDataWithContentNoVariables(origContent string) (*UserData, error) {
    emptyTemplateVariables := make([]TemplateVariable, 0)
    return NewUserDataWithContent(origContent, emptyTemplateVariables)
}

//NewUserDataWithContent return a UserData type fully initialized
func NewUserDataWithContent(origContent string, variables []TemplateVariable) (*UserData, error) {
    content := renderTemplate(origContent, variables)

    return &UserData{
        OriginalContent: origContent,
        Content:         content,
        Variables:       variables,
    }, nil
}

//NewUserDataNoVariables return a UserData type fully initialized
func NewUserDataNoVariables(fileName string) (*UserData, error) {
    emptyTemplateVariables := make([]TemplateVariable, 0)
    return NewUserData(fileName, emptyTemplateVariables)
}

//NewUserData return a UserData type fully initialized
func NewUserData(fileName string, variables []TemplateVariable) (*UserData, error) {
    fileContent, err := Util.ReadFile(fileName) // just pass the file name
    if err != nil {
        return nil, err
    }

    userData, err := NewUserDataWithContent(*fileContent, variables)
    userData.FileName = fileName
    return userData, err
}

func renderTemplate(template string, variables []TemplateVariable) string {
    result := template
    for _, variable := range variables {
        fmt.Println("Key:", variable.Key, "Value:", variable.Value)
        keyReplace := fmt.Sprintf("{{ %s }}", variable.Key)
        if variable.Type == STRING {

            result = strings.ReplaceAll(result, keyReplace, variable.Value)
        } else if variable.Type == STRING_PROPERTY {

            valueReplace := fmt.Sprintf("%s=%s", variable.Key, variable.Value)
            result = strings.ReplaceAll(result, keyReplace, valueReplace)
        } else if variable.Type == ENVIRONMENT {

            envVariable, ok2 := os.LookupEnv(variable.Value)
            if ok2 {
                result = strings.ReplaceAll(result, keyReplace, envVariable)
            } else {
                result = strings.ReplaceAll(result, keyReplace, "")
            }
        } else if variable.Type == ENV_PROPERTY {
            envVariable, ok2 := os.LookupEnv(variable.Value)
            valueReplace := fmt.Sprintf("%s=%s", variable.Key, envVariable)
            if ok2 {
                result = strings.ReplaceAll(result, keyReplace, valueReplace)
            } else {
                result = strings.ReplaceAll(result, keyReplace, "")
            }
        }
    }

    return result
}
