package query

import (
	"strings"

	"mtgBuilder/card"
)

type Set struct {
	Name string
}

func (s Set) Matches(c *card.Card) bool {
	return s.Name == strings.ToLower(c.Set)
}

type SetType struct {
	Name string
}

func (s SetType) Matches(c *card.Card) bool {
	return s.Name == strings.ToLower(c.SetType)
}
