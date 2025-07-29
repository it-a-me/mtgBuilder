package query

import (
	"regexp"
	"strings"

	"mtgBuilder/card"
)

type FullOracleText struct {
	Substr string
}

func (o FullOracleText) Matches(c *card.Card) bool {
	if c.OracleText != nil {
		if strings.Contains(strings.ToLower(*c.OracleText), o.Substr) {
			return true
		}
	}
	for _, face := range c.CardFaces {
		if face.OracleText != nil {
			if strings.Contains(strings.ToLower(*face.OracleText), o.Substr) {
				return true
			}
		}
	}
	return false
}

type OracleText struct {
	Substr string
}

func (o OracleText) Matches(c *card.Card) bool {
	if c.OracleText != nil {
		t := strings.ToLower(StripParens.ReplaceAllLiteralString(*c.OracleText, ""))
		if strings.Contains(t, o.Substr) {
			return true
		}
	}
	for _, face := range c.CardFaces {
		if face.OracleText != nil {
			t := strings.ToLower(StripParens.ReplaceAllLiteralString(*face.OracleText, ""))
			if strings.Contains(t, o.Substr) {
				return true
			}
		}
	}
	return false
}

type OracleTextRegex struct {
	Re *regexp.Regexp
}

func (o OracleTextRegex) Matches(c *card.Card) bool {
	if c.OracleText != nil {
		t := StripParens.ReplaceAllLiteralString(*c.OracleText, "")
		if o.Re.MatchString(t) {
			return true
		}
	}
	for _, face := range c.CardFaces {
		if face.OracleText != nil {
			t := StripParens.ReplaceAllLiteralString(*face.OracleText, "")
			if o.Re.MatchString(t) {
				return true
			}
		}
	}
	return false
}

type FullOracleTextRegex struct {
	Re *regexp.Regexp
}

func (o FullOracleTextRegex) Matches(c *card.Card) bool {
	if c.OracleText != nil {
		if o.Re.MatchString(*c.OracleText) {
			return true
		}
	}
	for _, face := range c.CardFaces {
		if face.OracleText != nil {
			if o.Re.MatchString(*face.OracleText) {
				return true
			}
		}
	}
	return false
}

type Keyword struct {
	Word string
}

var StripParens = regexp.MustCompile(`\(.*?\)`)
