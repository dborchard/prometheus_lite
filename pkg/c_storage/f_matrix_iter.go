package storage

import (
	"fmt"
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"math"
)

// NewBuffer returns a new iterator that buffers the values within the time range
// of the current element and the duration of delta before, initialized with an
// empty iterator. Use Reset() to set an actual iterator to be buffered.
func NewBuffer(delta int64) *BufferedSeriesIterator {
	return NewBufferIterator(chunkenc.NewNopIterator(), delta)
}

// NewBufferIterator returns a new iterator that buffers the values within the
// time range of the current element and the duration of delta before.
func NewBufferIterator(it chunkenc.Iterator, delta int64) *BufferedSeriesIterator {
	bit := &BufferedSeriesIterator{
		delta: delta,
	}
	bit.Reset(it)

	return bit
}

// Reset r

// BufferedSeriesIterator wraps an iterator with a look-back buffer.
type BufferedSeriesIterator struct {
	it    chunkenc.Iterator
	buf   *sampleRing
	delta int64

	lastTime  int64
	valueType chunkenc.ValueType
}

// Reset re-uses the buffer with a new iterator, resetting the buffered time
// delta to its original value.
func (b *BufferedSeriesIterator) Reset(it chunkenc.Iterator) {
	b.it = it
	b.lastTime = math.MinInt64
	b.valueType = it.Next()
}

// Seek advances the iterator to the element at time t or greater.
func (b *BufferedSeriesIterator) Seek(t int64) chunkenc.ValueType {
	t0 := t - b.buf.delta

	// If the delta would cause us to seek backwards, preserve the buffer
	// and just continue regular advancement while filling the buffer on the way.
	if b.valueType != chunkenc.ValNone && t0 > b.lastTime {
		b.buf.reset()

		b.valueType = b.it.Seek(t0)
		switch b.valueType {
		case chunkenc.ValNone:
			return chunkenc.ValNone
		case chunkenc.ValFloat, chunkenc.ValHistogram, chunkenc.ValFloatHistogram:
			b.lastTime = b.AtT()
		default:
			panic(fmt.Errorf("BufferedSeriesIterator: unknown value type %v", b.valueType))
		}
	}

	if b.lastTime >= t {
		return b.valueType
	}
	for {
		if b.valueType = b.Next(); b.valueType == chunkenc.ValNone || b.lastTime >= t {
			return b.valueType
		}
	}
}

// Next advances the iterator to the next element.
func (b *BufferedSeriesIterator) Next() chunkenc.ValueType {
	// Add current element to buffer before advancing.
	switch b.valueType {
	case chunkenc.ValNone:
		return chunkenc.ValNone
	case chunkenc.ValFloat:
		t, f := b.it.At()
		b.buf.addF(fSample{t: t, f: f})
	case chunkenc.ValHistogram:
		t, h := b.it.AtHistogram()
		b.buf.addH(hSample{t: t, h: h})
	case chunkenc.ValFloatHistogram:
		t, fh := b.it.AtFloatHistogram()
		b.buf.addFH(fhSample{t: t, fh: fh})
	default:
		panic(fmt.Errorf("BufferedSeriesIterator: unknown value type %v", b.valueType))
	}

	b.valueType = b.it.Next()
	if b.valueType != chunkenc.ValNone {
		b.lastTime = b.AtT()
	}
	return b.valueType
}

// At returns the current float element of the iterator.
func (b *BufferedSeriesIterator) At() (int64, float64) {
	return b.it.At()
}

// AtHistogram returns the current histogram element of the iterator.
func (b *BufferedSeriesIterator) AtHistogram() (int64, *histogram.Histogram) {
	return b.it.AtHistogram()
}

// AtFloatHistogram returns the current float-histogram element of the iterator.
func (b *BufferedSeriesIterator) AtFloatHistogram() (int64, *histogram.FloatHistogram) {
	return b.it.AtFloatHistogram()
}

// AtT returns the current timestamp of the iterator.
func (b *BufferedSeriesIterator) AtT() int64 {
	return b.it.AtT()
}

// Err returns the last encountered error.
func (b *BufferedSeriesIterator) Err() error {
	return b.it.Err()
}

// ReduceDelta lowers the buffered time delta, for the current SeriesIterator only.
func (b *BufferedSeriesIterator) ReduceDelta(delta int64) bool {
	return b.buf.reduceDelta(delta)
}

// Buffer returns an iterator over the buffered data. Invalidates previously
// returned iterators.
func (b *BufferedSeriesIterator) Buffer() chunkenc.Iterator {
	return b.buf.iterator()
}
