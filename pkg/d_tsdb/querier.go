package tsdb

import (
	"context"
	"errors"
	"github.com/oklog/ulid"
	storage "prometheus_lite/pkg/c_storage"
	"prometheus_lite/pkg/d_tsdb/tombstones"
	"prometheus_lite/pkg/z_model/labels"
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

func newBlockBaseQuerier(b BlockReader, mint, maxt int64) (*blockBaseQuerier, error) {
	indexr, _ := b.Index()
	chunkr, _ := b.Chunks()
	tombsr, _ := b.Tombstones()

	return &blockBaseQuerier{
		blockID:    b.Meta().ULID,
		mint:       mint,
		maxt:       maxt,
		index:      indexr,
		chunks:     chunkr,
		tombstones: tombsr,
	}, nil
}

func (q *blockQuerier) Select(ctx context.Context, sortSeries bool, hints *storage.SelectHints, ms ...*labels.Matcher) storage.SeriesSet {
	mint := q.mint
	maxt := q.maxt
	disableTrimming := false

	p, _ := PostingsForMatchers(ctx, q.index, ms...)
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
	return nil
}
