package promql

import (
	"context"
	"fmt"
	"prometheus_lite/pkg/b_promql/parser"
	storage "prometheus_lite/pkg/c_storage"
	"sync"
	"time"
)

type EngineOpts struct {
}

type Engine struct {
	activeQueryTracker QueryTracker
	queryLoggerLock    sync.RWMutex
}

func NewEngine(opts *EngineOpts) *Engine {
	return &Engine{}
}

func (ng *Engine) NewInstantQuery(ctx context.Context, q storage.Queryable, opts QueryOpts, qs string, ts time.Time) (Query, error) {
	pExpr, qry := ng.newQuery(q, qs, opts, ts, ts, 0)
	finishQueue, err := ng.queueActive(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer finishQueue()

	expr, err := parser.ParseExpr(qs)
	if err != nil {
		return nil, err
	}
	*pExpr = PreprocessExpr(expr, ts, ts)

	return qry, nil
}

func (ng *Engine) NewRangeQuery(ctx context.Context, q storage.Queryable, opts QueryOpts, qs string, start, end time.Time, interval time.Duration) (Query, error) {
	//TODO implement me
	panic("implement me")
}

func (ng *Engine) exec(ctx context.Context, q *query) (v parser.Value, err error) {
	return nil, nil
}

func (ng *Engine) queueActive(ctx context.Context, q *query) (func(), error) {
	if ng.activeQueryTracker == nil {
		return func() {}, nil
	}
	queryIndex, err := ng.activeQueryTracker.Insert(ctx, q.q)
	return func() { ng.activeQueryTracker.Delete(queryIndex) }, err
}

func (ng *Engine) newQuery(q storage.Queryable, qs string, opts QueryOpts, start, end time.Time, interval time.Duration) (*parser.Expr, *query) {
	if opts == nil {
		opts = NewPrometheusQueryOpts(false, 0)
	}

	es := &parser.EvalStmt{
		Start:    start,
		End:      end,
		Interval: interval,
	}
	qry := &query{
		q:         qs,
		stmt:      es,
		ng:        ng,
		queryable: q,
	}
	return &es.Expr, qry
}

func PreprocessExpr(expr parser.Expr, start, end time.Time) parser.Expr {
	isStepInvariant := preprocessExprHelper(expr, start, end)
	if isStepInvariant {
		return newStepInvariantExpr(expr)
	}
	return expr
}

// preprocessExprHelper wraps the child nodes of the expression
// with a StepInvariantExpr wherever it's step invariant. The returned boolean is true if the
// passed expression qualifies to be wrapped by StepInvariantExpr.
// It also resolves the preprocessors.
func preprocessExprHelper(expr parser.Expr, start, end time.Time) bool {
	switch n := expr.(type) {
	case *parser.AggregateExpr:
		return preprocessExprHelper(n.Expr, start, end)

	case *parser.BinaryExpr:
		isInvariant1, isInvariant2 := preprocessExprHelper(n.LHS, start, end), preprocessExprHelper(n.RHS, start, end)
		if isInvariant1 && isInvariant2 {
			return true
		}

		if isInvariant1 {
			n.LHS = newStepInvariantExpr(n.LHS)
		}
		if isInvariant2 {
			n.RHS = newStepInvariantExpr(n.RHS)
		}

		return false

	case *parser.Call:
		var isStepInvariant bool
		isStepInvariantSlice := make([]bool, len(n.Args))
		for i := range n.Args {
			isStepInvariantSlice[i] = preprocessExprHelper(n.Args[i], start, end)
			isStepInvariant = isStepInvariant && isStepInvariantSlice[i]
		}

		if isStepInvariant {
			// The function and all arguments are step invariant.
			return true
		}

		for i, isi := range isStepInvariantSlice {
			if isi {
				n.Args[i] = newStepInvariantExpr(n.Args[i])
			}
		}
		return false

	case *parser.UnaryExpr:
		return preprocessExprHelper(n.Expr, start, end)

	case *parser.StringLiteral, *parser.NumberLiteral:
		return true
	}

	panic(fmt.Sprintf("found unexpected node %#v", expr))
}

func newStepInvariantExpr(expr parser.Expr) parser.Expr {
	return &parser.StepInvariantExpr{Expr: expr}
}
