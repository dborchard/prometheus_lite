package promql

import "context"

type QueryOrigin struct{}

// NewOriginContext returns a new context with data about the origin attached.
func NewOriginContext(ctx context.Context, data map[string]interface{}) context.Context {
	return context.WithValue(ctx, QueryOrigin{}, data)
}
