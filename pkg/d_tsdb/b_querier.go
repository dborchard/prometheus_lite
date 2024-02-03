package tsdb

import (
	"context"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"github.com/oklog/ulid"
	"time"
)

type BlockQuerier struct {
	*blockBaseQuerier
}

type blockBaseQuerier struct {
	blockID ulid.ULID

	closed bool

	mint, maxt int64
}

func (q *BlockQuerier) Select(ctx context.Context, sortSeries bool, hints *storage.SelectHints, ms ...*labels.Matcher) storage.SeriesSet {
	return &blockSeriesSet{
		count: 0,
	}
}

func (q *blockBaseQuerier) LabelValues(ctx context.Context, name string, matchers ...*labels.Matcher) ([]string, error) {
	return nil, nil
}

func (q *blockBaseQuerier) LabelNames(ctx context.Context, matchers ...*labels.Matcher) ([]string, error) {
	return nil, nil
}

func (q *blockBaseQuerier) Close() error {
	return nil
}

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
