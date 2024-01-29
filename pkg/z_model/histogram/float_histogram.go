package histogram

import (
	"fmt"
	"math"
	"strings"
)

// FloatHistogram is similar to Histogram but uses float64 for all
// counts. Additionally, bucket counts are absolute and not deltas.
//
// A FloatHistogram is needed by PromQL to handle operations that might result
// in fractional counts. Since the counts in a histogram are unlikely to be too
// large to be represented precisely by a float64, a FloatHistogram can also be
// used to represent a histogram with integer counts and thus serves as a more
// generalized representation.
type FloatHistogram struct {
	// Counter reset information.
	//CounterResetHint CounterResetHint
	// Currently valid schema numbers are -4 <= n <= 8.  They are all for
	// base-2 bucket schemas, where 1 is a bucket boundary in each case, and
	// then each power of two is divided into 2^n logarithmic buckets.  Or
	// in other words, each bucket boundary is the previous boundary times
	// 2^(2^-n).
	Schema int32
	// Width of the zero bucket.
	ZeroThreshold float64
	// Observations falling into the zero bucket. Must be zero or positive.
	ZeroCount float64
	// Total number of observations. Must be zero or positive.
	Count float64
	// Sum of observations. This is also used as the stale marker.
	Sum float64
	// Spans for positive and negative buckets (see Span below).
	//PositiveSpans, NegativeSpans []Span
	// Observation counts in buckets. Each represents an absolute count and
	// must be zero or positive.
	PositiveBuckets, NegativeBuckets []float64

	PositiveSpans, NegativeSpans []Span
}

// TestExpression returns the string representation of this histogram as it is used in the internal PromQL testing
// framework as well as in promtool rules unit tests.
// The syntax is described in https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/#series
func (h *FloatHistogram) TestExpression() string {
	var res []string
	m := h.Copy()

	m.Compact(math.MaxInt) // Compact to reduce the number of positive and negative spans to 1.

	if m.Schema != 0 {
		res = append(res, fmt.Sprintf("schema:%d", m.Schema))
	}
	if m.Count != 0 {
		res = append(res, fmt.Sprintf("count:%g", m.Count))
	}
	if m.Sum != 0 {
		res = append(res, fmt.Sprintf("sum:%g", m.Sum))
	}
	if m.ZeroCount != 0 {
		res = append(res, fmt.Sprintf("z_bucket:%g", m.ZeroCount))
	}
	if m.ZeroThreshold != 0 {
		res = append(res, fmt.Sprintf("z_bucket_w:%g", m.ZeroThreshold))
	}

	addBuckets := func(kind, bucketsKey, offsetKey string, buckets []float64, spans []Span) []string {
		if len(spans) > 1 {
			panic(fmt.Sprintf("histogram with multiple %s spans not supported", kind))
		}
		for _, span := range spans {
			if span.Offset != 0 {
				res = append(res, fmt.Sprintf("%s:%d", offsetKey, span.Offset))
			}
		}

		var bucketStr []string
		for _, bucket := range buckets {
			bucketStr = append(bucketStr, fmt.Sprintf("%g", bucket))
		}
		if len(bucketStr) > 0 {
			res = append(res, fmt.Sprintf("%s:[%s]", bucketsKey, strings.Join(bucketStr, " ")))
		}
		return res
	}
	res = addBuckets("positive", "buckets", "offset", m.PositiveBuckets, m.PositiveSpans)
	res = addBuckets("negative", "n_buckets", "n_offset", m.NegativeBuckets, m.NegativeSpans)
	return "{{" + strings.Join(res, " ") + "}}"
}

// Copy returns a deep copy of the Histogram.
func (h *FloatHistogram) Copy() *FloatHistogram {
	c := *h
	if h.PositiveSpans != nil {
		c.PositiveSpans = make([]Span, len(h.PositiveSpans))
		copy(c.PositiveSpans, h.PositiveSpans)
	}
	if h.NegativeSpans != nil {
		c.NegativeSpans = make([]Span, len(h.NegativeSpans))
		copy(c.NegativeSpans, h.NegativeSpans)
	}
	if h.PositiveBuckets != nil {
		c.PositiveBuckets = make([]float64, len(h.PositiveBuckets))
		copy(c.PositiveBuckets, h.PositiveBuckets)
	}
	if h.NegativeBuckets != nil {
		c.NegativeBuckets = make([]float64, len(h.NegativeBuckets))
		copy(c.NegativeBuckets, h.NegativeBuckets)
	}

	return &c
}

// A Span defines a continuous sequence of buckets.
type Span struct {
	// Gap to previous span (always positive), or starting index for the 1st
	// span (which can be negative).
	Offset int32
	// Length of the span.
	Length uint32
}

// Compact eliminates empty buckets at the beginning and end of each span, then
// merges spans that are consecutive or at most maxEmptyBuckets apart, and
// finally splits spans that contain more consecutive empty buckets than
// maxEmptyBuckets. (The actual implementation might do something more efficient
// but with the same result.)  The compaction happens "in place" in the
// receiving histogram, but a pointer to it is returned for convenience.
//
// The ideal value for maxEmptyBuckets depends on circumstances. The motivation
// to set maxEmptyBuckets > 0 is the assumption that is less overhead to
// represent very few empty buckets explicitly within one span than cutting the
// one span into two to treat the empty buckets as a gap between the two spans,
// both in terms of storage requirement as well as in terms of encoding and
// decoding effort. However, the tradeoffs are subtle. For one, they are
// different in the exposition format vs. in a TSDB chunk vs. for the in-memory
// representation as Go types. In the TSDB, as an additional aspects, the span
// layout is only stored once per chunk, while many histograms with that same
// chunk layout are then only stored with their buckets (so that even a single
// empty bucket will be stored many times).
//
// For the Go types, an additional Span takes 8 bytes. Similarly, an additional
// bucket takes 8 bytes. Therefore, with a single separating empty bucket, both
// options have the same storage requirement, but the single-span solution is
// easier to iterate through. Still, the safest bet is to use maxEmptyBuckets==0
// and only use a larger number if you know what you are doing.
func (h *FloatHistogram) Compact(maxEmptyBuckets int) *FloatHistogram {

	return h
}

// String returns a string representation of the Histogram.
func (h *FloatHistogram) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{count:%g, sum:%g", h.Count, h.Sum)

	//var nBuckets []Bucket[float64]
	//for it := h.NegativeBucketIterator(); it.Next(); {
	//	bucket := it.At()
	//	if bucket.Count != 0 {
	//		nBuckets = append(nBuckets, it.At())
	//	}
	//}
	//for i := len(nBuckets) - 1; i >= 0; i-- {
	//	fmt.Fprintf(&sb, ", %s", nBuckets[i].String())
	//}
	//
	//if h.ZeroCount != 0 {
	//	fmt.Fprintf(&sb, ", %s", h.ZeroBucket().String())
	//}
	//
	//for it := h.PositiveBucketIterator(); it.Next(); {
	//	bucket := it.At()
	//	if bucket.Count != 0 {
	//		fmt.Fprintf(&sb, ", %s", bucket.String())
	//	}
	//}

	sb.WriteRune('}')
	return sb.String()
}

// Div works like Mul but divides instead of multiplies.
// When dividing by 0, everything will be set to Inf.
func (h *FloatHistogram) Div(scalar float64) *FloatHistogram {
	h.ZeroCount /= scalar
	h.Count /= scalar
	h.Sum /= scalar
	for i := range h.PositiveBuckets {
		h.PositiveBuckets[i] /= scalar
	}
	for i := range h.NegativeBuckets {
		h.NegativeBuckets[i] /= scalar
	}
	return h
}
