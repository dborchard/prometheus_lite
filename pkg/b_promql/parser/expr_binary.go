package parser

import (
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser/posrange"
)

// BinaryExpr represents a binary expression between two child expressions.
type BinaryExpr struct {
	Op       ItemType // The operation of the expression.
	LHS, RHS Expr     // The operands on the respective sides of the operator.

	// The matching behavior for the operation if both operands are Vectors.
	// If they are not this field is nil.
	VectorMatching *VectorMatching

	// If a comparison operator, return 0/1 rather than filtering.
	ReturnBool bool
}

func (b *BinaryExpr) String() string {
	//TODO implement me
	panic("implement me")
}

func (b *BinaryExpr) Pretty(level int) string {
	//TODO implement me
	panic("implement me")
}

func (b *BinaryExpr) PositionRange() posrange.PositionRange {
	//TODO implement me
	panic("implement me")
}

func (b *BinaryExpr) Type() ValueType {
	//TODO implement me
	panic("implement me")
}

func (b *BinaryExpr) PromQLExpr() {
	//TODO implement me
	panic("implement me")
}

// VectorMatching describes how elements from two Vectors in a binary
// operation are supposed to be matched.
type VectorMatching struct {
	// The cardinality of the two Vectors.
	Card VectorMatchCardinality
	// MatchingLabels contains the labels which define equality of a pair of
	// elements from the Vectors.
	MatchingLabels []string
	// On includes the given label names from matching,
	// rather than excluding them.
	On bool
	// Include contains additional labels that should be included in
	// the result from the side with the lower cardinality.
	Include []string
}

// VectorMatchCardinality describes the cardinality relationship
// of two Vectors in a binary operation.
type VectorMatchCardinality int

const (
	CardOneToOne VectorMatchCardinality = iota
	CardManyToOne
	CardOneToMany
	CardManyToMany
)
