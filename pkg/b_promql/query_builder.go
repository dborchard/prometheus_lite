package promql

import (
	parser "github.com/dborchard/prometheus_lite/pkg/b_parser"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"time"
)

func (ng *Engine) newQuery(q storage.Queryable, qs string, opts QueryOpts, start, end time.Time, interval time.Duration) (*parser.Expr, *query) {
	if opts == nil {
		opts = NewPrometheusQueryOpts(false, 0)
	}

	es := &parser.EvalStmt{
		Start:    start,
		End:      end,
		Interval: interval,
	}
	qry := &query{
		q:    qs,
		stmt: es,
		ng:   ng,
		cancel: func() {

		},
		queryable: q,
	}
	return &es.Expr, qry
}
