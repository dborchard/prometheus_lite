package promql

import (
	"math"
	"prometheus_lite/pkg/b_promql/parser"
	"prometheus_lite/pkg/z_model/labels"
	"regexp"
)

// FunctionCalls is a list of all functions supported by PromQL, including their types.
var FunctionCalls = map[string]FunctionCall{
	"abs": funcAbs,
}

// === abs(Vector parser.ValueTypeVector) (Vector, Annotations) ===
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

type FunctionCall func(vals []parser.Value, args parser.Expressions, enh *EvalNodeHelper) Vector

// EvalNodeHelper stores extra information and caches for evaluating a single node across steps.
type EvalNodeHelper struct {
	// Evaluation timestamp.
	Ts int64
	// Vector that can be used for output.
	Out Vector

	// Caches.
	// DropMetricName and label_*.
	Dmn map[uint64]labels.Labels
	// funcHistogramQuantile for classic histograms.
	//signatureToMetricWithBuckets map[string]*metricWithBuckets
	// label_replace.
	regex *regexp.Regexp

	lb           *labels.Builder
	lblBuf       []byte
	lblResultBuf []byte

	// For binary vector matching.
	rightSigs    map[string]Sample
	matchedSigs  map[string]map[uint64]struct{}
	resultMetric map[string]labels.Labels
}
