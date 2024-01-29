package storage

//
//import (
//	"context"
//	"github.com/dborchard/prometheus_lite/pkg/z_model/labels"
//)
//
//type mergeGenericQuerier struct {
//	queriers []genericQuerier
//
//	// mergeFn is used when we see series from different queriers Selects with the same labels.
//	mergeFn genericSeriesMergeFunc
//
//	// TODO(bwplotka): Remove once remote queries are asynchronous. False by default.
//	concurrentSelect bool
//}
//
//func (m *mergeGenericQuerier) LabelValues(ctx context.Context, name string, matchers ...*labels.Matcher) ([]string, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mergeGenericQuerier) LabelNames(ctx context.Context, matchers ...*labels.Matcher) ([]string, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mergeGenericQuerier) Close() error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mergeGenericQuerier) Select(ctx context.Context, b bool, hints *SelectHints, matcher ...*labels.Matcher) genericSeriesSet {
//	//TODO implement me
//	panic("implement me")
//}
//
//// NewMergeQuerier returns a new Querier that merges results of given primary and secondary queriers.
//// See NewFanout commentary to learn more about primary vs secondary differences.
////
//// In case of overlaps between the data given by primaries' and secondaries' Selects, merge function will be used.
//func NewMergeQuerier(primaries, secondaries []Querier, mergeFn VerticalSeriesMergeFunc) Querier {
//	queriers := make([]genericQuerier, 0, len(primaries)+len(secondaries))
//	for _, q := range primaries {
//		queriers = append(queriers, newGenericQuerierFrom(q))
//	}
//	//for _, q := range secondaries {
//	//	//queriers = append(queriers, newSecondaryQuerierFrom(q))
//	//
//	//}
//
//	concurrentSelect := false
//	if len(secondaries) > 0 {
//		concurrentSelect = true
//	}
//	return &querierAdapter{&mergeGenericQuerier{
//		mergeFn:          (&seriesMergerAdapter{VerticalSeriesMergeFunc: mergeFn}).Merge,
//		queriers:         queriers,
//		concurrentSelect: concurrentSelect,
//	}}
//}
//
//type seriesMergerAdapter struct {
//	VerticalSeriesMergeFunc
//}
//
//func (a *seriesMergerAdapter) Merge(s ...Labels) Labels {
//	buf := make([]Series, 0, len(s))
//	for _, ser := range s {
//		buf = append(buf, ser.(Series))
//	}
//	return a.VerticalSeriesMergeFunc(buf...)
//}
//
//// VerticalSeriesMergeFunc returns merged series implementation that merges series with same labels together.
//// It has to handle time-overlapped series as well.
//type VerticalSeriesMergeFunc func(...Series) Series
//
//type querierAdapter struct {
//	genericQuerier
//}
//
//func (q *querierAdapter) Select(ctx context.Context, sortSeries bool, hints *SelectHints, matchers ...*labels.Matcher) SeriesSet {
//	//TODO implement me
//	panic("implement me")
//}
