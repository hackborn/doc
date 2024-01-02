package doc

type Optional struct {
	Options []any
}

func Option[T any](options Optional) (*T, bool) {
	for _, opt := range options.Options {
		switch t := opt.(type) {
		case T:
			return &t, true
		}

	}
	return nil, false
}

type DocOption func(*DocOptions)

type DocOptions struct {
	Local Local
}

func MakeDocOptions(opts ...DocOption) *DocOptions {
	ans := &DocOptions{}
	for _, opt := range opts {
		opt(ans)
	}
	return ans
}

type Local interface {
}

func WithLocal(local Local) DocOption {
	return func(n *DocOptions) {
		n.Local = local
	}
}

type ItemOpt[T any] struct {
	Item *T
}
