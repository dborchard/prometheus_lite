package parser

import (
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser/posrange"
)

type Node interface {
	// String representation of the node that returns the given node when parsed
	// as part of a valid query.
	String() string

	// Pretty returns the prettified representation of the node.
	// It uses the level information to determine at which level/depth the current
	// node is in the AST and uses this to apply indentation.
	Pretty(level int) string

	// PositionRange returns the position of the AST Node in the query string.
	PositionRange() posrange.PositionRange
}

var _ Node = (*Expressions)(nil)
var _ Node = (*EvalStmt)(nil)
