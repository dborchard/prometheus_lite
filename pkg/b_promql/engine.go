package promql

import (
	"context"
	"fmt"
	"prometheus_lite/pkg/b_promql/parser"
	storage "prometheus_lite/pkg/c_storage"
	"sort"
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

	finishQueue, err := ng.queueActive(ctx, q)
	if err != nil {
		return nil, err
	}
	defer finishQueue()

	// Cancel when execution is done or an error was raised.
	defer q.cancel()

	switch s := q.Statement().(type) {
	case *parser.EvalStmt:
		return ng.execEvalStmt(ctx, q, s)
	}

	return nil, nil
}

func (ng *Engine) execEvalStmt(ctx context.Context, query *query, s *parser.EvalStmt) (parser.Value, error) {
	ctxPrepare := ctx
	mint, maxt := FindMinMaxTime(s)

	querier, _ := query.queryable.Querier(mint, maxt)
	defer querier.Close()

	ng.populateSeries(ctxPrepare, querier, s)

	if s.Start == s.End && s.Interval == 0 {
		panic("instant query not implemented")
	} else {
		// range query
		// Range evaluation.
		evaluator := &evaluator{
			startTimestamp: timeMilliseconds(s.Start),
			endTimestamp:   timeMilliseconds(s.End),
			interval:       durationMilliseconds(s.Interval),
			ctx:            ctxPrepare,
		}
		val, _ := evaluator.Eval(s.Expr)

		mat, ok := val.(Matrix)
		if !ok {
			panic(fmt.Errorf("promql.Engine.exec: invalid expression type %q", val.Type()))
		}
		query.matrix = mat

		sort.Sort(mat)

		return mat, nil
	}
}

func (ng *Engine) populateSeries(ctx context.Context, querier storage.Querier, s *parser.EvalStmt) {
	// Whenever a MatrixSelector is evaluated, evalRange is set to the corresponding range.
	// The evaluation of the VectorSelector inside then evaluates the given range and unsets
	// the variable.
	var evalRange time.Duration

	parser.Inspect(s.Expr, func(node parser.Node, path []parser.Node) error {
		switch n := node.(type) {
		case *parser.VectorSelector:
			start, end := getTimeRangesForSelector(s, n, path, evalRange)
			interval := ng.getLastSubqueryInterval(path)
			if interval == 0 {
				interval = s.Interval
			}
			hints := &storage.SelectHints{
				Start: start,
				End:   end,
				Step:  durationMilliseconds(interval),
				Range: durationMilliseconds(evalRange),
				//Func:  extractFuncFromPath(path),
			}
			evalRange = 0
			hints.By, hints.Grouping = extractGroupsFromPath(path)
			n.UnexpandedSeriesSet = querier.Select(ctx, false, hints, n.LabelMatchers...)

		case *parser.MatrixSelector:
			evalRange = n.Range
		}
		return nil
	})
}

func FindMinMaxTime(s *parser.EvalStmt) (int64, int64) {
	return 0, 0
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
