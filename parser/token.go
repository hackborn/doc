package parser

// ------------------------------------------------------------
// TOKEN-T

// tokenT stores a single token type, and associated behaviour.
type tokenT struct {
	Symbol       symbol
	Text         string
	BindingPower int
	nud          nudFn
	led          ledFn
}

// any answers true if my symbol is any of the supplied symbols.
func (t tokenT) any(symbols ...symbol) bool {
	for _, s := range symbols {
		if t.Symbol == s {
			return true
		}
	}
	return false
}

// inside answers true if my symbol is after start and before end.
func (t tokenT) inside(start, end symbol) bool {
	return t.Symbol > start && t.Symbol < end
}

// ------------------------------------------------------------
// FUNC

type nudFn func(*nodeT, *parserT) (*nodeT, error)

type ledFn func(*nodeT, *parserT, *nodeT) (*nodeT, error)

// ------------------------------------------------------------
// TOKEN FUNCS

func emptyNud(n *nodeT, p *parserT) (*nodeT, error) {
	return n, nil
}

func emptyLed(n *nodeT, p *parserT, left *nodeT) (*nodeT, error) {
	return n, nil
}

func binaryLed(n *nodeT, p *parserT, left *nodeT) (*nodeT, error) {
	n.addChild(left)
	right, err := p.Expression(n.Token.BindingPower)
	if err != nil {
		return nil, err
	}
	n.addChild(right)
	return n, nil
}

func enclosedNud(n *nodeT, p *parserT) (*nodeT, error) {
	enclosed, err := p.Expression(n.Token.BindingPower)
	if err != nil {
		return nil, err
	}
	next, err := p.Next()
	if err != nil {
		return nil, err
	}
	if next == nil {
		return nil, newParseError("missing next for " + n.Text)
	}
	if next.Token.Symbol != closeToken {
		return nil, newParseError("missing close for " + n.Text)
	}
	n.addChild(enclosed)
	return n, nil
}
