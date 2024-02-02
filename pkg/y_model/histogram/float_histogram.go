package histogram

// FloatHistogram represents a histogram for float values.
type FloatHistogram struct {
	buckets map[float64]int
}

// NewFloatHistogram creates a new FloatHistogram.
func NewFloatHistogram() *FloatHistogram {
	return &FloatHistogram{buckets: make(map[float64]int)}
}

// Div divides each entry in the histogram by a scalar value.
func (h *FloatHistogram) Div(scalar float64) *FloatHistogram {
	result := NewFloatHistogram()
	for value, count := range h.buckets {
		result.buckets[value/scalar] = count
	}
	return result
}

// Add adds the values of another FloatHistogram to this one.
func (h *FloatHistogram) Add(other *FloatHistogram) *FloatHistogram {
	result := h.Copy()
	for value, count := range other.buckets {
		result.buckets[value] += count
	}
	return result
}

// Compact reduces the number of buckets in the histogram to a maximum number.
func (h *FloatHistogram) Compact(maxEmptyBuckets int) *FloatHistogram {
	// Implement the logic based on your specific compaction strategy
	// This could involve merging buckets, removing empty buckets, etc.
	return h
}

// Copy creates a deep copy of the FloatHistogram.
func (h *FloatHistogram) Copy() *FloatHistogram {
	newHistogram := NewFloatHistogram()
	for value, count := range h.buckets {
		newHistogram.buckets[value] = count
	}
	return newHistogram
}

// Size returns the number of buckets in the histogram.
func (h *FloatHistogram) Size() int {
	return len(h.buckets)
}
