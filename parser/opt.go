package parser

// ------------------------------------------------------------
// OPT

// Opt contains options for evaluation.
type Opt struct {
	// Strict causes sloppy conditions to become errors. For example, comparing a
	// number to a string is false if strict is off, but error if it's on.
	Strict bool
	// OnError is a value returned when one of the typed Eval() statements returns an error.
	// Must match the type. For example, the value must be assigend a string if using EvalString().
	OnError interface{}
}

func (o Opt) onErrorBool() bool {
	if v, ok := o.OnError.(bool); ok {
		return v
	}
	return false
}

func (o Opt) onErrorFloat64() float64 {
	if v, ok := o.OnError.(float64); ok {
		return v
	}
	return 0.0
}

func (o Opt) onErrorInt() int {
	if v, ok := o.OnError.(int); ok {
		return v
	}
	return 0
}

func (o Opt) onErrorString() string {
	if v, ok := o.OnError.(string); ok {
		return v
	}
	return ""
}

func (o Opt) onErrorStringInterfaceMap() map[string]interface{} {
	if v, ok := o.OnError.(map[string]interface{}); ok {
		return v
	}
	return nil
}

func (o Opt) onErrorStringStringMap() map[string]string {
	if v, ok := o.OnError.(map[string]string); ok {
		return v
	}
	return nil
}

func (o Opt) onErrorStringSlice() []string {
	if v, ok := o.OnError.([]string); ok {
		return v
	}
	return nil
}
