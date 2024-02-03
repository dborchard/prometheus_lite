package parser

import "fmt"

// Visitor allows visiting a Node and its child nodes. The Visit method is
// invoked for each node with the path leading to the node provided additionally.
// If the result visitor w is not nil and no error, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil, nil).
type Visitor interface {
	Visit(node Node, path []Node) (w Visitor, err error)
}

type inspector func(Node, []Node) error

func (f inspector) Visit(node Node, path []Node) (Visitor, error) {
	if err := f(node, path); err != nil {
		return nil, err
	}

	return f, nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node, path); node must not be nil. If f returns a nil error, Inspect invokes f
// for all the non-nil children of node, recursively.
func Inspect(node Node, f inspector) {
	//nolint: errcheck
	Walk(f, node, nil)
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node, path); node must not be nil. If the visitor w returned by
// v.Visit(node, path) is not nil and the visitor returns no error, Walk is
// invoked recursively with visitor w for each of the non-nil children of node,
// followed by a call of w.Visit(nil), returning an error
// As the tree is descended the path of previous nodes is provided.
func Walk(v Visitor, node Node, path []Node) error {
	var err error
	if v, err = v.Visit(node, path); v == nil || err != nil {
		return err
	}
	path = append(path, node)

	for _, e := range Children(node) {
		if err := Walk(v, e, path); err != nil {
			return err
		}
	}

	_, err = v.Visit(nil, nil)
	return err
}

// Children returns a list of all child nodes of a syntax tree node.
func Children(node Node) []Node {
	// For some reasons these switches have significantly better performance than interfaces
	switch n := node.(type) {
	case *EvalStmt:
		return []Node{n.Expr}
	case Expressions:
		// golang cannot convert slices of interfaces
		ret := make([]Node, len(n))
		for i, e := range n {
			ret[i] = e
		}
		return ret
	case *AggregateExpr:
		// While this does not look nice, it should avoid unnecessary allocations
		// caused by slice resizing
		switch {
		case n.Expr == nil && n.Param == nil:
			return nil
		case n.Expr == nil:
			return []Node{n.Param}
		case n.Param == nil:
			return []Node{n.Expr}
		default:
			return []Node{n.Expr, n.Param}
		}
	case *BinaryExpr:
		return []Node{n.LHS, n.RHS}
	case *Call:
		// golang cannot convert slices of interfaces
		ret := make([]Node, len(n.Args))
		for i, e := range n.Args {
			ret[i] = e
		}
		return ret
	case *UnaryExpr:
		return []Node{n.Expr}
	case *MatrixSelector:
		return []Node{n.VectorSelector}
	case *StepInvariantExpr:
		return []Node{n.Expr}
	case *NumberLiteral, *StringLiteral, *VectorSelector:
		// nothing to do
		return []Node{}
	default:
		panic(fmt.Errorf("promql.Children: unhandled node type %T", node))
	}
}
