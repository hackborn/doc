package doc

// GetRequest provides parameters to a Get operation.
type GetRequest struct {
	Optional
	Condition Expr    // Matching condition to accept a record.
	Fields    *Fields // Limit response to these specific fields.
	Limit     int     // Limit to the number of items
	Flags     GetFlags
}

func (r GetRequest) With(opt any) GetRequest {
	r.Options = append(r.Options, opt)
	return r
}

// GetResponse provides output from a Get operation.
type GetResponse[T any] struct {
	Optional
	Results []*T
}

func (r GetResponse[T]) With(opt any) GetResponse[T] {
	r.Options = append(r.Options, opt)
	return r
}

// GetOneResponse provides output from a GetOne operation.
type GetOneResponse[T any] struct {
	Optional
	Result *T
}

func (r GetOneResponse[T]) With(opt any) GetOneResponse[T] {
	r.Options = append(r.Options, opt)
	return r
}

// SetRequest provides parameters to a Set operation.
type SetRequest[T any] struct {
	Optional

	// The item to set.
	Item T

	// Filter is used to determine which fields will be set.
	// An empty filter will set all fields.
	Filter Filter
}

func (r SetRequest[T]) With(opt any) SetRequest[T] {
	r.Options = append(r.Options, opt)
	return r
}

func (r SetRequest[T]) ItemAny() any {
	return &r.Item
}

func (r SetRequest[T]) GetFilter() Filter {
	return r.Filter
}

// SetResponse provides output from a Set operation.
type SetResponse[T any] struct {
	Optional
	Item *T
}

func (r SetResponse[T]) With(opt any) SetResponse[T] {
	r.Options = append(r.Options, opt)
	return r
}

// DeleteRequest
type DeleteRequest[T any] struct {
	Optional
	Item T
}

func (r DeleteRequest[T]) ItemAny() any {
	return &r.Item
}

// SetRequestAny provides access to the item of a Set request.
type SetRequestAny interface {
	ItemAny() any
	GetFilter() Filter
}

// DeleteRequestAny provides access to the item of a Delete request.
type DeleteRequestAny interface {
	ItemAny() any
}

type DeleteResponse struct {
	Optional
}
