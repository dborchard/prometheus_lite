package promql

import "prometheus_lite/pkg/b_promql/parser"

type Result struct {
	Err   error
	Value parser.Value
}

type Matrix []Series

type Series struct {
	//Metric     labels.Labels `json:"metric"`
	Floats     []FPoint `json:"values,omitempty"`
	Histograms []HPoint `json:"histograms,omitempty"`
}

// FPoint represents a single float data point for a given timestamp.
type FPoint struct {
	T int64
	F float64
}

// HPoint represents a single histogram data point for a given timestamp.
// H must never be nil.
type HPoint struct {
	T int64
	//H *histogram.FloatHistogram
}
