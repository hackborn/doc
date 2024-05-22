package doc

import (
	"github.com/hackborn/doc/parser"
)

type GetFlags uint64

const (
	GetUnique GetFlags = (1 << iota) // Return only items that have unique values.
)

const (
	// Expose parser-private data for the drivers.
	AndKeyword    = parser.AndKeyword
	AssignKeyword = parser.AssignKeyword
	ListKeyword   = parser.ListKeyword
	OrKeyword     = parser.OrKeyword
)
