package query

import (
	"strings"

	"mtgBuilder/card"
)

type Type struct {
	Text string
}

func (t Type) Matches(c *card.Card) bool {
	if strings.Contains(strings.ToLower(c.TypeLine), t.Text) {
		return true
	}
	for _, face := range c.CardFaces {
		if face.TypeLine != nil && strings.Contains(strings.ToLower(c.TypeLine), t.Text) {
			return true
		}
	}
	return false
}
