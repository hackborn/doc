package doc

import (
	"reflect"
)

type Allocator interface {
	// New answers a new object
	New() any
	// TypeName answers the typename of the object type I create.
	TypeName() string
}

type trackingAllocator[T any] struct {
	All []*T
}

func (a *trackingAllocator[T]) New() any {
	t := new(T)
	a.All = append(a.All, t)
	return t
}

func (a *trackingAllocator[T]) TypeName() string {
	var t T
	return reflect.TypeOf(t).Name()
}
