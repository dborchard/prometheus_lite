package labels

import (
	"regexp"
	"regexp/syntax"
)

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

// NewMatcher returns a matcher object.
func NewMatcher(t MatchType, n, v string) (*Matcher, error) {
	m := &Matcher{
		Type:  t,
		Name:  n,
		Value: v,
	}
	if t == MatchRegexp || t == MatchNotRegexp {
		re, err := NewFastRegexMatcher(v)
		if err != nil {
			return nil, err
		}
		m.re = re
	}
	return m, nil
}

func NewFastRegexMatcher(v string) (*FastRegexMatcher, error) {
	//if isLiteral(v) {
	//	return &FastRegexMatcher{literal: true, value: v}, nil
	//}
	re, err := regexp.Compile("^(?:" + v + ")$")
	if err != nil {
		return nil, err
	}

	parsed, err := syntax.Parse(v, syntax.Perl)
	if err != nil {
		return nil, err
	}

	m := &FastRegexMatcher{
		re: re,
	}

	if parsed.Op == syntax.OpConcat {
		//m.prefix, m.suffix, m.contains = optimizeConcatRegex(parsed)
	}

	return m, nil
}
