package promql

import (
	storage "github.com/dborchard/prometheus_lite/pkg/c_storage"
	"github.com/dborchard/prometheus_lite/pkg/d_tsdb/chunkenc"
	"github.com/dborchard/prometheus_lite/pkg/y_model/value"
)

// matrixIterSlice populates a matrix vector covering the requested range for a
// single time series, with points retrieved from an iterator.
//
// As an optimization, the matrix vector may already contain points of the same
// time series from the evaluation of an earlier step (with lower mint and maxt
// values). Any such points falling before mint are discarded; points that fall
// into the [mint, maxt] range are retained; only points with later timestamps
// are populated from the iterator.
func (ev *evaluator) matrixIterSlice(
	it *storage.BufferedSeriesIterator, mint, maxt int64,
	floats []FPoint, histograms []HPoint,
) ([]FPoint, []HPoint) {
	mintFloats, mintHistograms := mint, mint

	// First floats...
	if len(floats) > 0 && floats[len(floats)-1].T >= mint {
		// There is an overlap between previous and current ranges, retain common
		// points. In most such cases:
		//   (a) the overlap is significantly larger than the eval step; and/or
		//   (b) the number of samples is relatively small.
		// so a linear search will be as fast as a binary search.
		var drop int
		for drop = 0; floats[drop].T < mint; drop++ {
		}
		ev.currentSamples -= drop
		copy(floats, floats[drop:])
		floats = floats[:len(floats)-drop]
		// Only append points with timestamps after the last timestamp we have.
		mintFloats = floats[len(floats)-1].T + 1
	} else {
		ev.currentSamples -= len(floats)
		if floats != nil {
			floats = floats[:0]
		}
	}

	// ...then the same for histograms. TODO(beorn7): Use generics?
	if len(histograms) > 0 && histograms[len(histograms)-1].T >= mint {
		// There is an overlap between previous and current ranges, retain common
		// points. In most such cases:
		//   (a) the overlap is significantly larger than the eval step; and/or
		//   (b) the number of samples is relatively small.
		// so a linear search will be as fast as a binary search.
		var drop int
		for drop = 0; histograms[drop].T < mint; drop++ {
		}
		copy(histograms, histograms[drop:])
		histograms = histograms[:len(histograms)-drop]
		ev.currentSamples -= totalHPointSize(histograms)
		// Only append points with timestamps after the last timestamp we have.
		mintHistograms = histograms[len(histograms)-1].T + 1
	} else {
		ev.currentSamples -= totalHPointSize(histograms)
		if histograms != nil {
			histograms = histograms[:0]
		}
	}

	soughtValueType := it.Seek(maxt)
	if soughtValueType == chunkenc.ValNone {
		if it.Err() != nil {
			panic(it.Err())
		}
	}

	buf := it.Buffer()
loop:
	for {
		switch buf.Next() {
		case chunkenc.ValNone:
			break loop
		case chunkenc.ValFloatHistogram, chunkenc.ValHistogram:
			t, h := buf.AtFloatHistogram()
			// Values in the buffer are guaranteed to be smaller than maxt.
			if t >= mintHistograms {
				if ev.currentSamples >= ev.maxSamples {
					panic("too many samples")
				}
				point := HPoint{T: t, H: h}
				if histograms == nil {
					histograms = getHPointSlice(16)
				}
				histograms = append(histograms, point)
				ev.currentSamples += point.size()
			}
		case chunkenc.ValFloat:
			t, f := buf.At()
			if value.IsStaleNaN(f) {
				continue loop
			}
			// Values in the buffer are guaranteed to be smaller than maxt.
			if t >= mintFloats {
				if ev.currentSamples >= ev.maxSamples {
					panic("too many samples")
				}
				ev.currentSamples++
				if floats == nil {
					floats = getFPointSlice(16)
				}
				floats = append(floats, FPoint{T: t, F: f})
			}
		}
	}
	// The sought sample might also be in the range.
	switch soughtValueType {
	case chunkenc.ValFloatHistogram, chunkenc.ValHistogram:
		t, h := it.AtFloatHistogram()
		if t == maxt {
			if ev.currentSamples >= ev.maxSamples {
				panic("too many samples")
			}
			if histograms == nil {
				histograms = getHPointSlice(16)
			}
			point := HPoint{T: t, H: h}
			histograms = append(histograms, point)
			ev.currentSamples += point.size()
		}
	case chunkenc.ValFloat:
		t, f := it.At()
		if t == maxt {
			if ev.currentSamples >= ev.maxSamples {
				panic("too many samples")
			}
			if floats == nil {
				floats = getFPointSlice(16)
			}
			floats = append(floats, FPoint{T: t, F: f})
			ev.currentSamples++
		}
	default:
		panic("unhandled default case")
	}
	return floats, histograms
}
