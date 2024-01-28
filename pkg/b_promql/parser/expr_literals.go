package parser

import "prometheus_lite/pkg/b_promql/parser/posrange"

// NumberLiteral represents a number.
type NumberLiteral struct {
	Val float64

	PosRange posrange.PositionRange
}

// StringLiteral represents a string.
type StringLiteral struct {
	Val      string
	PosRange posrange.PositionRange
}

// StepInvariantExpr represents a query which evaluates to the same result
// irrespective of the evaluation time given the raw samples from TSDB remain unchanged.
// Currently this is only used for engine optimisations and the parser does not produce this.
type StepInvariantExpr struct {
	Expr Expr
}

func (n *NumberLiteral) String() string {
	//TODO implement me
	panic("implement me")
}

func (n *NumberLiteral) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (n *NumberLiteral) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (n *NumberLiteral) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (n *NumberLiteral) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}

func (s *StringLiteral) String() string {
	//TODO implement me
	panic("implement me")
}

func (s *StringLiteral) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (s *StringLiteral) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (s *StringLiteral) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (s *StringLiteral) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}

func (s *StepInvariantExpr) String() string {
	//TODO implement me
	panic("implement me")
}

func (s *StepInvariantExpr) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (s *StepInvariantExpr) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (s *StepInvariantExpr) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (s *StepInvariantExpr) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}
