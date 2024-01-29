package tsdb

import (
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
)

type DB struct {
}

func NewDB() *DB {
	return &DB{}
}

func (D DB) Querier(mint, maxt int64) (storage.Querier, error) {
	return &BlockQuerier{}, nil
}

func (D DB) ChunkQuerier(mint, maxt int64) (storage.ChunkQuerier, error) {
	//TODO implement me
	panic("implement me")
}

func (D DB) StartTime() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (D DB) Close() error {
	//TODO implement me
	panic("implement me")
}
