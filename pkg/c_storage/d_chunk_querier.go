package storage

import (
	"context"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
)

// ChunkQuerier provides querying access over time series data of a fixed time range.
type ChunkQuerier interface {
	LabelQuerier

	// Select returns a set of series that matches the given label matchers.
	// Caller can specify if it requires returned series to be sorted. Prefer not requiring sorting for better performance.
	// It allows passing hints that can help in optimising select, but it's up to implementation how this is used if used at all.
	Select(ctx context.Context, sortSeries bool, hints *SelectHints, matchers ...*labels.Matcher) ChunkSeriesSet
}

// ChunkSeriesSet contains a set of chunked series.
type ChunkSeriesSet interface {
	Next() bool
	At() ChunkSeries // At returns full chunk series. Returned series should be iterable even after Next is called.
	Err() error
}

// ChunkSeries exposes a single time series and allows iterating over chunks.
type ChunkSeries interface {
	Labels
	ChunkIterable
}

type ChunkIterable interface {
	// Iterator returns an iterator that iterates over potentially overlapping
	// chunks of the series, sorted by min time.
	//Iterator(chunks.Iterator) chunks.Iterator
}
