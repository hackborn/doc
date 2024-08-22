package parser

// Used to extract info from an AST.
// Experimental.
type ExtractBinary interface {
	BinaryConjunction(keyword string) error
	BinaryAssignment(lhs string, rhs any) error
}
