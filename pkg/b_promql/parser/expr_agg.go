package parser

import "prometheus_lite/pkg/b_promql/parser/posrange"

// AggregateExpr represents an aggregation operation on a Vector.
type AggregateExpr struct {
	Op       ItemType // The used aggregation operation.
	Expr     Expr     // The Vector expression over which is aggregated.
	Param    Expr     // Parameter used by some aggregators.
	Grouping []string // The labels by which to group the Vector.
	Without  bool     // Whether to drop the given labels rather than keep them.
	PosRange posrange.PositionRange
}

func (a *AggregateExpr) String() string {
	//TODO implement me
	panic("implement me")
}

func (a *AggregateExpr) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (a *AggregateExpr) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (a *AggregateExpr) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (a *AggregateExpr) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}
