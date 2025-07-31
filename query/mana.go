package query

import (
	"fmt"
	"log/slog"

	"mtgBuilder/card"
)

type Mana struct {
	Colors card.Colors
}

func (m Mana) Matches(c *card.Card) bool {
	panic("todo")
}

type Manavalue struct {
	Relationship relationship
	Value        float32
}

func (m Manavalue) Matches(c *card.Card) bool {
	if c.Cmc == nil {
		return false
	}
	slog.Debug("Manavalue Matching", "value", m.Value, "cmc", *c.Cmc, "name", c.Name)
	switch m.Relationship {
	case Less:
		return m.Value < *c.Cmc
	case LessEqual:
		return m.Value <= *c.Cmc
	case Equal, Colon:
		return m.Value == *c.Cmc
	case GreaterEqual:
		return m.Value >= *c.Cmc
	case Greater:
		return m.Value > *c.Cmc
	}
	panic(fmt.Sprintf("Invalid relationship: %+v", m.Relationship))
}
