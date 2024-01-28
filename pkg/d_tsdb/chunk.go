package tsdb

// ChunkReader provides reading access of serialized time series data.
type ChunkReader interface {
	// ChunkOrIterable returns the series data for the given chunks.Meta.
	// Either a single chunk will be returned, or an iterable.
	// A single chunk should be returned if chunks.Meta maps to a chunk that
	// already exists and doesn't need modifications.
	// An iterable should be returned if chunks.Meta maps to a subset of the
	// samples in a stored chunk, or multiple chunks. (E.g. OOOHeadChunkReader
	// could return an iterable where multiple histogram samples have counter
	// resets. There can only be one counter reset per histogram chunk so
	// multiple chunks would be created from the iterable in this case.)
	// Only one of chunk or iterable should be returned. In some cases you may
	// always expect a chunk to be returned. You can check that iterable is nil
	// in those cases.
	ChunkOrIterable(meta chunks.Meta) (chunkenc.Chunk, chunkenc.Iterable, error)

	// Close releases all underlying resources of the reader.
	Close() error
}
