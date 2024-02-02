package promql

import (
	"fmt"
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"strings"
)

type Result struct {
	Err   error
	Value parser.Value
}

var _ parser.Value = new(Matrix)
var _ parser.Value = new(Vector)
var _ parser.Value = new(Scalar)
var _ parser.Value = new(String)

func (Matrix) Type() parser.ValueType { return parser.ValueTypeMatrix }
func (Vector) Type() parser.ValueType { return parser.ValueTypeVector }
func (Scalar) Type() parser.ValueType { return parser.ValueTypeScalar }
func (String) Type() parser.ValueType { return parser.ValueTypeString }

/*
	Matrix & Vector
	+---------+----+----+----+
	|    x    | T0 | T1 | T2 |
	+---------+----+----+----+
	| metric1 |  1 |  2 |  3 | <- Series
	| metric2 |  2 |  3 |  4 | <- Series / Vector
	+---------+----+----+----+
*/

type Matrix []Series // Matrix is a table with x-axis as metrics and y-axis as time. Ie list of columns.
type Series struct { // It is kind of a column in a matrix.
	Metric     labels.Labels `json:"metric"`
	Floats     []FPoint      `json:"values,omitempty"`
	Histograms []HPoint      `json:"histograms,omitempty"`
}

type Vector []Sample // Vector is kind of a singe column of homogenous data points.
type Sample struct { // IT is single sample belonging to a metric. It is either a float sample or a histogram sample.
	T int64
	F float64
	H *histogram.FloatHistogram

	Metric labels.Labels
}

type Scalar struct { // Scalar is a data point that's explicitly not associated with a metric.
	T int64
	V float64
}
type String struct { // String represents a string value.
	T int64
	V string
}

func (s Series) String() string {
	// This currently renders floats first and then
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
