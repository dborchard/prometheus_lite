package tsdb

import (
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"time"
)

// populateWithDelSeriesIterator allows to iterate over samples for the single series.
type populateWithDelSeriesIterator struct {
	count int
	curr  chunkenc.Iterator
}

func (p *populateWithDelSeriesIterator) Next() chunkenc.ValueType {
	if p.count == 0 {
		return chunkenc.ValNone
	}
	p.count--
	return chunkenc.ValFloat
}

func (p *populateWithDelSeriesIterator) Seek(t int64) chunkenc.ValueType {
	return chunkenc.ValFloat
}

func (p *populateWithDelSeriesIterator) At() (int64, float64) {
	// NOTE: this the entry point for storage.
	return 0, float64(p.count)
}

func (p *populateWithDelSeriesIterator) AtHistogram() (int64, *histogram.Histogram) {
	//TODO implement me
	panic("implement me")
}

func (p *populateWithDelSeriesIterator) AtFloatHistogram() (int64, *histogram.FloatHistogram) {
	//TODO implement me
	panic("implement me")
}

func (p *populateWithDelSeriesIterator) AtT() int64 {
	// NOTE: this is the place where we send the Vector Output.
	return time.Now().Unix()
}

func (p *populateWithDelSeriesIterator) Err() error {
	return nil
}
