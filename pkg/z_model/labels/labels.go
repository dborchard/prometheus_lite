package labels

import (
	"slices"
	"strings"
)

const (
	MetricName   = "__name__"
	AlertName    = "alertname"
	BucketLabel  = "le"
	InstanceName = "instance"

	labelSep = '\xfe'
)

// Labels is a sorted set of labels. Order has to be guaranteed upon
// instantiation.
type Labels []Label

// Hash returns a hash value for the label set.
// Note: the result is not guaranteed to be consistent across different runs of Prometheus.
func (ls Labels) Hash() uint64 {
	//return labels.Hash(ls)
	return 0
}

// Range calls f on each label.
func (ls Labels) Range(f func(l Label)) {
	for _, l := range ls {
		f(l)
	}
}

// Label is a key/value pair of strings.
type Label struct {
	Name, Value string
}

// Builder allows modifying Labels.
type Builder struct {
	base Labels
	del  []string
	add  []Label
}

// NewBuilder returns a new LabelsBuilder.
func NewBuilder(base Labels) *Builder {
	b := &Builder{
		del: make([]string, 0, 5),
		add: make([]Label, 0, 5),
	}
	b.Reset(base)
	return b
}

// Reset clears all current state for the builder.
func (b *Builder) Reset(base Labels) {
	b.base = base
	b.del = b.del[:0]
	b.add = b.add[:0]
	b.base.Range(func(l Label) {
		if l.Value == "" {
			b.del = append(b.del, l.Name)
		}
	})
}

// Set the name/value pair as a label. A value of "" means delete that label.
func (b *Builder) Set(n, v string) *Builder {
	if v == "" {
		// Empty labels are the same as missing labels.
		return b.Del(n)
	}
	for i, a := range b.add {
		if a.Name == n {
			b.add[i].Value = v
			return b
		}
	}
	b.add = append(b.add, Label{Name: n, Value: v})

	return b
}

// Del deletes the label of the given name.
func (b *Builder) Del(ns ...string) *Builder {
	for _, n := range ns {
		for i, a := range b.add {
			if a.Name == n {
				b.add = append(b.add[:i], b.add[i+1:]...)
			}
		}
		b.del = append(b.del, n)
	}
	return b
}

func (b *Builder) Labels() Labels {
	return b.base
}

// EmptyLabels returns n empty Labels value, for convenience.
func EmptyLabels() Labels {
	return Labels{}
}

// New returns a sorted Labels from the given labels.
// The caller has to guarantee that all label names are unique.
func New(ls ...Label) Labels {
	set := make(Labels, 0, len(ls))
	set = append(set, ls...)
	slices.SortFunc(set, func(a, b Label) int { return strings.Compare(a.Name, b.Name) })

	return set
}

// FromStrings creates new labels from pairs of strings.
func FromStrings(ss ...string) Labels {
	if len(ss)%2 != 0 {
		panic("invalid number of strings")
	}
	res := make(Labels, 0, len(ss)/2)
	for i := 0; i < len(ss); i += 2 {
		res = append(res, Label{Name: ss[i], Value: ss[i+1]})
	}

	slices.SortFunc(res, func(a, b Label) int { return strings.Compare(a.Name, b.Name) })
	return res
}
