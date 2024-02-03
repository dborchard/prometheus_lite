package storage

import (
	"context"
	"github.com/dborchard/prometheus_lite/pkg/y_model/labels"
)

// LabelQuerier provides querying access over labels.
type LabelQuerier interface {
	// LabelValues returns all potential values for a label name.
	// It is not safe to use the strings beyond the lifetime of the querier.
	// If matchers are specified the returned result set is reduced
	// to label values of metrics matching the matchers.
	LabelValues(ctx context.Context, name string, matchers ...*labels.Matcher) ([]string, error)

	// LabelNames returns all the unique label names present in the block in sorted order.
	// If matchers are specified the returned result set is reduced
	// to label names of metrics matching the matchers.
	LabelNames(ctx context.Context, matchers ...*labels.Matcher) ([]string, error)

	// Close releases the resources of the Querier.
	Close() error
}
