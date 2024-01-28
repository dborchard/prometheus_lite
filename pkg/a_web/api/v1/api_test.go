package v1

import (
	"context"
	promql "github.com/dborchard/prometheus_lite/pkg/b_promql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAPI_queryRange(t *testing.T) {
	api := NewAPI(promql.NewEngine(nil), nil)

	// 1. Args
	ctx := context.Background()
	start := time.Time{}
	end := time.Time{}
	step, _ := time.ParseDuration("1")
	qry, _ := api.QueryEngine.NewRangeQuery(ctx, api.Queryable, nil, "query", start, end, step)

	res := qry.Exec(ctx)
	assert.Nil(t, res.Err)
}
