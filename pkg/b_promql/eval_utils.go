package promql

import (
	"bytes"
	"context"
	parser "github.com/dborchard/prometheus_lite/pkg/b_parser"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"slices"
)

// unwrapParenExpr does the AST equivalent of removing parentheses around a expression.
func unwrapParenExpr(e *parser.Expr) {
	for {
		if p, ok := (*e).(*parser.ParenExpr); ok {
			*e = p.Expr
		} else {
			break
		}
	}
}

func unwrapStepInvariantExpr(e parser.Expr) parser.Expr {
	if p, ok := e.(*parser.StepInvariantExpr); ok {
		return p.Expr
	}
	return e
}

func expandSeriesSet(ctx context.Context, it storage.SeriesSet) (res []storage.Series, err error) {
	for it.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		res = append(res, it.At())
	}
	return res, it.Err()
}

// totalHPointSize returns the total number of samples in the given slice of HPoints.
func totalHPointSize(histograms []HPoint) int {
	var total int
	for _, h := range histograms {
		total += h.size()
	}
	return total
}

// resultMetric returns the metric for the given sample(s) based on the Vector
// binary operation and the matching options.
func resultMetric(lhs, rhs labels.Labels, op parser.ItemType, matching *parser.VectorMatching, enh *EvalNodeHelper) labels.Labels {
	if enh.resultMetric == nil {
		enh.resultMetric = make(map[string]labels.Labels, len(enh.Out))
	}

	enh.resetBuilder(lhs)
	buf := bytes.NewBuffer(enh.lblResultBuf[:0])
	enh.lblBuf = lhs.Bytes(enh.lblBuf)
	buf.Write(enh.lblBuf)
	enh.lblBuf = rhs.Bytes(enh.lblBuf)
	buf.Write(enh.lblBuf)
	enh.lblResultBuf = buf.Bytes()

	if ret, ok := enh.resultMetric[string(enh.lblResultBuf)]; ok {
		return ret
	}
	str := string(enh.lblResultBuf)

	if matching.Card == parser.CardOneToOne {
		enh.lb.Del(matching.MatchingLabels...)
	}
	for _, ln := range matching.Include {
		// Included labels from the `group_x` modifier are taken from the "one"-side.
		if v := rhs.Get(ln); v != "" {
			enh.lb.Set(ln, v)
		} else {
			enh.lb.Del(ln)
		}
	}

	ret := enh.lb.Labels()
	enh.resultMetric[str] = ret
	return ret
}

func (enh *EvalNodeHelper) resetBuilder(lbls labels.Labels) {
	if enh.lb == nil {
		enh.lb = labels.NewBuilder(lbls)
	} else {
		enh.lb.Reset(lbls)
	}
}

func signatureFunc(on bool, b []byte, names ...string) func(labels.Labels) string {
	names = append([]string{labels.MetricName}, names...)
	slices.Sort(names)
	return func(lset labels.Labels) string {
		return string(lset.BytesWithoutLabels(b, names...))
	}
}
