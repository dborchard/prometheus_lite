package parser

import (
	"prometheus_lite/pkg/b_promql/parser/posrange"
	"time"
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

// Statement is a generic interface for all statements.
type Statement interface {
	Node

	// PromQLStmt ensures that no other type accidentally implements the interface
	PromQLStmt()
}

type EvalStmt struct {
	Expr Expr // Expression to be evaluated.

	// The time boundaries for the evaluation. If Start equals End an instant
	// is evaluated.
	Start, End time.Time
	// Time between two evaluated instants for the range [Start:End].
	Interval time.Duration
}

func (e *EvalStmt) String() string {
	//TODO implement me
	panic("implement me")
}

func (e *EvalStmt) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (e *EvalStmt) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (e *EvalStmt) PromQLStmt() {
	//TODO implement me
	panic("implement me")
}
