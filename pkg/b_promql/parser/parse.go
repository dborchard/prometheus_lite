package parser

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

func (p *parser) ParseExpr() (Expr, error) {
	//TODO implement me
	panic("implement me")
}

func (p *parser) Close() {
	//TODO implement me
	panic("implement me")
}

// ParseExpr returns the expression parsed from the input.
func ParseExpr(input string) (expr Expr, err error) {
	p := NewParser(input)
	defer p.Close()
	return p.ParseExpr()
}
