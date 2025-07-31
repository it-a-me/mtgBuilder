package query

import "mtgBuilder/card"

type Format struct {
	Format   string
	Expected string
}

var formatAliases = map[string]string{
	"s": "standard",
	"f": "future",
	"h": "historic",
	"t": "timeless",
	"g": "gladiator",
	// "p": "pioneer",
	"m": "modern",
	"l": "legacy",
	"p": "pauper",
	"v": "vintage",
	// "p": "penny",
	"c":    "commander",
	"edh":  "commander",
	"o":    "oathbreaker",
	"sb":   "standardbrawl",
	"b":    "brawl",
	"a":    "alchemy",
	"pedh": "paupercommander",
	"d":    "duel",
	"os":   "oldschool",
	"pre":  "premodern",
	// "p": "predh",
}

func (f Format) Matches(c *card.Card) bool {
	status := c.Legalities[f.Format]
	return status == f.Expected
}
