package tsdb

import storage "github.com/dborchard/prometheus_lite/pkg/c_storage"

type blockSeriesSet struct {
	count int
}

func (b *blockSeriesSet) Next() bool {
	if b.count == 3 {
		return false
	}
	return true
}

func (b *blockSeriesSet) At() storage.Series {
	b.count++
	panic("implement me")
}

func (b *blockSeriesSet) Err() error {
	//TODO implement me
	panic("implement me")
}
