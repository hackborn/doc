package doc

// Filter determines what fields are allowed during an operation.
type Filter struct {
	// Rule provides some predefined rules to determine what fields are
	// allowed. Optional.
	Rule int64
}

// Filter macros

var FilterOff = Filter{}
var FilterSetItem = Filter{Rule: RuleSetItem}
var FilterCreateItem = Filter{Rule: RuleCreateItem}

const (
	RuleOff        = iota // No rule will be applied in the filter
	RuleSetItem           // The item is being set -- all fields are allowed.
	RuleCreateItem        // The item is being created -- all fields are set except auto-generated ones
)
