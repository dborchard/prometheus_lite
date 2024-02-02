package promql

import (
	"fmt"
	"strconv"
	"strings"
)

func (m Matrix) Len() int { return len(m) }

func (m Matrix) Less(i, j int) bool {
	//TODO implement me
	panic("implement me")
}

func (m Matrix) Swap(i, j int) {
	//TODO implement me
	panic("implement me")
}

// Strings
func (m Matrix) String() string {
	strs := make([]string, len(m))

	for i, ss := range m {
		strs[i] = ss.String()
	}

	return strings.Join(strs, "\n")
}

func (v Vector) String() string {
	//TODO implement me
	panic("implement me")
}

func (s Scalar) String() string {
	v := strconv.FormatFloat(s.V, 'f', -1, 64)
	return fmt.Sprintf("scalar: %v @[%v]", v, s.T)
}

func (s String) String() string {
	return s.V
}
