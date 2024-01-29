package v1

import (
	promql "github.com/dborchard/prometheus_lite/pkg/b_promql"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"github.com/dborchard/prometheus_lite/pkg/z_util/httputil"
	"net/http"
	"time"
)

type API struct {
	QueryEngine QueryEngine
	Queryable   storage.SampleAndChunkQueryable
}

func NewAPI(qe QueryEngine, q storage.SampleAndChunkQueryable) *API {
	return &API{
		QueryEngine: qe,
		Queryable:   q,
	}
}

func (api *API) queryRange(r *http.Request) (result apiFuncResult) {

	// 1. Args
	ctx := r.Context()
	start, _ := parseTime(r.FormValue("start"))
	end, _ := parseTime(r.FormValue("end"))
	step, _ := parseDuration(r.FormValue("step"))
	opts := promql.NewPrometheusQueryOpts(r.FormValue("stats") == "all", time.Duration(10))

	// 2. Create Logical Plan
	qry, _ := api.QueryEngine.NewRangeQuery(ctx, api.Queryable, opts, r.FormValue("query"), start, end, step)
	defer func() {
		if result.finalizer == nil {
			qry.Close()
		}
	}()

	// 3. Execute Logical Plan
	ctx = httputil.ContextFromRequest(ctx, r)
	res := qry.Exec(ctx)
	if res.Err != nil {
		return apiFuncResult{nil, nil, qry.Close}
	}

	// 4. Return Result
	return apiFuncResult{&QueryData{
		ResultType: res.Value.Type(),
		Result:     res.Value,
	}, nil, qry.Close}
}
