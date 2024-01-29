package storage

//
//import (
//	"log"
//)
//
//type fanout struct {
//	logger log.Logger
//	//
//	primary     Storage
//	secondaries []Storage
//}
//
//// NewFanout returns a new fanout Storage, which proxies reads and writes
//// through to multiple underlying storages.
////
//// The difference between primary and secondary Storage is only for read (Querier) path and it goes as follows:
//// * If the primary querier returns an error, then any of the Querier operations will fail.
//// * If any secondary querier returns an error the result from that queries is discarded. The overall operation will succeed,
//// and the error from the secondary querier will be returned as a warning.
////
//// NOTE: In the case of Prometheus, it treats all remote storages as secondary / best effort.
//func NewFanout(logger log.Logger, primary Storage, secondaries ...Storage) Storage {
//	return &fanout{
//		logger:      logger,
//		primary:     primary,
//		secondaries: secondaries,
//	}
//}
//
//func (f *fanout) Querier(mint, maxt int64) (Querier, error) {
//	primary, err := f.primary.Querier(mint, maxt)
//	if err != nil {
//		return nil, err
//	}
//
//	secondaries := make([]Querier, 0, len(f.secondaries))
//	for _, storage := range f.secondaries {
//		querier, err := storage.Querier(mint, maxt)
//		if err != nil {
//			panic(err)
//		}
//		secondaries = append(secondaries, querier)
//	}
//	return NewMergeQuerier([]Querier{primary}, secondaries, ChainedSeriesMerge), nil
//}
//
//func (f *fanout) ChunkQuerier(mint, maxt int64) (ChunkQuerier, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (f *fanout) StartTime() (int64, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (f *fanout) Close() error {
//	//TODO implement me
//	panic("implement me")
//}
//
//// ChainedSeriesMerge returns single series from many same, potentially overlapping series by chaining samples together.
//// If one or more samples overlap, one sample from random overlapped ones is kept and all others with the same
//// timestamp are dropped.
////
//// This works the best with replicated series, where data from two series are exactly the same. This does not work well
//// with "almost" the same data, e.g. from 2 Prometheus HA replicas. This is fine, since from the Prometheus perspective
//// this never happens.
////
//// It's optimized for non-overlap cases as well.
//func ChainedSeriesMerge(series ...Series) Series {
//	panic("implement me")
//}
