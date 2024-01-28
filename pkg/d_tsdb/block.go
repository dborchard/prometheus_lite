package tsdb

import "github.com/dborchard/prometheus_lite/pkg/d_tsdb/tombstones"

// BlockReader provides reading access to a data block.
type BlockReader interface {
	// Index returns an IndexReader over the block's data.
	Index() (IndexReader, error)

	// Chunks returns a ChunkReader over the block's data.
	Chunks() (ChunkReader, error)

	// Tombstones returns a tombstones.Reader over the block's deleted data.
	Tombstones() (tombstones.Reader, error)

	// Meta provides meta information about the block reader.
	Meta() BlockMeta

	// Size returns the number of bytes that the block takes up on disk.
	Size() int64
}
