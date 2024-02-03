package parser

import "github.com/dborchard/prometheus_lite/pkg/b_parser/posrange"

// ParenExpr wraps an expression so it cannot be disassembled as a consequence
// of operator precedence.
type ParenExpr struct {
	Expr     Expr
	PosRange posrange.PositionRange
}

func (p *ParenExpr) String() string {
	//TODO implement me
	panic("implement me")
}

func (p *ParenExpr) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (p *ParenExpr) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (p *ParenExpr) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (p *ParenExpr) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}
