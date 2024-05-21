package doc

import (
	"fmt"

	"github.com/hackborn/doc/parser"
)

type DB struct {
	driver Driver

	// Cache the driver format, wrapping it in a composite that
	// provides fallback formatting for anything not covered in the driver.
	// Guaranteed non-nil.
	format Format
}

func (d *DB) Close() error {
	driver := d.driver
	d.driver = nil
	if driver != nil {
		return driver.Close()
	}
	return nil
}

// Answer a new expression, using the optional validator and
// the driver's formatting.
func (d *DB) Expr(expr string, v Validator) Expr {
	return &rawExpression{term: expr, v: v, f: d.format}
}

func Open(driverName, dataSourceName string) (*DB, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("doc: unknown driver %v (forgotten import, or select from %v)", driverName, Drivers())
	}
	driver, err := driveri.Open(dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{driver: driver, format: newFormatFor(driver)}, nil
}

func newFormatFor(driver Driver) Format {
	f := driver.Format()
	if f != nil {
		return f
	}
	return parser.DefaultFormat()
}

// Get returns a list of items based on the get condition.
func Get[T any](d *DB, req GetRequest) (GetResponse[T], error) {
	var a trackingAllocator[T]
	opts, err := d.driver.Get(req, &a)
	if err != nil {
		return GetResponse[T]{}, err
	}
	resp := GetResponse[T]{Results: a.All}
	if opts != nil {
		resp.Options = opts.Options
	}
	return resp, nil
}

// GetOne returns a single item based on the get condition.
func GetOne[T any](d *DB, req GetRequest) (GetOneResponse[T], error) {
	req.Limit = 1
	resp, err := Get[T](d, req)
	oneResp := GetOneResponse[T]{}
	if err != nil {
		return oneResp, err
	}
	if len(resp.Results) > 0 {
		oneResp.Result = resp.Results[0]
	}
	return oneResp, nil
}

func Set[T any](d *DB, req SetRequest[T]) (SetResponse[T], error) {
	var a trackingAllocator[T]
	opts, err := d.driver.Set(req, &a)
	if err != nil {
		return SetResponse[T]{}, err
	}

	resp := SetResponse[T]{}
	if len(a.All) > 0 {
		resp.Item = a.All[0]
	}
	if opts != nil {
		resp.Options = opts.Options
	}
	return resp, nil
}

func Delete[T any](d *DB, req DeleteRequest[T]) (DeleteResponse, error) {
	var a trackingAllocator[T]
	var resp DeleteResponse
	opts, err := d.driver.Delete(req, &a)
	if opts != nil {
		resp.Options = opts.Options
	}
	return resp, err
}
