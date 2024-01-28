package promql

import (
	"github.com/dborchard/prometheus_lite/pkg/b_promql/parser"
	"time"
)

func timeMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond/time.Nanosecond)
}
func durationMilliseconds(d time.Duration) int64 {
	return int64(d / (time.Millisecond / time.Nanosecond))
}
func getTimeRangesForSelector(s *parser.EvalStmt, n *parser.VectorSelector, path []parser.Node, evalRange time.Duration) (int64, int64) {
	start := timeMilliseconds(s.Start)
	end := timeMilliseconds(s.End)
	if evalRange > 0 {
		start = end - durationMilliseconds(evalRange)
	}
	return start, end
}
func (ng *Engine) getLastSubqueryInterval(path []parser.Node) time.Duration {
	var interval time.Duration
	//for _, node := range path {
	//	//if n, ok := node.(*parser.SubqueryExpr); ok {
	//	//	interval = n.Step
	//	//	if n.Step == 0 {
	//	//		interval = time.Duration(ng.noStepSubqueryIntervalFn(durationMilliseconds(n.Range))) * time.Millisecond
	//	//	}
	//	//}
	//}
	return interval
}

// extractGroupsFromPath parses vector outer function and extracts grouping information if by or without was used.
func extractGroupsFromPath(p []parser.Node) (bool, []string) {
	if len(p) == 0 {
		return false, nil
	}
	if n, ok := p[len(p)-1].(*parser.AggregateExpr); ok {
		return !n.Without, n.Grouping
	}
	return false, nil
}
