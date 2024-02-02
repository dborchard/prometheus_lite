package promql

import (
	"github.com/dborchard/prometheus_lite/pkg/y_model/histogram"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
	"regexp"
)

// EvalNodeHelper stores extra information and caches for evaluating a single node across steps.
type EvalNodeHelper struct {
	// Evaluation timestamp.
	Ts int64
	// Vector that can be used for output.
	Out Vector

	// Caches.
	// DropMetricName and label_*.
	Dmn map[uint64]labels.Labels
	// funcHistogramQuantile for classic histograms.
	//signatureToMetricWithBuckets map[string]*metricWithBuckets
	// label_replace.
	regex *regexp.Regexp

	lb           *labels.Builder
	lblBuf       []byte
	lblResultBuf []byte

	// For binary vector matching.
	rightSigs    map[string]Sample
	matchedSigs  map[string]map[uint64]struct{}
	resultMetric map[string]labels.Labels
}

// EvalSeriesHelper stores extra information about a series.
type EvalSeriesHelper struct {
	// The grouping key used by aggregation.
	groupingKey uint64
	// Used to map left-hand to right-hand in binary operations.
	signature string
}

type groupedAggregation struct {
	hasFloat       bool // Has at least 1 float64 sample aggregated.
	hasHistogram   bool // Has at least 1 histogram sample aggregated.
	labels         labels.Labels
	floatValue     float64
	histogramValue *histogram.FloatHistogram
	floatMean      float64
	histogramMean  *histogram.FloatHistogram
	groupCount     int
}
