package doc

import (
	_ "fmt"
	"strings"

	"github.com/hackborn/doc/parser"
)

type Expr interface {
	// Compile the expression. No-op for already-compiled expressions.
	// Expressions that are not compiled may need to be compiled whenever
	// you ask for any info.
	Compile() (Expr, error)

	// Format answers the expression as a formatted string.
	Format() (string, error)
}

// NewExpr answers a new compiled expression based on the supplied tokens.
// Tokens will be interpreted as keywords if they exist. This bypasses
// the parsing stage -- it's an optimization for knowledgeable clients,
// but also can only handle simple phrases, so in general clients should
// prefer passing in an expression string.
func NewExpr(f Format, tokens ...any) (Expr, error) {
	ast, err := parser.NewAst(tokens...)
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	if f == nil {
		f = parser.DefaultFormat()
	}
	args := parser.FormatArgs{Writer: &sb, Format: f}
	err = ast.Format(args)
	if err != nil {
		return nil, err
	}
	return &compiledExpr{ast: ast, formatted: sb.String()}, nil
}

// compiledExpr is a parsed, cached expression.
type compiledExpr struct {
	ast parser.AstNode

	// Formatted is the formatted string produced by the parsed expression.
	formatted string

	fields []string
}

func (e *compiledExpr) Compile() (Expr, error) {
	return e, nil
}

func (e *compiledExpr) Format() (string, error) {
	return e.formatted, nil
}

// rawExpression contains a raw expression term and the information
// necessary to compile it.
type rawExpression struct {
	term string
	v    Validator
	f    Format
}

func (e *rawExpression) Compile() (Expr, error) {
	ast, err := parser.Parse(e.term)
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	args := parser.FormatArgs{Writer: &sb, Format: e.f}
	err = ast.Format(args)
	if err != nil {
		return nil, err
	}
	fa := &parser.FieldArgs{}
	err = ast.Fields(fa)
	if err != nil {
		return nil, err
	}
	return &compiledExpr{ast: ast, formatted: sb.String(), fields: fa.Fields}, nil
}

func (e *rawExpression) Format() (string, error) {
	expr, _ := e.Compile()
	if expr == nil {
		return "", nil
	}
	return expr.Format()
}
