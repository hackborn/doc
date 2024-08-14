package parser

import (
	"cmp"
	"fmt"
	"io"
	"strconv"

	"github.com/hackborn/onefunc/errors"
	"github.com/hackborn/onefunc/strings"
)

// ------------------------------------------------------------
// AST-NODE

// AstNode defines the basic interface for evaluating expressions.
type AstNode interface {
	Format(args FormatArgs) error
	Fields(args *FieldArgs) error
	Extract(any) error
}

type FormatArgs struct {
	Writer io.StringWriter
	Format Format
	Ctx    FormatContext
	Eb     errors.Block
}

type FieldArgs struct {
	Fields []string
	Ctx    FormatContext
}

// ------------------------------------------------------------
// BINARY-NODE

// binaryNode performs binary operations on the current interface{}.
type binaryNode struct {
	Op      symbol
	Keyword string
	Lhs     AstNode
	Rhs     AstNode
}

func (n *binaryNode) Format(args FormatArgs) error {
	if err := n.stateErr(); err != nil {
		return err
	}
	keyword := args.Format.Keyword(n.Keyword)
	if keyword == "" {
		return newSyntaxError("format returned empty for keyword \"" + n.Keyword + "\"")
	}

	lhW := strings.GetWriter(args.Eb)
	defer strings.PutWriter(lhW)
	rhW := strings.GetWriter(args.Eb)
	defer strings.PutWriter(rhW)

	newArgs := args
	newArgs.Writer = lhW
	// LHS
	newArgs.Ctx = n.lhsContext(args.Ctx)
	err := n.Lhs.Format(newArgs)
	if err != nil {
		return err
	}
	// RHS
	newArgs.Writer = rhW
	newArgs.Ctx = n.rhsContext(args.Ctx)
	err = n.Rhs.Format(newArgs)
	if err != nil {
		return err
	}

	args.Writer.WriteString(strings.String(lhW))
	args.Writer.WriteString(keyword)
	args.Writer.WriteString(strings.String(rhW))

	return cmp.Or(strings.StringErr(lhW), strings.StringErr(rhW))
}

func (n *binaryNode) Fields(args *FieldArgs) error {
	if err := n.stateErr(); err != nil {
		return err
	}

	prevctx := args.Ctx
	args.Ctx = n.lhsContext(prevctx)
	err := n.Lhs.Fields(args)
	args.Ctx = n.rhsContext(prevctx)
	return cmp.Or(err, n.Rhs.Fields(args))
}

func (n *binaryNode) Extract(_fn any) error {
	fn, ok := _fn.(ExtractBinary)
	if !ok {
		return nil
	}
	if err := n.stateErr(); err != nil {
		return err
	}
	switch n.Keyword {
	case AndKeyword, OrKeyword, ListKeyword:
		err := fn.BinaryConjunction(n.Keyword)
		err = cmp.Or(err, n.Lhs.Extract(_fn), n.Rhs.Extract(_fn))
		if err != nil {
			return err
		}
	case AssignKeyword:
		if lhs, rhs, err := n.valueStrings(); err == nil {
			err := fn.BinaryAssignment(lhs, rhs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *binaryNode) stateErr() error {
	if n.Lhs == nil || n.Rhs == nil {
		return newMalformedError("binary node")
	}
	if n.Keyword == "" {
		return newMalformedError("binary node missing keyword for symbol " + strconv.Itoa(int(n.Op)))
	}
	return nil
}

func (n *binaryNode) lhsContext(ctx FormatContext) FormatContext {
	return NoFormatContext
}

func (n *binaryNode) rhsContext(ctx FormatContext) FormatContext {
	switch n.Keyword {
	case AssignKeyword:
		return ValueContext
	default:
		return NoFormatContext
	}
}

func (n *binaryNode) valueStrings() (string, string, error) {
	lhn, ok1 := n.Lhs.(*valueNode)
	rhn, ok2 := n.Rhs.(*valueNode)
	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("Missing value node")
	}
	lhs, ok1 := lhn.Value.(string)
	rhs, ok2 := rhn.Value.(string)
	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("Missing value node string")
	}
	return lhs, rhs, nil
}

// ------------------------------------------------------------
// VALUE-NODE

// valueNode returns a constant value (string, float, etc.).
type valueNode struct {
	Value interface{}
}

func (n *valueNode) Format(args FormatArgs) error {
	s := ""
	var err error
	if args.Ctx == ValueContext {
		s, err = args.Format.Value(n.Value)
		if err != nil {
			return err
		}
	} else {
		s = fmt.Sprintf("%v", n.Value)
	}
	_, err = args.Writer.WriteString(s)
	return err
}

func (n *valueNode) Fields(args *FieldArgs) error {
	if args.Ctx == ValueContext {
		return nil
	}
	switch t := n.Value.(type) {
	case string:
		args.Fields = append(args.Fields, t)
	}
	return nil
}

func (n *valueNode) Extract(any) error {
	return nil
}

// ------------------------------------------------------------
// UNARY-NODE

// unaryNode performs a unary operation on the current interface{}.
type unaryNode struct {
	Op    symbol
	Child AstNode
}

func (n *unaryNode) Format(args FormatArgs) error {
	if err := n.stateErr(); err != nil {
		return err
	}

	switch n.Op {
	case openToken:
		args.Writer.WriteString("(")
		n.Child.Format(args)
		args.Writer.WriteString(")")
	default:
		return newUnhandledError("unary " + strconv.Itoa(int(n.Op)))
	}
	return nil
}

func (n *unaryNode) Fields(args *FieldArgs) error {
	if err := n.stateErr(); err != nil {
		return err
	}

	return n.Child.Fields(args)
}

func (n *unaryNode) Extract(any) error {
	return nil
}

func (n *unaryNode) stateErr() error {
	if n.Child == nil {
		return newMalformedError("unary node missing child")
	}
	return nil
}

// ------------------------------------------------------------
// CREATION

// NewAst answers a new AST based on the supplied tokens.
// Tokens will be interpreted as keywords if they exist. This bypasses
// the parsing stage so it can only be used to build simple ASTs.
func NewAst(tokens ...any) (AstNode, error) {
	// Wrap everyone in the final AST nodes
	newAstWrap(tokens)
	// Deal with not having any binding power or doing any actual
	// processing by going through multiple levels.
	err := newAstBinary(AssignKeyword, tokens)
	err = cmp.Or(err, newAstBinary("", tokens))
	if err != nil {
		return nil, err
	}
	return newAstFinalize(tokens)
}

var (
	newAstWrapSymbol = map[string]symbol{
		AndKeyword:    andToken,
		AssignKeyword: assignToken,
		ListKeyword:   listToken,
		OrKeyword:     orToken,
	}
)

func newAstWrap(tokens []any) {
	for i, t := range tokens {
		switch v := t.(type) {
		case string:
			if tok, ok := newAstWrapSymbol[v]; ok {
				tokens[i] = &binaryNode{Op: tok, Keyword: v}
			} else {
				tokens[i] = &valueNode{Value: t}
			}
		}
	}
}

func newAstBinary(keyword string, tokens []any) error {
	var last AstNode = nil
	lasti := -1
	for i, t := range tokens {
		if t == nil {
			continue
		}

		switch v := t.(type) {
		case *binaryNode:
			if keyword == "" || v.Keyword == keyword {
				if v.Lhs == nil {
					if last == nil {
						return fmt.Errorf("invalid syntax: missing lhs")
					}
					v.Lhs = last
					tokens[lasti] = nil
				}
				if v.Rhs == nil {
					v.Rhs = newAstBinaryNext(i+1, tokens)
					if v.Rhs == nil {
						return fmt.Errorf("invalid syntax: missing rhs")
					}
				}
			}
		}
		last = t.(AstNode)
		lasti = i
	}
	return nil
}

func newAstBinaryNext(_i int, tokens []any) AstNode {
	for i := _i; i < len(tokens); i++ {
		if tokens[i] != nil {
			ans := tokens[i]
			tokens[i] = nil
			return ans.(AstNode)
		}
	}
	return nil
}

func newAstFinalize(tokens []any) (AstNode, error) {
	var a AstNode = nil
	for _, t := range tokens {
		if t != nil {
			if a != nil {
				return nil, fmt.Errorf("invalid syntax")
			}
			a = t.(AstNode)
		}
	}
	if a == nil {
		return nil, fmt.Errorf("invalid syntax")
	}
	return a, nil
}
