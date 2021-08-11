package utility

import (
    "encoding/json"
    "fmt"
)

// Println prints any object as json to stdout
func Println(object interface{}) {
    out, err := json.Marshal(object)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(out))
}
