package storage

import (
	"context"
	"prometheus_lite/pkg/z_model/labels"
)

type genericQuerierAdapter struct {
	LabelQuerier

	// One-of. If both are set, Querier will be used.
	q  Querier
	cq ChunkQuerier
}

func (g *genericQuerierAdapter) Select(ctx context.Context, b bool, hints *SelectHints, matcher ...*labels.Matcher) genericSeriesSet {
	//TODO implement me
	panic("implement me")
}

func newGenericQuerierFrom(q Querier) genericQuerier {
	return &genericQuerierAdapter{LabelQuerier: q, q: q}
}

type genericSeriesMergeFunc func(...Labels) Labels

type genericQuerier interface {
	LabelQuerier
	Select(context.Context, bool, *SelectHints, ...*labels.Matcher) genericSeriesSet
}

type genericSeriesSet interface {
	Next() bool
	At() Labels
	Err() error
}
