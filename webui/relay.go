package webui

import (
	"errors"
	"fmt"
	"os"

	"github.com/fsufitch/fwooshfile/client"
)

// FwooshRelay is a struct representing one relay server
type FwooshRelay struct {
	PublicHost, PrivateHost string
}

// CheckAlive pings the given file relay and returns an error or nil,
// indicating the relay's status.
func (r FwooshRelay) CheckAlive() error {
	client := client.TransferClient{BaseURL: r.PrivateHost}
	return client.CheckAlive()
}

// Relays contains the known relays of this webserver
var Relays []FwooshRelay

// RoundRobinBalancer is a channel serving relays in a round robin manner
var RoundRobinBalancer chan FwooshRelay

var resetRoundRobin chan bool
var roundRobinRunning bool

func roundRobin() error {
	roundRobinRunning = true
	RoundRobinBalancer = make(chan FwooshRelay)
	var err error
	if Relays == nil || len(Relays) == 0 {
		err = new(NoRelaysDefinedError)
	}
	for err == nil {
		for _, relay := range Relays {
			stopLoop := false
			select {
			case RoundRobinBalancer <- relay:
				break
			case <-resetRoundRobin:
				stopLoop = true
			}
			if stopLoop {
				break
			}
		}
		if Relays == nil || len(Relays) == 0 {
			err = new(NoRelaysDefinedError)
		}
	}
	fmt.Fprintf(os.Stderr, "Round robin stopping due to error: %s\n", err)
	roundRobinRunning = false
	return err
}

// NoRelaysDefinedError indicates no relays are defined
type NoRelaysDefinedError string

func (err *NoRelaysDefinedError) Error() string {
	return string(*err)
}

// RestartRoundRobinBalancer does what it says on the lid
func RestartRoundRobinBalancer() error {
	if Relays == nil || len(Relays) == 0 {
		return errors.New("No relays defined! Cannot start round robin!")
	}
	if roundRobinRunning {
		resetRoundRobin <- true
	} else {
		go roundRobin()
	}
	return nil
}
