package promql

import (
	"fmt"
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
)

func (ev *evaluator) rangeEval(prepSeries func(labels.Labels, *EvalSeriesHelper),
	funcCall func([]parser.Value, [][]EvalSeriesHelper, *EvalNodeHelper) Vector,
	exprs ...parser.Expr) Matrix {
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
				if ss.ts == ts { // If we've seen this output series before at this timestamp, it's a duplicate.
					panic("vector cannot contain metrics with the same labelset")
				}
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

// aggregation implemented SUM aggFn.
func (ev *evaluator) aggregation(e *parser.AggregateExpr, vec Vector, seriesHelper []EvalSeriesHelper, enh *EvalNodeHelper) Vector {
	op := e.Op
	result := map[uint64]*groupedAggregation{}
	var orderedResult []*groupedAggregation

	for si, s := range vec {
		var groupingKey = seriesHelper[si].groupingKey

		group, ok := result[groupingKey]
		if !ok {
			newAgg := &groupedAggregation{
				labels:     enh.lb.Labels(),
				floatValue: s.F,
				floatMean:  s.F,
				groupCount: 1,
			}
			result[groupingKey] = newAgg
			orderedResult = append(orderedResult, newAgg)
		}

		switch op {
		case parser.SUM:

			if s.H != nil {
				group.hasHistogram = true
				if group.histogramValue != nil {
					group.histogramValue.Add(s.H)
				}
			} else {
				group.hasFloat = true
				group.floatValue += s.F
			}
		}
	}

	// Construct the result Vector from the aggregated groups.
	for _, agg := range orderedResult {
		switch op {
		case parser.SUM:
			if agg.hasHistogram {
				agg.histogramValue.Compact(0)
			}
		default:
		}

		enh.Out = append(enh.Out, Sample{
			Metric: agg.labels,
			F:      agg.floatValue,
			H:      agg.histogramValue,
		})
	}

	return enh.Out
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

// VectorBinop evaluates a binary operation between two Vectors, excluding set operators.
func (ev *evaluator) VectorBinop(op parser.ItemType, lhs, rhs Vector, matching *parser.VectorMatching, returnBool bool, lhsh, rhsh []EvalSeriesHelper, enh *EvalNodeHelper) Vector {
	if matching.Card == parser.CardManyToMany {
		panic("many-to-many only allowed for set operators")
	}
	if len(lhs) == 0 || len(rhs) == 0 {
		return nil // Short-circuit: nothing is going to match.
	}

	// The control flow below handles one-to-one or many-to-one matching.
	// For one-to-many, swap sidedness and account for the swap when calculating
	// values.
	if matching.Card == parser.CardOneToMany {
		lhs, rhs = rhs, lhs
		lhsh, rhsh = rhsh, lhsh
	}

	// All samples from the rhs hashed by the matching label/values.
	if enh.rightSigs == nil {
		enh.rightSigs = make(map[string]Sample, len(enh.Out))
	} else {
		for k := range enh.rightSigs {
			delete(enh.rightSigs, k)
		}
	}
	rightSigs := enh.rightSigs

	// Add all rhs samples to a map so we can easily find matches later.
	for i, rs := range rhs {
		sig := rhsh[i].signature
		rightSigs[sig] = rs
	}

	// Tracks the match-signature. For one-to-one operations the value is nil. For many-to-one
	// the value is a set of signatures to detect duplicated result elements.
	if enh.matchedSigs == nil {
		enh.matchedSigs = make(map[string]map[uint64]struct{}, len(rightSigs))
	} else {
		for k := range enh.matchedSigs {
			delete(enh.matchedSigs, k)
		}
	}
	matchedSigs := enh.matchedSigs

	// For all lhs samples find a respective rhs sample and perform
	// the binary operation.
	for i, ls := range lhs {
		sig := lhsh[i].signature

		rs, found := rightSigs[sig] // Look for a match in the rhs Vector.
		if !found {
			continue
		}

		// Account for potentially swapped sidedness.
		fl, fr := ls.F, rs.F
		hl, hr := ls.H, rs.H
		if matching.Card == parser.CardOneToMany {
			fl, fr = fr, fl
			hl, hr = hr, hl
		}
		floatValue, histogramValue, keep := vectorElemBinop(op, fl, fr, hl, hr)
		switch {
		case returnBool:
			if keep {
				floatValue = 1.0
			} else {
				floatValue = 0.0
			}
		case !keep:
			continue
		}
		metric := resultMetric(ls.Metric, rs.Metric, op, matching, enh)
		if returnBool {
			metric = enh.DropMetricName(metric)
		}
		insertedSigs, exists := matchedSigs[sig]
		if matching.Card == parser.CardOneToOne {
			if exists {
				panic("multiple matches for labels: grouping labels must ensure unique matches")
			}
			matchedSigs[sig] = nil // Set existence to true.
		} else {
			// In many-to-one matching the grouping labels have to ensure a unique metric
			// for the result Vector. Check whether those labels have already been added for
			// the same matching labels.
			insertSig := metric.Hash()

			if !exists {
				insertedSigs = map[uint64]struct{}{}
				matchedSigs[sig] = insertedSigs
			} else if _, duplicate := insertedSigs[insertSig]; duplicate {
				panic("multiple matches for labels: grouping labels must ensure unique matches")
			}
			insertedSigs[insertSig] = struct{}{}
		}

		enh.Out = append(enh.Out, Sample{
			Metric: metric,
			F:      floatValue,
			H:      histogramValue,
		})
	}
	return enh.Out
}

// vectorElemBinop evaluates a binary operation between two Vector elements.
func vectorElemBinop(op parser.ItemType, lhs, rhs float64, hlhs, hrhs *histogram.FloatHistogram) (float64, *histogram.FloatHistogram, bool) {
	switch op {
	case parser.ADD:
		if hlhs != nil && hrhs != nil {
			return 0, hrhs.Copy().Add(hlhs).Compact(0), true
		}
		return lhs + rhs, nil, true
	case parser.DIV:
		if hlhs != nil && hrhs == nil {
			return 0, hlhs.Copy().Div(rhs), true
		}
		return lhs / rhs, nil, true
	}
	panic(fmt.Errorf("operator %q not allowed for operations between Vectors", op))
}

// VectorscalarBinop evaluates a binary operation between a Vector and a Scalar.
func (ev *evaluator) VectorscalarBinop(op parser.ItemType, lhs Vector, rhs Scalar, swap, returnBool bool, enh *EvalNodeHelper) Vector {
	for _, lhsSample := range lhs {
		lf, rf := lhsSample.F, rhs.V
		var rh *histogram.FloatHistogram
		lh := lhsSample.H
		// lhs always contains the Vector. If the original position was different
		// swap for calculating the value.
		if swap {
			lf, rf = rf, lf
			lh, rh = rh, lh
		}
		float, histogram, keep := vectorElemBinop(op, lf, rf, lh, rh)
		// Catch cases where the scalar is the LHS in a scalar-vector comparison operation.
		// We want to always keep the vector element value as the output value, even if it's on the RHS.
		if returnBool {
			if keep {
				float = 1.0
			} else {
				float = 0.0
			}
			keep = true
		}
		if keep {
			lhsSample.F = float
			lhsSample.H = histogram
			enh.Out = append(enh.Out, lhsSample)
		}
	}
	return enh.Out
}
