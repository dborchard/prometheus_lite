package parser

import "github.com/dborchard/prometheus_lite/pkg/b_parser/posrange"

type ItemType int

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

// Expr is a generic interface for all expression types.
type Expr interface {
	Node

	// Type returns the type the expression evaluates to. It does not perform
	// in-depth checks as this is done at parsing-time.
	Type() ValueType
	// PromQLExpr ensures that no other types accidentally implement the interface.
	PromQLExpr()
}

var _ Expr = (*NumberLiteral)(nil)
var _ Expr = (*StringLiteral)(nil)
var _ Expr = (*StepInvariantExpr)(nil)
var _ Expr = (*ParenExpr)(nil)

var _ Expr = (*UnaryExpr)(nil)
var _ Expr = (*BinaryExpr)(nil)

var _ Expr = (*Call)(nil)
var _ Expr = (*AggregateExpr)(nil)

var _ Expr = (*MatrixSelector)(nil)
var _ Expr = (*VectorSelector)(nil)
