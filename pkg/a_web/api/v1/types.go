package v1

import (
	"context"
	parser "github.com/dborchard/prometheus_lite/pkg/b_parser"
	promql "github.com/dborchard/prometheus_lite/pkg/b_promql"
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"time"
)

type apiFuncResult struct {
	data      *QueryData
	err       *apiError
	finalizer func()
}

type QueryData struct {
	ResultType parser.ValueType `json:"resultType"`
	Result     parser.Value     `json:"result"`
}

type apiError struct {
	typ errorType
	err error
}

type errorType string

const (
	errorNone          errorType = ""
	errorTimeout       errorType = "timeout"
	errorCanceled      errorType = "canceled"
	errorExec          errorType = "execution"
	errorBadData       errorType = "bad_data"
	errorInternal      errorType = "internal"
	errorUnavailable   errorType = "unavailable"
	errorNotFound      errorType = "not_found"
	errorNotAcceptable errorType = "not_acceptable"
)

// QueryEngine defines the interface for the *promql.Engine, so it can be replaced, wrapped or mocked.
type QueryEngine interface {
	NewInstantQuery(ctx context.Context, q storage.Queryable, opts promql.QueryOpts, qs string, ts time.Time) (promql.Query, error)
	NewRangeQuery(ctx context.Context, q storage.Queryable, opts promql.QueryOpts, qs string, start, end time.Time, interval time.Duration) (promql.Query, error)
}
