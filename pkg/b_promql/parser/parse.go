package parser

import (
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser/posrange"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"github.com/prometheus/common/model"
)

type Parser interface {
	ParseExpr() (Expr, error)
	Close()
}

type parser struct {
}

var _ Parser = (*parser)(nil)

func NewParser(input string) *parser {
	return &parser{}
}

// ParseExpr parses an expression from the input.
// Right now Mocking 6/3 expr
//func (p *parser) ParseExpr() (Expr, error) {
//	return &BinaryExpr{
//		Op: DIV,
//		LHS: &NumberLiteral{
//			Val:      6,
//			PosRange: posrange.PositionRange{Start: 0, End: 1},
//		},
//		RHS: &NumberLiteral{
//			Val:      3,
//			PosRange: posrange.PositionRange{Start: 4, End: 5},
//		},
//	}, nil
//}

func (p *parser) ParseExpr() (Expr, error) {
	return &BinaryExpr{
		Op: DIV,
		LHS: &VectorSelector{
			Name: "foo",
			LabelMatchers: []*labels.Matcher{
				MustLabelMatcher(labels.MatchEqual, model.MetricNameLabel, "foo"),
			},
			PosRange: posrange.PositionRange{
				Start: 0,
				End:   3,
			},
		},
		RHS: &VectorSelector{
			Name: "bar",
			LabelMatchers: []*labels.Matcher{
				MustLabelMatcher(labels.MatchEqual, model.MetricNameLabel, "bar"),
			},
			PosRange: posrange.PositionRange{
				Start: 6,
				End:   9,
			},
		},
		VectorMatching: &VectorMatching{Card: CardOneToOne},
	}, nil
}

func (p *parser) Close() {
}

// ParseExpr returns the expression parsed from the input.
func ParseExpr(input string) (expr Expr, err error) {
	p := NewParser(input)
	defer p.Close()
	return p.ParseExpr()
}

func MustLabelMatcher(mt labels.MatchType, name, val string) *labels.Matcher {
	m, err := labels.NewMatcher(mt, name, val)
	if err != nil {
		panic(err)
	}
	return m
}
