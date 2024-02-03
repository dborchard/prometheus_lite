package promql

import (
	"context"
	"fmt"
	parser "github.com/dborchard/prometheus_lite/pkg/b_parser"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
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
	// 1. Build Logical Plan Builder
	pExpr, qry := ng.newQuery(q, qs, opts, ts, ts, 0)

	// 2. Check the concurrent queries status
	finishQueue, err := ng.queueActive(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer finishQueue()

	// 3. Parse the Query String
	expr, err := parser.ParseExpr(qs)

	if err != nil {
		return nil, err
	}
	// 4. Build Logical Plan
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
	//defer querier.Close()

	ng.populateSeries(ctxPrepare, querier, s)

	if s.Start == s.End && s.Interval == 0 {
		// i. instant query

		s.Start = s.Start.Add(time.Second * 1) // Just updating with 1 second to make it instant query.
		s.End = s.End.Add(time.Second * 1)
		s.Interval = 1

		evaluator := &evaluator{
			startTimestamp: timeMilliseconds(s.Start),
			endTimestamp:   timeMilliseconds(s.Start), // NOTE: single it is instant query, we keep start=end.
			interval:       1,
			ctx:            ctxPrepare,
			maxSamples:     1000,
		}
		val, _ := evaluator.Eval(s.Expr)

		mat, ok := val.(Matrix)
		if !ok {
			panic(fmt.Errorf("promql.Engine.exec: invalid expression type %q", val.Type()))
		}
		query.matrix = mat

		sort.Sort(mat)

		return mat, nil
	} else {
		// ii. range query
		evaluator := &evaluator{
			startTimestamp: timeMilliseconds(s.Start),
			endTimestamp:   timeMilliseconds(s.End),
			interval:       durationMilliseconds(s.Interval),
			ctx:            ctxPrepare,
			maxSamples:     1000,
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
			interval := s.Interval

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
