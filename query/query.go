// Package query implements query matching against cards
package query

import (
	"mtgBuilder/card"
)

type Query interface {
	Matches(c *card.Card) bool
}
