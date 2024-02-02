package promql

import "time"

type QueryOpts interface {
	LookbackDelta() time.Duration
}

var _ QueryOpts = &PrometheusQueryOpts{}

type PrometheusQueryOpts struct {
	// Enables recording per-step statistics if the engine has it enabled as well. Disabled by default.
	enablePerStepStats bool
	// Lookback delta duration for this query.
	lookbackDelta time.Duration
}

func NewPrometheusQueryOpts(enablePerStepStats bool, lookbackDelta time.Duration) QueryOpts {
	return &PrometheusQueryOpts{
		enablePerStepStats: enablePerStepStats,
		lookbackDelta:      lookbackDelta,
	}
}
func (p *PrometheusQueryOpts) LookbackDelta() time.Duration {
	//TODO implement me
	panic("implement me")
}
