package parser

import "github.com/dborchard/prometheus_lite/pkg/b_promql/parser/posrange"

// UnaryExpr represents a unary operation on another expression.
// Currently unary operations are only supported for Scalars.
type UnaryExpr struct {
	Op   ItemType
	Expr Expr

	StartPos posrange.Pos
}

func (u *UnaryExpr) String() string {
	//TODO implement me
	panic("implement me")
}

func (u *UnaryExpr) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (u *UnaryExpr) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (u *UnaryExpr) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (u *UnaryExpr) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}
