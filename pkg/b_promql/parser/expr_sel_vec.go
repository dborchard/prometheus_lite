package parser

import (
	"prometheus_lite/pkg/b_promql/parser/posrange"
	storage "prometheus_lite/pkg/c_storage"
	"prometheus_lite/pkg/z_model/labels"
	"time"
)

// VectorSelector represents a Vector selection.
type VectorSelector struct {
	Name string
	// OriginalOffset is the actual offset that was set in the query.
	// This never changes.
	OriginalOffset time.Duration
	// Offset is the offset used during the query execution
	// which is calculated using the original offset, at modifier time,
	// eval time, and subquery offsets in the AST tree.
	Offset     time.Duration
	Timestamp  *int64
	StartOrEnd ItemType // Set when @ is used with start() or end()

	// The unexpanded seriesSet populated at query preparation time.
	UnexpandedSeriesSet storage.SeriesSet
	Series              []storage.Series
	LabelMatchers       []*labels.Matcher

	PosRange posrange.PositionRange
}

func (v *VectorSelector) String() string {
	//TODO implement me
	panic("implement me")
}

func (v *VectorSelector) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (v *VectorSelector) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (v *VectorSelector) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (v *VectorSelector) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}
