package query

import (
	"mtgBuilder/card"
)

type Color struct {
	Mulicolor bool
	Operator  relationship
	Colors    card.Colors
}

func (q Color) Matches(c *card.Card) bool {
	var colors card.Colors
	if c.Colors != nil {
		colors.Add(*c.Colors)
	}
	for _, face := range c.CardFaces {
		if face.Colors != nil {
			colors.Add(*face.Colors)
		}
	}

	if q.Mulicolor {
		return len(colors) > 1
	}

	switch q.Operator {
	case Less:
		return colors.IsSubset(q.Colors) && len(colors) < len(q.Colors)
	case LessEqual:
		return colors.IsSubset(q.Colors)
	case Equal:
		return colors.Equal(q.Colors)
	case GreaterEqual:
		return q.Colors.IsSubset(colors)
	case Greater:
		return q.Colors.IsSubset(colors) && len(colors) > len(q.Colors)
	}
	panic("unreachable")
}

type ColorIdentity struct {
	Mulicolor bool
	Operator  relationship
	Colors    card.Colors
}

func (q ColorIdentity) Matches(c *card.Card) bool {
	if c.ColorIdentity == nil {
		return false
	}
	colors := *c.ColorIdentity
	if q.Mulicolor {
		return len(colors) > 1
	}

	switch q.Operator {
	case Less:
		return colors.IsSubset(q.Colors) && len(colors) < len(q.Colors)
	case LessEqual:
		return colors.IsSubset(q.Colors)
	case Equal:
		return colors.Equal(q.Colors)
	case GreaterEqual:
		return q.Colors.IsSubset(colors)
	case Greater:
		return q.Colors.IsSubset(colors) && len(colors) > len(q.Colors)
	}
	panic("unreachable")
}
