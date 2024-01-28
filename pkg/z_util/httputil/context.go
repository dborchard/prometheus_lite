package httputil

import (
	"context"
	promql "github.com/dborchard/prometheus_lite/pkg/b_promql"
	"net/http"
)

func ContextFromRequest(ctx context.Context, r *http.Request) context.Context {
	var ip string
	var path string
	return promql.NewOriginContext(ctx, map[string]interface{}{
		"httpRequest": map[string]string{
			"clientIP": ip,
			"method":   r.Method,
			"path":     path,
		},
	})
}
