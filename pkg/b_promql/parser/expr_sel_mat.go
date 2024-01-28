package parser

import (
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser/posrange"
	"time"
)

// MatrixSelector represents a Matrix selection.
type MatrixSelector struct {
	// It is safe to assume that this is an VectorSelector
	// if the parser hasn't returned an error.
	VectorSelector Expr
	Range          time.Duration

	EndPos posrange.Pos
}

func (m *MatrixSelector) String() string {
	//TODO implement me
	panic("implement me")
}

func (m *MatrixSelector) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (m *MatrixSelector) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (m *MatrixSelector) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (m *MatrixSelector) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}
