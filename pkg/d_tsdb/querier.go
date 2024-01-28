package d_tsdb

import (
	"context"
	"errors"
	"prometheus_lite/pkg/a_model/labels"
	storage "prometheus_lite/pkg/c_storage"
)

type blockQuerier struct {
	*blockBaseQuerier
}

type blockBaseQuerier struct {
	blockID    ulid.ULID
	index      IndexReader
	chunks     ChunkReader
	tombstones tombstones.Reader

	closed bool

	mint, maxt int64
}

func (q *blockQuerier) Select(ctx context.Context, sortSeries bool, hints *storage.SelectHints, ms ...*labels.Matcher) storage.SeriesSet {
	mint := q.mint
	maxt := q.maxt
	disableTrimming := false

	p, err := PostingsForMatchers(ctx, q.index, ms...)
	if err != nil {
		return storage.ErrSeriesSet(err)
	}
	if sortSeries {
		p = q.index.SortedPostings(p)
	}

	if hints != nil {
		mint = hints.Start
		maxt = hints.End
		disableTrimming = hints.DisableTrimming
		if hints.Func == "series" {
			// When you're only looking up metadata (for example series API), you don't need to load any chunks.
			return newBlockSeriesSet(q.index, newNopChunkReader(), q.tombstones, p, mint, maxt, disableTrimming)
		}
	}

	return newBlockSeriesSet(q.index, q.chunks, q.tombstones, p, mint, maxt, disableTrimming)
}
func (q *blockBaseQuerier) LabelValues(ctx context.Context, name string, matchers ...*labels.Matcher) ([]string, error) {
	res, err := q.index.SortedLabelValues(ctx, name, matchers...)
	return res, err
}

func (q *blockBaseQuerier) LabelNames(ctx context.Context, matchers ...*labels.Matcher) ([]string, error) {
	res, err := q.index.LabelNames(ctx, matchers...)
	return res, err
}

func (q *blockBaseQuerier) Close() error {
	if q.closed {
		return errors.New("block querier already closed")
	}

	errs := tsdb_errors.NewMulti(
		q.index.Close(),
		q.chunks.Close(),
		q.tombstones.Close(),
	)
	q.closed = true
	return errs.Err()
}
