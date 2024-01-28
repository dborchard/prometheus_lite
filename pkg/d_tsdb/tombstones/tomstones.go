package tombstones

import storage "github.com/dborchard/prometheus_lite/pkg/c_storage"

// Reader gives access to tombstone intervals by series reference.
type Reader interface {
	// Get returns deletion intervals for the series with the given reference.
	Get(ref storage.SeriesRef) (Intervals, error)

	// Iter calls the given function for each encountered interval.
	Iter(func(storage.SeriesRef, Intervals) error) error

	// Total returns the total count of tombstones.
	Total() uint64

	// Close any underlying resources
	Close() error
}
