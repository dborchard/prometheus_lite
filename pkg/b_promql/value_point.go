package promql

import (
	"fmt"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"strconv"
)

// FPoint represents a single float data point for a given timestamp.
type FPoint struct {
	T int64
	F float64
}

// HPoint represents a single histogram data point for a given timestamp.
// H must never be nil.
type HPoint struct {
	T int64
	H *histogram.FloatHistogram
}

func (p FPoint) String() string {
	s := strconv.FormatFloat(p.F, 'f', -1, 64)
	return fmt.Sprintf("%s @[%v]", s, p.T)
}

func (p HPoint) String() string {
	//return fmt.Sprintf("%s @[%v]", p.H.String(), p.T)
	return fmt.Sprintf("%s @[%v]", "TODO", p.T)
}

// size returns the size of the HPoint compared to the size of an FPoint.
// The total size is calculated considering the histogram timestamp (p.T - 8 bytes),
// and then a number of bytes in the histogram.
// This sum is divided by 16, as samples are 16 bytes.
func (p HPoint) size() int {
	return (p.H.Size() + 8) / 16
}
