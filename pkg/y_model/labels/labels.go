package labels

import (
	"bytes"
	"github.com/cespare/xxhash/v2"
)

const (
	MetricName   = "__name__"
	AlertName    = "alertname"
	BucketLabel  = "le"
	InstanceName = "instance"

	labelSep = '\xfe'
)

var seps = []byte{'\xff'}

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

func (ls Labels) HashForLabels(names ...string) uint64 {
	b := make([]byte, 0, 256)
	i, j := 0, 0
	for i < len(ls) && j < len(names) {
		switch {
		case names[j] < ls[i].Name:
			j++
		case ls[i].Name < names[j]:
			i++
		default:
			b = append(b, ls[i].Name...)
			b = append(b, seps[0])
			b = append(b, ls[i].Value...)
			b = append(b, seps[0])
			i++
			j++
		}
	}
	return xxhash.Sum64(b)
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

// BytesWithoutLabels is just as Bytes(), but only for labels not matching names.
// 'names' have to be sorted in ascending order.
func (ls Labels) BytesWithoutLabels(buf []byte, names ...string) []byte {
	b := bytes.NewBuffer(buf[:0])
	b.WriteByte(labelSep)
	j := 0
	for i := range ls {
		for j < len(names) && names[j] < ls[i].Name {
			j++
		}
		if j < len(names) && ls[i].Name == names[j] {
			continue
		}
		if b.Len() > 1 {
			b.WriteByte(seps[0])
		}
		b.WriteString(ls[i].Name)
		b.WriteByte(seps[0])
		b.WriteString(ls[i].Value)
	}
	return b.Bytes()
}

func (ls Labels) Bytes(buf []byte) []byte {
	return ls.BytesWithoutLabels(buf, "")
}

// Get returns the value for the label with the given name.
// Returns an empty string if the label doesn't exist.
func (ls Labels) Get(name string) string {
	for _, l := range ls {
		if l.Name == name {
			return l.Value
		}
	}
	return ""
}
