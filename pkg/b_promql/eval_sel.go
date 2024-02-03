package promql

import (
	"context"
	"fmt"
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
)

// matrixSelector evaluates a *parser.MatrixSelector expression.
func (ev *evaluator) matrixSelector(node *parser.MatrixSelector) Matrix {
	var (
		vs = node.VectorSelector.(*parser.VectorSelector)

		offset = durationMilliseconds(vs.Offset)
		maxt   = ev.startTimestamp - offset
		mint   = maxt - durationMilliseconds(node.Range)
		matrix = make(Matrix, 0, len(vs.Series))

		it = storage.NewBuffer(durationMilliseconds(node.Range))
	)
	err := checkAndExpandSeriesSet(ev.ctx, node)
	if err != nil {
		panic(fmt.Errorf("error expanding series set: %w", err))
	}

	var chkIter chunkenc.Iterator
	series := vs.Series
	for i, s := range series {
		chkIter = s.Iterator(chkIter)
		it.Reset(chkIter)
		ss := Series{
			Metric: series[i].Labels(),
		}

		ss.Floats, ss.Histograms = ev.matrixIterSlice(it, mint, maxt, nil, nil)
		totalSize := int64(len(ss.Floats)) + int64(totalHPointSize(ss.Histograms))

		if totalSize > 0 {
			matrix = append(matrix, ss)
		} else {

		}
	}
	return matrix
}

func checkAndExpandSeriesSet(ctx context.Context, expr parser.Expr) error {
	switch e := expr.(type) {
	case *parser.MatrixSelector:
		return checkAndExpandSeriesSet(ctx, e.VectorSelector)
	case *parser.VectorSelector:
		if e.Series != nil {
			return nil
		}
		series, err := expandSeriesSet(ctx, e.UnexpandedSeriesSet)
		e.Series = series
		return err
	}
	return nil
}

// vectorSelectorSingle evaluates an instant vector for the iterator of one time series.
func (ev *evaluator) vectorSelectorSingle(it *storage.MemoizedSeriesIterator, node *parser.VectorSelector, ts int64) (
	int64, float64, *histogram.FloatHistogram, bool,
) {
	refTime := ts - durationMilliseconds(node.Offset)
	var t int64
	var v float64
	var h *histogram.FloatHistogram

	valueType := it.Seek(refTime)
	switch valueType {
	case chunkenc.ValNone:
		if it.Err() != nil {
			panic(it.Err())
		}
	case chunkenc.ValFloat:
		t, v = it.At()
	case chunkenc.ValFloatHistogram:
		t, h = it.AtFloatHistogram()
	default:
		panic(fmt.Errorf("unknown value type %v", valueType))
	}
	//if valueType == chunkenc.ValNone || t > refTime {
	//	var ok bool
	//	t, v, h, ok = it.PeekPrev()
	//	if !ok || t < refTime-durationMilliseconds(0) {
	//		return 0, 0, nil, false
	//	}
	//}
	return t, v, h, true
}
