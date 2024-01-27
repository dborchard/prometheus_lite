package promql

import (
	"context"
	"time"
)

// QueryTracker provides access to two features:
//
// 1) Tracking of active query. If PromQL engine crashes while executing any query, such query should be present
// in the tracker on restart, hence logged. After the logging on restart, the tracker gets emptied.
//
// 2) Enforcement of the maximum number of concurrent queries.
type QueryTracker interface {
	// GetMaxConcurrent returns maximum number of concurrent queries that are allowed by this tracker.
	GetMaxConcurrent() int

	// Insert inserts query into query tracker. This call must block if maximum number of queries is already running.
	// If Insert doesn't return error then returned integer value should be used in subsequent Delete call.
	// Insert should return error if context is finished before query can proceed, and integer value returned in this case should be ignored by caller.
	Insert(ctx context.Context, query string) (int, error)

	// Delete removes query from activity tracker. InsertIndex is value returned by Insert call.
	Delete(insertIndex int)
}

type QueryOpts interface {
}

type PrometheusQueryOpts struct {
	// Enables recording per-step statistics if the engine has it enabled as well. Disabled by default.
	enablePerStepStats bool
	// Lookback delta duration for this query.
	lookbackDelta time.Duration
}

var _ QueryOpts = &PrometheusQueryOpts{}

func NewPrometheusQueryOpts(enablePerStepStats bool, lookbackDelta time.Duration) QueryOpts {
	return &PrometheusQueryOpts{
		enablePerStepStats: enablePerStepStats,
		lookbackDelta:      lookbackDelta,
	}
}

type QueryOrigin struct{}

// NewOriginContext returns a new context with data about the origin attached.
func NewOriginContext(ctx context.Context, data map[string]interface{}) context.Context {
	return context.WithValue(ctx, QueryOrigin{}, data)
}
