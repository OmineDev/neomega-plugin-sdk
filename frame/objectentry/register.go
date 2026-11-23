// Package objectentry connects a trusted Go object's init function to the host.
package objectentry

import (
	"fmt"
	"sync"

	"github.com/OmineDev/neomega-plugin-sdk/frame"
)

// Factory creates a component instance for one configuration.
type Factory func() frame.Component

var registration struct {
	sync.Mutex
	active  bool
	factory Factory
	count   int
}

var loading sync.Mutex

// Register registers the component factory from the plugin package's init function. Call it once per
// archive. Outside a dynamic load it has no effect. The plugin must call the host's copy of this function.
//
//go:noinline
func Register(factory Factory) {
	registration.Lock()
	defer registration.Unlock()
	if !registration.active {
		return
	}
	registration.factory = factory
	registration.count++
}

// DuringLoad runs load and returns the factory registered during that call. It serializes loads and
// rejects a missing, nil or duplicate registration.
func DuringLoad(load func() error) (Factory, error) {
	loading.Lock()
	defer loading.Unlock()
	registration.Lock()
	registration.active = true
	registration.Unlock()
	defer func() {
		registration.Lock()
		registration.active = false
		registration.factory = nil
		registration.count = 0
		registration.Unlock()
	}()
	if err := load(); err != nil {
		return nil, err
	}
	registration.Lock()
	defer registration.Unlock()
	if registration.count != 1 || registration.factory == nil {
		return nil, fmt.Errorf("expected one objectentry.Register call, got %d", registration.count)
	}
	return registration.factory, nil
}
