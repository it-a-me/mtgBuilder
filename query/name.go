package query

import (
	"regexp"
	"strings"

	"mtgBuilder/card"
)

type Name struct {
	Name string
}

func (n Name) Matches(c *card.Card) bool {
	if strings.Contains(strings.ToLower(c.Name), n.Name) {
		return true
	}
	for _, face := range c.CardFaces {
		if strings.Contains(strings.ToLower(face.Name), n.Name) {
			return true
		}
	}
	return false
}

type NameExact struct {
	Name string
}

func (n NameExact) Matches(c *card.Card) bool {
	if strings.ToLower(c.Name) == n.Name {
		return true
	}
	for _, face := range c.CardFaces {
		if strings.ToLower(face.Name) == n.Name {
			return true
		}
	}
	return false
}

type NameRegex struct {
	Re *regexp.Regexp
}

func (o NameRegex) Matches(c *card.Card) bool {
	if o.Re.MatchString(c.Name) {
		return true
	}
	for _, face := range c.CardFaces {
		if o.Re.MatchString(face.Name) {
			return true
		}
	}
	return false
}
