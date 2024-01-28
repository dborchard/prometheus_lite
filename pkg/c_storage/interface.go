package storage

// SeriesSet contains a set of series.
type SeriesSet interface {
	Next() bool
	// At returns full series. Returned series should be iterable even after Next is called.
	At() Series
	// The error that iteration as failed with.
	// When an error occurs, set cannot continue to iterate.
	Err() error
}

// Series exposes a single time series and allows iterating over samples.
type Series interface {
	Labels
	SampleIterable
}

// Labels represents an item that has labels e.g. time series.
type Labels interface {
	// Labels returns the complete set of labels. For series it means all labels identifying the series.
	Labels() labels.Labels
}

type SampleIterable interface {
	// Iterator returns an iterator of the data of the series.
	// The iterator passed as argument is for re-use, if not nil.
	// Depending on implementation, the iterator can
	// be re-used or a new iterator can be allocated.
	Iterator(chunkenc.Iterator) chunkenc.Iterator
}
