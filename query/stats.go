package query

import (
	"cmp"
	"fmt"
	"strconv"

	"mtgBuilder/card"
)

func fieldCompare[T cmp.Ordered](a T, rel relationship, b T) bool {
	c := cmp.Compare(a, b)
	switch rel {
	case Less:
		return c == -1
	case LessEqual:
		return c == 0 || c == -1
	case Equal:
		return c == 0
	case GreaterEqual:
		return c == 0 || c == 1
	case Greater:
		return c == 1
	}
	panic(fmt.Sprintf("invalid rel %+v", rel))
}

type Power struct {
	Relationship relationship
	Value        float32
}

func (p Power) Matches(c *card.Card) bool {
	if c.Power != nil {
		var power float32
		if f, err := strconv.ParseFloat(*c.Power, 32); err == nil {
			power = float32(f)
		}
		if fieldCompare(power, p.Relationship, p.Value) {
			return true
		}
	}
	for _, face := range c.CardFaces {
		if face.Power == nil {
			continue
		}
		var power float32
		if f, err := strconv.ParseFloat(*face.Power, 32); err != nil {
			power = float32(f)
		}
		if fieldCompare(power, p.Relationship, p.Value) {
			return true
		}
	}
	return false
}

type Toughness struct {
	Relationship relationship
	Value        float32
}

func (t Toughness) Matches(c *card.Card) bool {
	if c.Toughness != nil {
		var toughness float32
		if f, err := strconv.ParseFloat(*c.Toughness, 32); err == nil {
			toughness = float32(f)
		}
		if fieldCompare(toughness, t.Relationship, t.Value) {
			return true
		}
	}
	for _, face := range c.CardFaces {
		if face.Toughness == nil {
			continue
		}
		var power float32
		if f, err := strconv.ParseFloat(*face.Toughness, 32); err != nil {
			power = float32(f)
		}
		if fieldCompare(power, t.Relationship, t.Value) {
			return true
		}
	}
	return false
}
