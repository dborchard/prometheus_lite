package promql

import (
	"fmt"
	"prometheus_lite/pkg/b_promql/parser"
	"prometheus_lite/pkg/z_model/histogram"
	"prometheus_lite/pkg/z_model/labels"
	"strconv"
	"strings"
)

func (Matrix) Type() parser.ValueType { return parser.ValueTypeMatrix }
func (Vector) Type() parser.ValueType { return parser.ValueTypeVector }
func (Scalar) Type() parser.ValueType { return parser.ValueTypeScalar }
func (String) Type() parser.ValueType { return parser.ValueTypeString }

type Result struct {
	Err   error
	Value parser.Value
}

type Matrix []Series

func (m Matrix) Len() int {
	//TODO implement me
	panic("implement me")
}

func (m Matrix) Less(i, j int) bool {
	//TODO implement me
	panic("implement me")
}

func (m Matrix) Swap(i, j int) {
	//TODO implement me
	panic("implement me")
}

// Vector is basically only an alias for []Sample, but the contract is that
// in a Vector, all Samples have the same timestamp.
type Vector []Sample

var _ parser.Value = new(Matrix)
var _ parser.Value = new(Vector)

// Sample is a single sample belonging to a metric. It represents either a float
// sample or a histogram sample. If H is nil, it is a float sample. Otherwise,
// it is a histogram sample.
type Sample struct {
	T int64
	F float64
	H *histogram.FloatHistogram

	Metric labels.Labels
}

type Series struct {
	Metric     labels.Labels `json:"metric"`
	Floats     []FPoint      `json:"values,omitempty"`
	Histograms []HPoint      `json:"histograms,omitempty"`
}

func (s Series) String() string {
	// TODO(beorn7): This currently renders floats first and then
	// histograms, each sorted by timestamp. Maybe, in mixed series, that's
	// fine. Maybe, however, primary sorting by timestamp is preferred, in
	// which case this has to be changed.
	vals := make([]string, 0, len(s.Floats)+len(s.Histograms))
	for _, f := range s.Floats {
		vals = append(vals, f.String())
	}
	for _, h := range s.Histograms {
		vals = append(vals, h.String())
	}
	return fmt.Sprintf("%s =>\n%s", s.Metric, strings.Join(vals, "\n"))
}

// FPoint represents a single float data point for a given timestamp.
type FPoint struct {
	T int64
	F float64
}

func (p FPoint) String() string {
	s := strconv.FormatFloat(p.F, 'f', -1, 64)
	return fmt.Sprintf("%s @[%v]", s, p.T)
}

// String represents a string value.
type String struct {
	T int64
	V string
}

// HPoint represents a single histogram data point for a given timestamp.
// H must never be nil.
type HPoint struct {
	T int64
	H *histogram.FloatHistogram
}

func (p HPoint) String() string {
	//return fmt.Sprintf("%s @[%v]", p.H.String(), p.T)
	return fmt.Sprintf("%s @[%v]", "TODO", p.T)
}

func (s String) String() string {
	return s.V
}

func (m Matrix) String() string {
	//TODO implement me
	panic("implement me")
}

func (v Vector) String() string {
	//TODO implement me
	panic("implement me")
}

// Scalar is a data point that's explicitly not associated with a metric.
type Scalar struct {
	T int64
	V float64
}

func (s Scalar) String() string {
	v := strconv.FormatFloat(s.V, 'f', -1, 64)
	return fmt.Sprintf("scalar: %v @[%v]", v, s.T)
}
