package doc

func NewFields(names ...string) *Fields {
	return &Fields{names: names}
}

// Fields stores a list of field names.
type Fields struct {
	names []string
}

func (f *Fields) Names() []string {
	return f.names
}
