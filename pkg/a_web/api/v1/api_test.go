package v1

import (
	"context"
	promql "github.com/dborchard/prometheus_lite/pkg/b_promql"
	tsdb "github.com/dborchard/prometheus_lite/pkg/d_tsdb"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAPI_queryRange(t *testing.T) {
	api := NewAPI(promql.NewEngine(nil), tsdb.NewDB())

	// 1. Args
	ctx := context.Background()
	qry, _ := api.QueryEngine.NewInstantQuery(ctx, api.Queryable, nil, "query", time.Time{})

	res := qry.Exec(ctx)
	assert.Nil(t, res.Err)
	assert.NotNil(t, res.Value)
	//assert.Equal(t, "[] =>\n2 @[-6795364577871]", res.Value.String())
}
