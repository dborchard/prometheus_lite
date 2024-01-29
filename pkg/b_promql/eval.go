package promql

import (
	"context"
	"fmt"
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser"
	"github.com/dborchard/prometheus_lite/pkg/z_model/histogram"
	"github.com/dborchard/prometheus_lite/pkg/z_model/labels"
	"slices"
)

// An evaluator evaluates the given expressions over the given fixed
// timestamps. It is attached to an engine through which it connects to a
// querier and reports errors. On timeout or cancellation of its context it
// terminates.
type evaluator struct {
	ctx context.Context

	startTimestamp int64 // Start time in milliseconds.
	endTimestamp   int64 // End time in milliseconds.
	interval       int64 // Interval in milliseconds.
}

func (ev *evaluator) Eval(expr parser.Expr) (v parser.Value, err error) {
	defer ev.recover(expr, &err)

	v = ev.eval(expr)
	return v, nil
}

func (ev *evaluator) eval(expr parser.Expr) parser.Value {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		// Grouping labels must be sorted (expected both by generateGroupingKey() and aggregation()).
		sortedGrouping := e.Grouping
		slices.Sort(sortedGrouping)

		initSeries := func(series labels.Labels, h *EvalSeriesHelper) {
			h.groupingKey = generateGroupingKey(series, sortedGrouping)
		}

		unwrapParenExpr(&e.Param)
		param := unwrapStepInvariantExpr(e.Param)
		unwrapParenExpr(&param)
		if s, ok := param.(*parser.StringLiteral); ok {
			return ev.rangeEval(initSeries, func(v []parser.Value, sh [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
				return ev.aggregation(e, sortedGrouping, s.Val, v[0].(Vector), sh[0], enh)
			}, e.Expr)
		}

		return ev.rangeEval(initSeries, func(v []parser.Value, sh [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
			var param float64
			if e.Param != nil {
				param = v[0].(Vector)[0].F
			}
			return ev.aggregation(e, sortedGrouping, param, v[1].(Vector), sh[1], enh)
		}, e.Param, e.Expr)

	case *parser.UnaryExpr:
		val := ev.eval(e.Expr)
		mat := val.(Matrix)
		return mat

	case *parser.BinaryExpr:
		switch lt, rt := e.LHS.Type(), e.RHS.Type(); {
		case lt == parser.ValueTypeScalar && rt == parser.ValueTypeScalar:
			return ev.rangeEval(nil, func(v []parser.Value, _ [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
				val := scalarBinop(e.Op, v[0].(Vector)[0].F, v[1].(Vector)[0].F)
				return append(enh.Out, Sample{F: val})
			}, e.LHS, e.RHS)
		}

	case *parser.NumberLiteral:
		return ev.rangeEval(nil, func(v []parser.Value, _ [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
			return append(enh.Out, Sample{F: e.Val, Metric: labels.EmptyLabels()})
		})

	case *parser.StringLiteral:
		return String{V: e.Val, T: ev.startTimestamp}

	//case *parser.VectorSelector:
	//	ws, err := checkAndExpandSeriesSet(ev.ctx, e)
	//	if err != nil {
	//		ev.error(errWithWarnings{fmt.Errorf("expanding series: %w", err), ws})
	//	}
	//	mat := make(Matrix, 0, len(e.Series))
	//	it := storage.NewMemoizedEmptyIterator(durationMilliseconds(ev.lookbackDelta))
	//	var chkIter chunkenc.Iterator
	//	for i, s := range e.Series {
	//		chkIter = s.Iterator(chkIter)
	//		it.Reset(chkIter)
	//		ss := Series{
	//			Metric: e.Series[i].Labels(),
	//		}
	//
	//		for ts, step := ev.startTimestamp, -1; ts <= ev.endTimestamp; ts += ev.interval {
	//			step++
	//			_, f, h, ok := ev.vectorSelectorSingle(it, e, ts)
	//			if ok {
	//				if ev.currentSamples < ev.maxSamples {
	//					if h == nil {
	//						if ss.Floats == nil {
	//							ss.Floats = getFPointSlice(numSteps)
	//						}
	//						ss.Floats = append(ss.Floats, FPoint{F: f, T: ts})
	//						ev.currentSamples++
	//						ev.samplesStats.IncrementSamplesAtStep(step, 1)
	//					} else {
	//						if ss.Histograms == nil {
	//							ss.Histograms = getHPointSlice(numSteps)
	//						}
	//						point := HPoint{H: h, T: ts}
	//						ss.Histograms = append(ss.Histograms, point)
	//						histSize := point.size()
	//						ev.currentSamples += histSize
	//						ev.samplesStats.IncrementSamplesAtStep(step, int64(histSize))
	//					}
	//				}
	//			}
	//		}
	//
	//		if len(ss.Floats)+len(ss.Histograms) > 0 {
	//			mat = append(mat, ss)
	//		}
	//	}
	//	return mat
	//
	//case *parser.MatrixSelector:
	//	if ev.startTimestamp != ev.endTimestamp {
	//		panic(errors.New("cannot do range evaluation of matrix selector"))
	//	}
	//	return ev.matrixSelector(e)

	case *parser.StepInvariantExpr:
		switch ce := e.Expr.(type) {
		case *parser.StringLiteral, *parser.NumberLiteral:
			return ev.eval(ce)
		}

		newEv := &evaluator{
			startTimestamp: ev.startTimestamp,
			endTimestamp:   ev.startTimestamp, // Always a single evaluation.
			interval:       ev.interval,
			ctx:            ev.ctx,
		}
		res := newEv.eval(e.Expr)
		for ts, step := ev.startTimestamp, -1; ts <= ev.endTimestamp; ts += ev.interval {
			step++
		}
		switch e.Expr.(type) {
		case *parser.MatrixSelector:
			// We do not duplicate results for range selectors since result is a matrix
			// with their unique timestamps which does not depend on the step.
			return res
		}

		// For every evaluation while the value remains same, the timestamp for that
		// value would change for different eval times. Hence we duplicate the result
		// with changed timestamps.
		mat, ok := res.(Matrix)
		if !ok {
			panic(fmt.Errorf("unexpected result in StepInvariantExpr evaluation: %T", expr))
		}
		for i := range mat {
			if len(mat[i].Floats)+len(mat[i].Histograms) != 1 {
				panic(fmt.Errorf("unexpected number of samples"))
			}
			for ts := ev.startTimestamp + ev.interval; ts <= ev.endTimestamp; ts += ev.interval {
				if len(mat[i].Floats) > 0 {
					mat[i].Floats = append(mat[i].Floats, FPoint{
						T: ts,
						F: mat[i].Floats[0].F,
					})
				} else {
					point := HPoint{
						T: ts,
						H: mat[i].Histograms[0].H,
					}
					mat[i].Histograms = append(mat[i].Histograms, point)

				}
			}
		}
		return res
	}
	panic(fmt.Errorf("unexpected expression type %T", expr))
}

// grouping labels.
func generateGroupingKey(metric labels.Labels, grouping []string) uint64 {
	//if without {
	//	return metric.HashWithoutLabels(buf, grouping...)
	//}
	return 0
	//return metric.HashForLabels(buf, grouping...)
}

// unwrapParenExpr does the AST equivalent of removing parentheses around a expression.
func unwrapParenExpr(e *parser.Expr) {
	//for {
	//	if p, ok := (*e).(*parser.ParenExpr); ok {
	//		*e = p.Expr
	//	} else {
	//		break
	//	}
	//}
}
func unwrapStepInvariantExpr(e parser.Expr) parser.Expr {
	if p, ok := e.(*parser.StepInvariantExpr); ok {
		return p.Expr
	}
	return e
}

func (ev *evaluator) rangeEval(prepSeries func(labels.Labels, *EvalSeriesHelper), funcCall func([]parser.Value, [][]EvalSeriesHelper, *EvalNodeHelper) Vector, exprs ...parser.Expr) Matrix {
	numSteps := 1

	matrixes := make([]Matrix, len(exprs))
	origMatrixes := make([]Matrix, len(exprs))

	for i, e := range exprs {
		// Functions will take string arguments from the expressions, not the values.
		if e != nil && e.Type() != parser.ValueTypeString {
			// ev.currentSamples will be updated to the correct value within the ev.eval call.
			val := ev.eval(e)
			matrixes[i] = val.(Matrix)

			// Keep a copy of the original point slices so that they
			// can be returned to the pool.
			origMatrixes[i] = make(Matrix, len(matrixes[i]))
			copy(origMatrixes[i], matrixes[i])
		}
	}

	vectors := make([]Vector, len(exprs))    // Input vectors for the function.
	args := make([]parser.Value, len(exprs)) // Argument to function.
	// Create an output vector that is as big as the input matrix with
	// the most time series.
	biggestLen := 1
	for i := range exprs {
		vectors[i] = make(Vector, 0, len(matrixes[i]))
		if len(matrixes[i]) > biggestLen {
			biggestLen = len(matrixes[i])
		}
	}

	enh := &EvalNodeHelper{Out: make(Vector, 0, biggestLen)}
	type seriesAndTimestamp struct {
		Series
		ts int64
	}
	seriess := make(map[uint64]seriesAndTimestamp, biggestLen) // Output series by series hash.
	var (
		seriesHelpers [][]EvalSeriesHelper
		bufHelpers    [][]EvalSeriesHelper // Buffer updated on each step
	)

	// If the series preparation function is provided, we should run it for
	// every single series in the matrix.
	//if prepSeries != nil {
	seriesHelpers = make([][]EvalSeriesHelper, len(exprs))
	bufHelpers = make([][]EvalSeriesHelper, len(exprs))

	if prepSeries != nil {
		seriesHelpers = make([][]EvalSeriesHelper, len(exprs))
		bufHelpers = make([][]EvalSeriesHelper, len(exprs))

		for i := range exprs {
			seriesHelpers[i] = make([]EvalSeriesHelper, len(matrixes[i]))
			bufHelpers[i] = make([]EvalSeriesHelper, len(matrixes[i]))

			for si, series := range matrixes[i] {
				prepSeries(series.Metric, &seriesHelpers[i][si])
			}
		}
	}

	for ts := ev.startTimestamp; ts <= ev.endTimestamp; ts += ev.interval {

		for i := range exprs {
			vectors[i] = vectors[i][:0]

			if prepSeries != nil {
				bufHelpers[i] = bufHelpers[i][:0]
			}

			for si, series := range matrixes[i] {
				switch {
				case len(series.Floats) > 0 && series.Floats[0].T == ts:
					vectors[i] = append(vectors[i], Sample{Metric: series.Metric, F: series.Floats[0].F, T: ts})
					// Move input vectors forward so we don't have to re-scan the same
					// past points at the next step.
					matrixes[i][si].Floats = series.Floats[1:]
				case len(series.Histograms) > 0 && series.Histograms[0].T == ts:
					vectors[i] = append(vectors[i], Sample{Metric: series.Metric, H: series.Histograms[0].H, T: ts})
					matrixes[i][si].Histograms = series.Histograms[1:]
				default:
					continue
				}
				if prepSeries != nil {
					bufHelpers[i] = append(bufHelpers[i], seriesHelpers[i][si])
				}
			}

			args[i] = vectors[i]
		}
		// Make the function call.
		enh.Ts = ts
		result := funcCall(args, bufHelpers, enh)
		enh.Out = result[:0] // Reuse result vector.

		// Add samples in output vector to output series.
		for _, sample := range result {
			h := sample.Metric.Hash()
			ss, ok := seriess[h]
			if ok {
				//if ss.ts == ts { // If we've seen this output series before at this timestamp, it's a duplicate.
				//	panic("vector cannot contain metrics with the same labelset")
				//}
				ss.ts = ts
			} else {
				ss = seriesAndTimestamp{Series{Metric: sample.Metric}, ts}
			}
			if sample.H == nil {
				if ss.Floats == nil {
					ss.Floats = getFPointSlice(numSteps)
				}
				ss.Floats = append(ss.Floats, FPoint{T: ts, F: sample.F})
			} else {
				if ss.Histograms == nil {
					ss.Histograms = getHPointSlice(numSteps)
				}
				ss.Histograms = append(ss.Histograms, HPoint{T: ts, H: sample.H})
			}
			seriess[h] = ss
		}
	}

	// Assemble the output matrix. By the time we get here we know we don't have too many samples.
	mat := make(Matrix, 0, len(seriess))
	for _, ss := range seriess {
		mat = append(mat, ss.Series)
	}

	return mat
}

func getFPointSlice(sz int) []FPoint {
	return make([]FPoint, 0, sz)
}

// getHPointSlice will return a HPoint slice of size max(maxPointsSliceSize, sz).
// This function is called with an estimated size which often can be over-estimated.
func getHPointSlice(sz int) []HPoint {

	return make([]HPoint, 0, sz)
}
func (ev *evaluator) recover(expr parser.Expr, errp *error) {

}

func (ev *evaluator) aggregation(e *parser.AggregateExpr, grouping []string, param interface{}, vec Vector, seriesHelper []EvalSeriesHelper, enh *EvalNodeHelper) Vector {
	op := e.Op
	result := map[uint64]*groupedAggregation{}

	for _, s := range vec {
		metric := s.Metric

		// We can use the pre-computed grouping key unless grouping labels have changed.
		var groupingKey uint64
		groupingKey = generateGroupingKey(metric, grouping)

		group, _ := result[groupingKey]

		switch op {
		case parser.SUM:
			group.hasFloat = true
			group.floatValue += s.F
		}
	}

	return nil
}

type groupedAggregation struct {
	hasFloat       bool // Has at least 1 float64 sample aggregated.
	hasHistogram   bool // Has at least 1 histogram sample aggregated.
	labels         labels.Labels
	floatValue     float64
	histogramValue *histogram.FloatHistogram
	floatMean      float64
	histogramMean  *histogram.FloatHistogram
	groupCount     int
}

// scalarBinop evaluates a binary operation between two Scalars.
func scalarBinop(op parser.ItemType, lhs, rhs float64) float64 {
	switch op {
	case parser.ADD:
		return lhs + rhs
	case parser.DIV:
		return lhs / rhs
	}
	panic(fmt.Errorf("operator %q not allowed for Scalar operations", op))
}
