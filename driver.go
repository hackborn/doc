package doc

import (
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

type Driver interface {
	Open(dataSourceName string) (Driver, error)
	Close() error

	// Format provides driver-specific formatting rules
	// when converting an expression to a string. Return
	// nil for the default rules.
	Format() Format

	// Get performs a query based on the request parameters.
	// The responses are collected in the allocator, and the
	// driver can return any additional data in the first return param.
	Get(req GetRequest, a Allocator) (*Optional, error)

	Set(req SetRequestAny, a Allocator) (*Optional, error)

	Delete(req DeleteRequestAny, a Allocator) (*Optional, error)
}

type PrivateDriver interface {
	Private(any) error
}

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("doc: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("doc: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	list := make([]string, 0, len(drivers))
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}
