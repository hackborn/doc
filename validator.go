package doc

// Validator is used to validate an expression.
type Validator interface {
	// Return true if the field name is valid.
	AcceptField(name string) bool
}
