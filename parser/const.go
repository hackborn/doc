package parser

const (
	AndKeyword    = "AND"
	AssignKeyword = "="
	ListKeyword   = ","
	OrKeyword     = "OR"
)

type symbol int

const (
	// Special tokens
	illegalToken symbol = iota

	// Raw values
	intToken    // 12345
	floatToken  // 123.45
	stringToken // "abc"

	// Assignment
	assignToken // =

	// Negation
	negToken // -

	// List. A binary that combines left and right.
	listToken // ,

	// Comparison
	startComparison

	eqlToken // ==
	neqToken // !=

	endComparison

	// -- CONDITIONALS. All conditional operators must be after this
	startConditional

	andToken // AND, and
	orToken  // OR, or

	// -- END CONDITIONALS.
	endConditional

	// -- UNARIES. All unary operators must be after this
	startUnary

	// Enclosures
	openToken  // (
	closeToken // ) // All closes must be after the opens

	// -- END UNARIES.
	endUnary
)

type FormatContext int

const (
	NoFormatContext FormatContext = iota
	ValueContext                  // A value, i.e. the RHS of an assignment
)

var (
	tokenMap = map[symbol]*tokenT{
		illegalToken: &tokenT{illegalToken, "", 0, emptyNud, emptyLed},
		intToken:     &tokenT{intToken, "", 10, emptyNud, emptyLed},
		floatToken:   &tokenT{floatToken, "", 10, emptyNud, emptyLed},
		stringToken:  &tokenT{stringToken, "", 10, emptyNud, emptyLed},
		assignToken:  &tokenT{assignToken, AssignKeyword, 80, emptyNud, binaryLed},
		negToken:     &tokenT{negToken, "", 10, emptyNud, emptyLed},
		eqlToken:     &tokenT{eqlToken, "==", 70, emptyNud, binaryLed},
		neqToken:     &tokenT{neqToken, "!=", 70, emptyNud, binaryLed},
		listToken:    &tokenT{listToken, ListKeyword, 60, emptyNud, binaryLed},
		andToken:     &tokenT{andToken, AndKeyword, 60, emptyNud, binaryLed},
		orToken:      &tokenT{orToken, OrKeyword, 60, emptyNud, binaryLed},
		openToken:    &tokenT{openToken, "(", 0, enclosedNud, emptyLed},
		closeToken:   &tokenT{closeToken, ")", 0, emptyNud, emptyLed},
	}
	keywordMap = map[string]*tokenT{
		AssignKeyword: tokenMap[assignToken],
		`-`:           tokenMap[negToken],
		`==`:          tokenMap[eqlToken],
		`!=`:          tokenMap[neqToken],
		ListKeyword:   tokenMap[listToken],
		AndKeyword:    tokenMap[andToken],
		OrKeyword:     tokenMap[orToken],
		`(`:           tokenMap[openToken],
		`)`:           tokenMap[closeToken],
	}
)

var (
	_defaultFormat = &_format{keywords: map[string]string{
		AndKeyword:    ` ` + AndKeyword + ` `,
		AssignKeyword: ` ` + AssignKeyword + ` `,
		ListKeyword:   ListKeyword + ` `,
		OrKeyword:     ` ` + OrKeyword + ` `,
	}}
)
