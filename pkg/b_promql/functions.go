package promql

import (
	parser "github.com/dborchard/prometheus_lite/pkg/b_parser"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"math"
)

type FunctionCall func(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) Vector

// FunctionCalls is a list of all functions supported by PromQL, including their types.
var FunctionCalls = map[string]FunctionCall{
	"abs": funcAbs,
}

// === abs(Vector b_parser.ValueTypeVector) (Vector, Annotations) ===
func funcAbs(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) Vector {
	return simpleFunc(vals, enh, math.Abs)
}

func simpleFunc(vals []parser.Value, enh *EvalNodeHelper, f func(float64) float64) Vector {
	for _, el := range vals[0].(Vector) {
		if el.H == nil { // Process only float samples.
			enh.Out = append(enh.Out, Sample{
				Metric: enh.DropMetricName(el.Metric),
				F:      f(el.F),
			})
		}
	}
	return enh.Out
}

// DropMetricName is a cached version of DropMetricName.
func (enh *EvalNodeHelper) DropMetricName(l labels.Labels) labels.Labels {
	if enh.Dmn == nil {
		enh.Dmn = make(map[uint64]labels.Labels, len(enh.Out))
	}
	h := l.Hash()
	ret, ok := enh.Dmn[h]
	if ok {
		return ret
	}
	ret = dropMetricName(l)
	enh.Dmn[h] = ret
	return ret
}

func dropMetricName(l labels.Labels) labels.Labels {
	return labels.NewBuilder(l).Del(labels.MetricName).Labels()
}
