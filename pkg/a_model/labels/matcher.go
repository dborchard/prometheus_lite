package labels

import "regexp"

// Matcher models the matching of a label.
type Matcher struct {
	Type  MatchType
	Name  string
	Value string

	re *FastRegexMatcher
}

// MatchType is an enum for label matching types.
type MatchType int

// Possible MatchTypes.
const (
	MatchEqual MatchType = iota
	MatchNotEqual
	MatchRegexp
	MatchNotRegexp
)

type FastRegexMatcher struct {
	re       *regexp.Regexp
	prefix   string
	suffix   string
	contains string

	// shortcut for literals
	literal bool
	value   string
}
