package promql

import (
	"context"
	"errors"
	"fmt"
	parser "github.com/dborchard/prometheus_lite/pkg/b_parser"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
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

	currentSamples int
	maxSamples     int
}

func (ev *evaluator) Eval(expr parser.Expr) (v parser.Value, err error) {
	defer ev.recover(expr, &err)

	v = ev.eval(expr)
	return v, nil
}

func (ev *evaluator) eval(expr parser.Expr) parser.Value {
	numSteps := int((ev.endTimestamp-ev.startTimestamp)/ev.interval) + 1

	switch e := expr.(type) {
	case *parser.AggregateExpr:
		sortedGrouping := e.Grouping
		slices.Sort(sortedGrouping)

		prepSeries := func(series labels.Labels, h *EvalSeriesHelper) {
			h.groupingKey = generateGroupingKey(series, sortedGrouping)
		}

		funcCall := func(v []parser.Value, sh [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
			return ev.aggregation(e, v[1].(Vector), sh[1], enh)
		}

		return ev.rangeEval(prepSeries, funcCall, e.Param, e.Expr)
	case *parser.Call:
		call := FunctionCalls[e.Func.Name]
		var (
			matrixArgIndex int
			matrixArg      bool
		)
		for i := range e.Args {
			unwrapParenExpr(&e.Args[i])
			a := unwrapStepInvariantExpr(e.Args[i])
			unwrapParenExpr(&a)
			if _, ok := a.(*parser.MatrixSelector); ok {
				matrixArgIndex = i
				matrixArg = true
				break
			}
		}
		if !matrixArg {
			// Does not have a matrix argument.
			return ev.rangeEval(nil, func(v []parser.Value, _ [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
				vec := call(v, e.Args, enh)
				return vec
			}, e.Args...)
		} else {
			// Evaluate any non-matrix arguments.
			inArgs := make([]parser.Value, len(e.Args))
			otherArgs := make([]Matrix, len(e.Args))
			otherInArgs := make([]Vector, len(e.Args))
			for i, e := range e.Args {
				if i != matrixArgIndex {
					val := ev.eval(e)
					otherArgs[i] = val.(Matrix)
					otherInArgs[i] = Vector{Sample{}}
					inArgs[i] = otherInArgs[i]
				}
			}

			unwrapParenExpr(&e.Args[matrixArgIndex])
			arg := unwrapStepInvariantExpr(e.Args[matrixArgIndex])
			unwrapParenExpr(&arg)
			sel := arg.(*parser.MatrixSelector)
			selVS := sel.VectorSelector.(*parser.VectorSelector)

			err := checkAndExpandSeriesSet(ev.ctx, sel)
			if err != nil {
				panic(err)
			}
			mat := make(Matrix, 0, len(selVS.Series)) // Output matrix.
			offset := durationMilliseconds(selVS.Offset)
			selRange := durationMilliseconds(sel.Range)
			stepRange := selRange
			if stepRange > ev.interval {
				stepRange = ev.interval
			}

			// NOTE: The unknown code starts from here.
			var floats []FPoint // Reuse objects across steps to save memory allocations.
			var histograms []HPoint
			inMatrix := make(Matrix, 1)
			inArgs[matrixArgIndex] = inMatrix
			enh := &EvalNodeHelper{Out: make(Vector, 0, 1)}

			it := storage.NewBuffer(selRange) // Process all the calls for one time series at a time.
			var chkIter chunkenc.Iterator
			for i, s := range selVS.Series {
				ev.currentSamples -= len(floats) + totalHPointSize(histograms)
				if floats != nil {
					floats = floats[:0]
				}
				if histograms != nil {
					histograms = histograms[:0]
				}
				chkIter = s.Iterator(chkIter)
				it.Reset(chkIter)
				metric := selVS.Series[i].Labels()
				ss := Series{
					Metric: dropMetricName(metric),
				}
				inMatrix[0].Metric = selVS.Series[i].Labels()
				for ts, step := ev.startTimestamp, -1; ts <= ev.endTimestamp; ts += ev.interval {
					step++
					// Set the non-matrix arguments.
					// They are scalar, so it is safe to use the step number
					// when looking up the argument, as there will be no gaps.
					for j := range e.Args {
						if j != matrixArgIndex {
							otherInArgs[j][0].F = otherArgs[j][0].Floats[step].F
						}
					}
					maxt := ts - offset
					mint := maxt - selRange
					// Evaluate the matrix selector for this series for this step.
					floats, histograms = ev.matrixIterSlice(it, mint, maxt, floats, histograms)
					if len(floats)+len(histograms) == 0 {
						continue
					}
					inMatrix[0].Floats = floats
					inMatrix[0].Histograms = histograms
					enh.Ts = ts
					outVec := call(inArgs, e.Args, enh)

					enh.Out = outVec[:0]
					if len(outVec) > 0 {
						if outVec[0].H == nil {
							if ss.Floats == nil {
								ss.Floats = getFPointSlice(numSteps)
							}
							ss.Floats = append(ss.Floats, FPoint{F: outVec[0].F, T: ts})
						} else {
							if ss.Histograms == nil {
								ss.Histograms = getHPointSlice(numSteps)
							}
							ss.Histograms = append(ss.Histograms, HPoint{H: outVec[0].H, T: ts})
						}
					}
					// Only buffer stepRange milliseconds from the second step on.
					it.ReduceDelta(stepRange)
				}
				histSamples := totalHPointSize(ss.Histograms)
				if len(ss.Floats)+histSamples > 0 {
					if ev.currentSamples+len(ss.Floats)+histSamples <= ev.maxSamples {
						mat = append(mat, ss)
						ev.currentSamples += len(ss.Floats) + histSamples
					}
				}
			}
			ev.currentSamples -= len(floats) + totalHPointSize(histograms)
		}
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

		case lt == parser.ValueTypeVector && rt == parser.ValueTypeVector:
			// Function to compute the join signature for each series.
			buf := make([]byte, 0, 1024)
			sigFn := signatureFunc(e.VectorMatching.On, buf, e.VectorMatching.MatchingLabels...)
			prepFn := func(series labels.Labels, h *EvalSeriesHelper) {
				h.signature = sigFn(series)
			}
			switch e.Op {
			default:
				callFn := func(v []parser.Value, sh [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
					return ev.VectorBinop(e.Op, v[0].(Vector), v[1].(Vector), e.VectorMatching, e.ReturnBool, sh[0], sh[1], enh)
				}
				return ev.rangeEval(prepFn, callFn, e.LHS, e.RHS)
			}

		case lt == parser.ValueTypeVector && rt == parser.ValueTypeScalar:
			return ev.rangeEval(nil, func(v []parser.Value, _ [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
				return ev.VectorscalarBinop(e.Op, v[0].(Vector), Scalar{V: v[1].(Vector)[0].F}, false, e.ReturnBool, enh)
			}, e.LHS, e.RHS)

		case lt == parser.ValueTypeScalar && rt == parser.ValueTypeVector:
			return ev.rangeEval(nil, func(v []parser.Value, _ [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
				return ev.VectorscalarBinop(e.Op, v[1].(Vector), Scalar{V: v[0].(Vector)[0].F}, true, e.ReturnBool, enh)
			}, e.LHS, e.RHS)
		}
	case *parser.NumberLiteral:
		return ev.rangeEval(nil, func(v []parser.Value, _ [][]EvalSeriesHelper, enh *EvalNodeHelper) Vector {
			return append(enh.Out, Sample{F: e.Val, Metric: labels.EmptyLabels()})
		})
	case *parser.StringLiteral:
		return String{V: e.Val, T: ev.startTimestamp}
	case *parser.VectorSelector:
		err := checkAndExpandSeriesSet(ev.ctx, e)
		if err != nil {
			panic(err)
		}
		mat := make(Matrix, 0, len(e.Series))
		it := storage.NewMemoizedEmptyIterator(durationMilliseconds(0))
		var chkIter chunkenc.Iterator
		for i, s := range e.Series {
			chkIter = s.Iterator(chkIter)
			it.Reset(chkIter)
			ss := Series{
				Metric: e.Series[i].Labels(),
			}

			for ts, step := ev.startTimestamp, -1; ts <= ev.endTimestamp; ts += ev.interval {
				step++
				_, f, h, ok := ev.vectorSelectorSingle(it, e, ts)
				if ok {
					if ev.currentSamples <= ev.maxSamples {
						if h == nil {
							if ss.Floats == nil {
								ss.Floats = getFPointSlice(numSteps)
							}
							ss.Floats = append(ss.Floats, FPoint{F: f, T: ts})
							ev.currentSamples++
						} else {
							if ss.Histograms == nil {
								ss.Histograms = getHPointSlice(numSteps)
							}
							point := HPoint{H: h, T: ts}
							ss.Histograms = append(ss.Histograms, point)
							histSize := point.size()
							ev.currentSamples += histSize
						}
					}
				}
			}

			if len(ss.Floats)+len(ss.Histograms) > 0 {
				mat = append(mat, ss)
			}
		}
		return mat
	case *parser.MatrixSelector:
		if ev.startTimestamp != ev.endTimestamp {
			panic(errors.New("cannot do range evaluation of matrix selector"))
		}
		return ev.matrixSelector(e)
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
	return metric.HashForLabels(grouping...)
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
