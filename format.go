package doc

import (
	"github.com/hackborn/doc/parser"
)

// Expose the parser format to consumers of the doc package
// without requiring them to import the parser package.
// Conceptually, parser is completely private.
type Format = parser.Format

// FormatWithDefaults answers the supplied format, adding
// default behaviour for anything the supplied format leaves unhandled.
// f can be nil if you just want the default formatting.
func FormatWithDefaults(f Format) Format {
	fallback := parser.DefaultFormat()
	if f == nil {
		return fallback
	}
	return &compositeFormat{main: f, fallback: fallback}
}

// compositeFormat provides two levels of formatting rules.
type compositeFormat struct {
	main     Format
	fallback Format
}

func (f *compositeFormat) Keyword(s string) string {
	ans := f.main.Keyword(s)
	if ans != "" {
		return ans
	}
	return f.fallback.Keyword(s)
}

func (f *compositeFormat) Value(v interface{}) (string, error) {
	s, err := f.main.Value(v)
	if err == nil {
		return s, err
	}
	return f.fallback.Value(v)
}
