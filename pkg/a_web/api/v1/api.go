package v1

import "net/http"

type API struct {
	QueryEngine QueryEngine
}

func (api *API) queryRange(r *http.Request) (result apiFuncResult) {

	qry, _ := api.QueryEngine.NewRangeQuery(ctx, api.Queryable, opts, r.FormValue("query"), start, end, step)
	defer func() {
		if result.finalizer == nil {
			qry.Close()
		}
	}()
	ctx = httputil.ContextFromRequest(ctx, r)
	res := qry.Exec(ctx)
	if res.Err != nil {
		return apiFuncResult{nil, returnAPIError(res.Err), res.Warnings, qry.Close}
	}

	// Optional stats field in response if parameter "stats" is not empty.
	sr := api.statsRenderer
	if sr == nil {
		sr = defaultStatsRenderer
	}
	qs := sr(ctx, qry.Stats(), r.FormValue("stats"))

	return apiFuncResult{&QueryData{
		ResultType: res.Value.Type(),
		Result:     res.Value,
		Stats:      qs,
	}, nil, res.Warnings, qry.Close}
}
