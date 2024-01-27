package promql

import (
	"context"
	"prometheus_lite/pkg/b_promql/parser"
	storage "prometheus_lite/pkg/c_storage"
)

// A Query is derived from an a raw query string and can be run against an engine
// it is associated with.
type Query interface {
	// Exec processes the query. Can only be called once.
	Exec(ctx context.Context) *Result
	// Close recovers memory used by the query result.
	Close()
	// Statement returns the parsed statement of the query.
	Statement() parser.Statement
	// Cancel signals that a running query execution should be aborted.
	Cancel()
	// String returns the original query string.
	String() string
}

var _ Query = (*query)(nil)

type query struct {
	queryable storage.Queryable
	q         string
	stmt      parser.Statement
	matrix    Matrix
	cancel    func()
	ng        *Engine
}

func (q *query) Statement() parser.Statement {
	return q.stmt
}

func (q *query) String() string {
	return q.q
}

func (q *query) Cancel() {
	if q.cancel != nil {
		q.cancel()
	}
}

func (q *query) Close() {
	for _, s := range q.matrix {
		print(s.Floats)
		//putFPointSlice(s.Floats)
		//putHPointSlice(s.Histograms)
	}
}

func (q *query) Exec(ctx context.Context) *Result {
	// Exec query.
	res, err := q.ng.exec(ctx, q)
	return &Result{Err: err, Value: res}
}
