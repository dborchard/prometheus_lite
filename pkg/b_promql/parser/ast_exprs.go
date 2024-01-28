package parser

import (
	"prometheus_lite/pkg/b_promql/parser/posrange"
)

// Expressions is a list of expression nodes that implements Node.
type Expressions []Expr

func (e Expressions) String() string {
	//TODO implement me
	panic("implement me")
}

func (e Expressions) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (e Expressions) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}
