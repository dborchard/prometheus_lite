package chunkenc

import (
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"math"
)

// ValueType defines the type of a value an Iterator points to.
type ValueType uint8

// Possible values for ValueType.
const (
	ValNone           ValueType = iota // No value at the current position.
	ValFloat                           // A simple float, retrieved with At.
	ValHistogram                       // A histogram, retrieve with AtHistogram, but AtFloatHistogram works, too.
	ValFloatHistogram                  // A floating-point histogram, retrieve with AtFloatHistogram.
)

// Iterator is a simple iterator that can only get the next value.
// Iterator iterates over the samples of a time series, in timestamp-increasing order.
type Iterator interface {
	// Next advances the iterator by one and returns the type of the value
	// at the new position (or ValNone if the iterator is exhausted).
	Next() ValueType
	// Seek advances the iterator forward to the first sample with a
	// timestamp equal or greater than t. If the current sample found by a
	// previous `Next` or `Seek` operation already has this property, Seek
	// has no effect. If a sample has been found, Seek returns the type of
	// its value. Otherwise, it returns ValNone, after which the iterator is
	// exhausted.
	Seek(t int64) ValueType
	// At returns the current timestamp/value pair if the value is a float.
	// Before the iterator has advanced, the behaviour is unspecified.
	At() (int64, float64)
	// AtHistogram returns the current timestamp/value pair if the value is
	// a histogram with integer counts. Before the iterator has advanced,
	// the behaviour is unspecified.
	AtHistogram() (int64, *histogram.Histogram)
	// AtFloatHistogram returns the current timestamp/value pair if the
	// value is a histogram with floating-point counts. It also works if the
	// value is a histogram with integer counts, in which case a
	// FloatHistogram copy of the histogram is returned. Before the iterator
	// has advanced, the behaviour is unspecified.
	AtFloatHistogram() (int64, *histogram.FloatHistogram)
	// AtT returns the current timestamp.
	// Before the iterator has advanced, the behaviour is unspecified.
	AtT() int64
	// Err returns the current error. It should be used only after the
	// iterator is exhausted, i.e. `Next` or `Seek` have returned ValNone.
	Err() error
}

// NewNopIterator returns a new chunk iterator that does not hold any data.
func NewNopIterator() Iterator {
	return nopIterator{}
}

type nopIterator struct{}

func (nopIterator) Next() ValueType                                      { return ValNone }
func (nopIterator) Seek(int64) ValueType                                 { return ValNone }
func (nopIterator) At() (int64, float64)                                 { return math.MinInt64, 0 }
func (nopIterator) AtHistogram() (int64, *histogram.Histogram)           { return math.MinInt64, nil }
func (nopIterator) AtFloatHistogram() (int64, *histogram.FloatHistogram) { return math.MinInt64, nil }
func (nopIterator) AtT() int64                                           { return math.MinInt64 }
func (nopIterator) Err() error                                           { return nil }
