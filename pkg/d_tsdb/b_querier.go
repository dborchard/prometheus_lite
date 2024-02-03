package tsdb

import (
	"context"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"github.com/oklog/ulid"
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
