package actors

import (
	"fmt"
)

type connector struct {
	actions chan<- func()
	quit    chan struct{}
}

// Connector this interface define the basic functions for all the actor Connector implementations.
type Connector interface {
	Connect(address string) string
	Stop()
}

func newConnector() connector {
	actions := make(chan func())
	c := connector{
		actions: actions,
		quit:    make(chan struct{}),
	}
	go c.loop(actions)
	return c
}

func (c *connector) loop(actions <-chan func()) {
	for {
		select {
		case f := <-actions:
			f()
		case <-c.quit:
			return
		}
	}
}

// register the result of a mesh.Router.NewGossip.
func (c *connector) Connect(address string) string {
	resultChan := make(chan string, 0)
	c.actions <- func() {
		fmt.Printf("dummy connection established to %s", address)
		resultChan <- address
	}
	return <-resultChan
}

func (c *connector) Stop() {
	close(c.quit)
}
