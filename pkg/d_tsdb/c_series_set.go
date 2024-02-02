package tsdb

import storage "github.com/dborchard/prometheus_lite/pkg/c_storage"

type blockSeriesSet struct {
	count int
}

func (b *blockSeriesSet) Next() bool {
	if b.count == 3 {
		return false
	}
	b.count++
	return true
}

func (b *blockSeriesSet) At() storage.Series {
	// At can be looped over before iterating, so save the current values locally.
	return &blockSeriesEntry{}
}

func (b *blockSeriesSet) Err() error {
	return nil
}
