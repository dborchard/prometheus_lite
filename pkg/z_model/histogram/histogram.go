package histogram

// Histogram encodes a sparse, high-resolution histogram. See the design
// document for full details:
// https://docs.google.com/document/d/1cLNv3aufPZb3fNfaJgdaRBZsInZKKIHo9E6HinJVbpM/edit#
//
// The most tricky bit is how bucket indices represent real bucket boundaries.
// An example for schema 0 (by which each bucket is twice as wide as the
// previous bucket):
//
//	Bucket boundaries →              [-2,-1)  [-1,-0.5) [-0.5,-0.25) ... [-0.001,0.001] ... (0.25,0.5] (0.5,1]  (1,2] ....
//	                                    ↑        ↑           ↑                  ↑                ↑         ↑      ↑
//	Zero bucket (width e.g. 0.001) →    |        |           |                  ZB               |         |      |
//	Positive bucket indices →           |        |           |                          ...     -1         0      1    2    3
//	Negative bucket indices →  3   2    1        0          -1       ...
//
// Which bucket indices are actually used is determined by the spans.
type Histogram struct {
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
	// Observations falling into the zero bucket.
	ZeroCount uint64
	// Total number of observations.
	Count uint64
	// Sum of observations. This is also used as the stale marker.
	Sum float64
	// Spans for positive and negative buckets (see Span below).
	//PositiveSpans, NegativeSpans []Span
	// Observation counts in buckets. The first element is an absolute
	// count. All following ones are deltas relative to the previous
	// element.
	PositiveBuckets, NegativeBuckets []int64
}
