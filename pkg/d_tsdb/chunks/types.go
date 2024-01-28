package chunks

import "prometheus_lite/pkg/d_tsdb/chunkenc"

// Iterator iterates over the chunks of a single time series.
type Iterator interface {
	// At returns the current meta.
	// It depends on implementation if the chunk is populated or not.
	At() Meta
	// Next advances the iterator by one.
	Next() bool
	// Err returns optional error if Next is false.
	Err() error
}

// Meta holds information about one or more chunks.
// For examples of when chunks.Meta could refer to multiple chunks, see
// ChunkReader.ChunkOrIterable().
type Meta struct {
	// Ref and Chunk hold either a reference that can be used to retrieve
	// chunk data or the data itself.
	// If Chunk is nil, call ChunkReader.ChunkOrIterable(Meta.Ref) to get the
	// chunk and assign it to the Chunk field. If an iterable is returned from
	// that method, then it may not be possible to set Chunk as the iterable
	// might form several chunks.
	Ref   ChunkRef
	Chunk chunkenc.Chunk

	// Time range the data covers.
	// When MaxTime == math.MaxInt64 the chunk is still open and being appended to.
	MinTime, MaxTime int64

	// OOOLastRef, OOOLastMinTime and OOOLastMaxTime are kept as markers for
	// overlapping chunks.
	// These fields point to the last created out of order Chunk (the head) that existed
	// when Series() was called and was overlapping.
	// Series() and Chunk() method responses should be consistent for the same
	// query even if new data is added in between the calls.
	OOOLastRef                     ChunkRef
	OOOLastMinTime, OOOLastMaxTime int64
}
