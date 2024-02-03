package tsdb

import (
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"github.com/oklog/ulid"
)

type blockSeriesEntry struct {
	blockID ulid.ULID
}

func (s *blockSeriesEntry) Labels() labels.Labels {
	return nil
}

func (s *blockSeriesEntry) Iterator(it chunkenc.Iterator) chunkenc.Iterator {
	pi := &populateWithDelSeriesIterator{
		count: 3,
	}
	return pi
}
