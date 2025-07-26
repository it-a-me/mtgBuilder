// Package query implements query matching against cards
package query

import (
	"regexp"

	"mtgBuilder/card"
)

type Query interface {
	Matches(c *card.Card) bool
}

type OracleText struct {
	Re regexp.Regexp
}

// NewBasicMatcher returns a new matcher with multiline mode and case insensitivity enabled
func NewBasicMatcher(content string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(content)
	return regexp.Compile(`(?im)` + escaped)
}

// NewRegexMatcher returns a new matcher with multiline mode and case insensitivity enabled
func NewRegexMatcher(content string) (*regexp.Regexp, error) {
	return regexp.Compile(`(?im)` + content)
}

var StripParens = regexp.MustCompile(`\(.*?\)`)

func (o OracleText) Matches(c *card.Card) bool {
	lines := c.GetOracleText()
	for _, line := range lines {
		cleaned := StripParens.ReplaceAllLiteralString(line, "")
		if StripParens.MatchString(cleaned) {
			return true
		}
	}
	return false
}
