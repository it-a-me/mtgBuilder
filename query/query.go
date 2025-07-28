// Package query implements query matching against cards
package query

import (
	"regexp"
	"slices"
	"strings"

	"mtgBuilder/card"
)

type Query interface {
	Matches(c *card.Card) bool
}

type OracleText struct {
	Re *regexp.Regexp
}

var StripParens = regexp.MustCompile(`\(.*?\)`)

func (o OracleText) Matches(c *card.Card) bool {
	lines := c.GetOracleText()
	for _, line := range lines {
		cleaned := StripParens.ReplaceAllLiteralString(line, "")
		if o.Re.MatchString(cleaned) {
			return true
		}
	}
	return false
}

// NewBasicMatcher returns a new matcher with multiline mode and case insensitivity enabled
func NewBasicMatcher(content string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(content)
	return regexp.Compile(`(?im)` + escaped)
}

// NewRegexMatcher returns a new matcher with multiline mode and case insensitivity enabled
func NewRegexMatcher(content string) (*regexp.Regexp, error) {
	// return regexp.Compile(`(?im)` + content)
	return regexp.Compile(`(?m)` + content)
}

type Name struct {
	Name  string
	Exact bool
}

func NewNameQuery(name string, exact bool) Query {
	return Name{name, exact}
}

func (n Name) Matches(c *card.Card) bool {
	comp := func(s string) bool { return strings.Contains(s, n.Name) }
	if n.Exact {
		comp = func(s string) bool { return s == n.Name }
	}
	if comp(c.Name) {
		return true
	}
	return slices.ContainsFunc(c.CardFaces, func(c card.CardFace) bool { return comp(c.Name) })
}
