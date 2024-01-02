package parser

var (
	badRequestErr = newBadRequestError("")
	conditionErr  = newConditionError("")
	evalErr       = newEvalError("")
	syntaxErr     = newSyntaxError("")
	malformedErr  = newMalformedError("")
	mismatchErr   = newMismatchError("")
	parseErr      = newParseError("")
	unhandledErr  = newUnhandledError("")
)

// --------------------------------
// SQI-ERROR

func newBadRequestError(msg string) error {
	return &docErr{badRequestErrCode, msg, nil}
}

func newConditionError(msg string) error {
	return &docErr{conditionErrCode, msg, nil}
}

func newEvalError(msg string) error {
	return &docErr{evalErrCode, msg, nil}
}

func newSyntaxError(msg string) error {
	return &docErr{syntaxErrCode, msg, nil}
}

func newMalformedError(msg string) error {
	return &docErr{malformedErrCode, msg, nil}
}

func newMismatchError(msg string) error {
	return &docErr{mismatchErrCode, msg, nil}
}

func newParseError(msg string) error {
	return &docErr{parseErrCode, msg, nil}
}

func newUnhandledError(msg string) error {
	return &docErr{unhandledErrCode, msg, nil}
}

type docErr struct {
	code int
	msg  string
	err  error
}

func (e *docErr) ErrorCode() int {
	return e.code
}

func (e *docErr) Error() string {
	var label string
	switch e.code {
	case badRequestErrCode:
		label = "doc: bad request"
	case conditionErrCode:
		label = "doc: condition"
	case evalErrCode:
		label = "doc: eval"
	case syntaxErrCode:
		label = "doc: invalid syntax"
	case malformedErrCode:
		label = "doc: malformed"
	case mismatchErrCode:
		label = "doc: mismatch"
	case parseErrCode:
		label = "doc: parse"
	case unhandledErrCode:
		label = "doc: unhandled"
	default:
		label = "doc: error"
	}
	if e.msg != "" {
		label += " (" + e.msg + ")"
	}
	if e.err != nil {
		label += " (" + e.err.Error() + ")"
	}
	return label
}

// --------------------------------
// CONST and VAR

const (
	badRequestErrCode = 1000 + iota
	conditionErrCode
	evalErrCode
	syntaxErrCode
	malformedErrCode
	mismatchErrCode
	parseErrCode
	unhandledErrCode
)
