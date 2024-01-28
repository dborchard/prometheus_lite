package parser

import "prometheus_lite/pkg/b_promql/parser/posrange"

// Call represents a function call.
type Call struct {
	Func *Function   // The function that was called.
	Args Expressions // Arguments used in the call.

	PosRange posrange.PositionRange
}

func (c *Call) String() string {
	//TODO implement me
	panic("implement me")
}

func (c *Call) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (c *Call) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (c *Call) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (c *Call) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}

// Function represents a function of the expression language and is
// used by function nodes.
type Function struct {
	Name         string
	ArgTypes     []ValueType
	Variadic     int
	ReturnType   ValueType
	Experimental bool
}
