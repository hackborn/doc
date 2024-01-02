package parser

import (
	_ "fmt"
	"strings"
)

// Parse converts an expression string into an AST.
func Parse(term string) (AstNode, error) {
	tokens, err := scan(term)
	if err != nil {
		return nil, err
	}
	p := newParser(tokens)
	tree, err := p.Expression(0)
	if err != nil {
		return nil, err
	}
	if tree == nil {
		return nil, newParseError("no result")
	}
	ast, err := tree.asAst()
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func validate(term string, f Format) (string, error) {
	ast, err := Parse(term)
	if err != nil {
		return "", err
	}
	if ast == nil {
		return "", newEvalError("missing AST")
	}
	var sb strings.Builder
	if f == nil {
		f = _defaultFormat
	}
	args := FormatArgs{Writer: &sb, Format: f}
	err = ast.Format(args)
	return sb.String(), err
}

// ------------------------------------------------------------
// PARSER-T

type parser interface {
	// Next answers the next node, or nil if we're finished.
	// Note that a finished condition is both the node and error being nil;
	// any error response is always an actual error.
	Next() (*nodeT, error)
	// Peek the next value. Note that this is never nil; an illegal is returned if we're at the end.
	Peek() *nodeT
	Expression(rbp int) (*nodeT, error)
}

type parserT struct {
	tokens   []*nodeT
	position int
	illegal  *nodeT
}

func newParser(tokens []*nodeT) parser {
	illegal := &nodeT{Token: tokenMap[illegalToken]}
	return &parserT{tokens: tokens, position: 0, illegal: illegal}
}

func (p *parserT) Next() (*nodeT, error) {
	if p.position >= len(p.tokens) {
		return nil, nil
	}
	pos := p.position
	p.position++
	return p.tokens[pos], nil
}

func (p *parserT) Peek() *nodeT {
	if p.position >= len(p.tokens) {
		return p.illegal
	}
	return p.tokens[p.position]
}

func (p *parserT) Expression(rbp int) (*nodeT, error) {
	n, err := p.Next()
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, newParseError("premature stop")
	}
	//	fmt.Println("Expression on rbp", rbp, "next \"", n.Text, "\"", n.Token)
	left, err := n.Token.nud(n, p)
	//	fmt.Println("\tat", n.Text, "left", left, "err", err)
	if err != nil {
		return nil, err
	}
	//	fmt.Println("rbp binding", rbp, "peek binding", p.Peek().Token.BindingPower, "token", p.Peek().Token.Text, p.Peek().Token.Symbol)

	for rbp < p.Peek().Token.BindingPower {
		n, err = p.Next()
		//		fmt.Println("\tloop rbp", rbp, "next \"", n.Text, "\"", n.Token)
		if err != nil {
			return nil, err
		}
		if n == nil {
			return nil, newParseError("premature stop")
		}
		left, err = n.Token.led(n, p, left)
		if err != nil {
			return nil, err
		}
	}
	//	fmt.Println("returning left", left.Text, left.Token.Text)
	return left, nil
}
