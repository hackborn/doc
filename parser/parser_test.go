package parser

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------
// TEST-NEW-AST
func TestNewAst(t *testing.T) {
	table := []struct {
		tokens  []any
		want    string
		wantErr error
	}{
		{[]any{"id", "=", "10"}, "id = 10", nil},
		{[]any{"id", "=", "10", "AND", "plan", "=", "9"}, "id = 10 AND plan = 9", nil},
	}
	for i, v := range table {
		ast, haveErr := NewAst(v.tokens...)
		have := ""
		if haveErr == nil {
			var sb strings.Builder
			args := FormatArgs{Writer: &sb, Format: _defaultFormat}
			haveErr = ast.Format(args)
			have = sb.String()
		}

		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestNewAst %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestNewAst %v has no error but exptected %v", i, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestNewAst %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}

// ---------------------------------------------------------
// TEST-VALIDATE
func TestValidate(t *testing.T) {
	table := []struct {
		term    string
		want    string
		wantErr error
	}{
		{"id = 10", "id = 10", nil},
		{"id = tree", "id = tree", nil},
		{"id = \"tree\"", "id = tree", nil},
		{`id = "tree"`, "id = tree", nil},
		{"(id = 10)", "(id = 10)", nil},
		{"id = 10 AND step = 1", "id = 10 AND step = 1", nil},
		{"id = 10 and step = 1", "id = 10 AND step = 1", nil},
		{"id = 10 OR step = 1", "id = 10 OR step = 1", nil},
		{"id = 10 or step = 1", "id = 10 OR step = 1", nil},
		{"id = 10 AND form = \"wd-20\"", "id = 10 AND form = wd-20", nil},
		{"id = 10 and form = wd-20", "", syntaxErr},
		{"id = 10 && form = wd20", "form = wd20", nil}, // XXX This should produce an error but not sure how to identify that
		{"id, form", "id, form", nil},
		{"id,  form", "id, form", nil},
		{"id, form, type", "id, form, type", nil},
	}
	for i, v := range table {
		have, haveErr := validate(v.term, nil)
		//		panic(nil)
		if v.wantErr == nil && haveErr != nil {
			t.Fatalf("TestValidate %v expected no error but has %v", i, haveErr)
		} else if v.wantErr != nil && haveErr == nil {
			t.Fatalf("TestValidate %v has no error but exptected %v", i, v.wantErr)
		} else if have != v.want {
			t.Fatalf("TestValidate %v has \"%v\" but wanted \"%v\"", i, have, v.want)
		}
	}
}
