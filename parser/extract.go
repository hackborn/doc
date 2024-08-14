package parser

// Used to extract info from an AST.
// Experimental.
type ExtractBinary interface {
	BinaryConjunction(keyword string)
	BinaryAssignment(lhs, rhs string)
}
