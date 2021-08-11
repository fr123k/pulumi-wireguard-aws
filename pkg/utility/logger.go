package utility

import (
    "encoding/json"
    "fmt"

    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Logger struct {
    Ctx *pulumi.Context
}

func (logger Logger) Info(msg string, args ...interface{}) {
    logger.Ctx.Log.Info(fmt.Sprintf(msg, args...), nil)
}

func (logger Logger) Error(msg string, args ...interface{}) {
    logger.Ctx.Log.Error(fmt.Sprintf(msg, args...), nil)
}

func (logger Logger) Debug(msg string, args ...interface{}) {
    logger.Ctx.Log.Debug(fmt.Sprintf(msg, args...), nil)
}

// Println prints any object as json to stdout
func Println(object interface{}) {
    out, err := json.Marshal(object)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(out))
}
