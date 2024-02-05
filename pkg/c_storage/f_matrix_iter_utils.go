package storage

import (
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
)

type bufType int

const (
	noBuf bufType = iota // Nothing yet stored in sampleRing.
	iBuf
	fBuf
	hBuf
	fhBuf
)

type sampleRing struct {
	delta int64

	// Lookback buffers. We use iBuf for mixed samples, but one of the three
	// concrete ones for homogenous samples. (Only one of the four bufs is
	// allowed to be populated!) This avoids the overhead of the interface
	// wrapper for the happy (and by far most common) case of homogenous
	// samples.
	fBuf     []fSample
	hBuf     []hSample
	fhBuf    []fhSample
	bufInUse bufType

	i int // Position of most recent element in ring buffer.
	f int // Position of first element in ring buffer.
	l int // Number of elements in buffer.
}

func (r *sampleRing) reset() {
	r.l = 0
	r.i = -1
	r.f = 0
}

// addF is a version of the add method specialized for fSample.
func (r *sampleRing) addF(s fSample) {

}

// addH is a version of the add method specialized for hSample.
func (r *sampleRing) addH(s hSample) {

}

// addFH is a version of the add method specialized for fhSample.
func (r *sampleRing) addFH(s fhSample) {

}

// reduceDelta lowers the buffered time delta, dropping any samples that are
// out of the new delta range.
func (r *sampleRing) reduceDelta(delta int64) bool {
	if delta > r.delta {
		return false
	}
	r.delta = delta

	if r.l == 0 {
		return true
	}

	return true
}

// Returns the current iterator. Invalidates previously returned iterators.
func (r *sampleRing) iterator() chunkenc.Iterator {
	return nil
}

type Sample interface {
	T() int64
	F() float64
	H() *histogram.Histogram
	FH() *histogram.FloatHistogram
	Type() chunkenc.ValueType
}

type fSample struct {
	t int64
	f float64
}

type hSample struct {
	t int64
	h *histogram.Histogram
}
type fhSample struct {
	t  int64
	fh *histogram.FloatHistogram
}
