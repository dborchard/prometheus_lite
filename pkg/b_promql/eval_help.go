package promql

// EvalSeriesHelper stores extra information about a series.
type EvalSeriesHelper struct {
	// The grouping key used by aggregation.
	groupingKey uint64
	// Used to map left-hand to right-hand in binary operations.
	signature string
}
