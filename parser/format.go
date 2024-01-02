package parser

import (
	"fmt"
)

// Format provides the tokens and rules
// used when converting an AST to a string
type Format interface {
	// Translate a keyword to the local format.
	// Example, "AND" might become "&&".
	Keyword(s string) string

	// Convert a value to a string.
	Value(v interface{}) (string, error)
}

func DefaultFormat() Format {
	return _defaultFormat
}

// _format provides default formatting rules.
type _format struct {
	keywords map[string]string
}

func (f *_format) Keyword(s string) string {
	s, _ = f.keywords[s]
	return s
}

func (f *_format) Value(v interface{}) (string, error) {
	return fmt.Sprintf("%v", v), nil
}
