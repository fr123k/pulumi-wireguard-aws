package actors

import (
	"fmt"
)

// SSHConnector the ssh implementation of the Connector actor
type SDKActor struct {
	connector
	execute func(string) string
}

// type execute func() string

// NewSSHConnector initialize an ssh connector
func NewSDKActor(execute func(string) string) SDKActor {
	sdkActor := SDKActor{}
	sdkActor.connector = newConnector()
	sdkActor.execute = execute
	return sdkActor
}

// Connect this function is called when the virtual instance is created and can recevie connection.
func (c *SDKActor) Connect(address string) string {
	resultChan := make(chan string, 0)
	c.actions <- func() {
		fmt.Printf("Can open connection to %s", address)
		
		result := c.execute(address)

		resultChan <- result
	}
	return <-resultChan
}
