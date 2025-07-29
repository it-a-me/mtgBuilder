package query

import (
	"mtgBuilder/card"
)

type Negation struct {
	Query Query
}

func (n Negation) Matches(c *card.Card) bool {
	return !n.Query.Matches(c)
}

type Union struct {
	Queries []Query
}

func (u Union) Matches(c *card.Card) bool {
	for _, q := range u.Queries {
		if q.Matches(c) {
			return true
		}
	}
	return false
}

type Intersection struct {
	Queries []Query
}

func (i Intersection) Matches(c *card.Card) bool {
	for _, q := range i.Queries {
		if !q.Matches(c) {
			return false
		}
	}
	return true
}

var DefaultFilter = Intersection{[]Query{
	Negation{Type{"vanguard"}},
	Negation{Type{"plane"}},
	Negation{Type{"scheme"}},
	Negation{Type{"phenomenon"}},
	Negation{SetType{"memorabilia"}},
	Negation{SetType{"minigame"}},
	Negation{Set{"unk"}},
}}

func Parse(queryline string, withDefault bool) ([]Query, error) {
	runes := []rune(queryline)
	tokens, err := scan(queryline)
	if err != nil {
		return nil, err
	}
	inequalities, err := groupTokens(runes, tokens)
	if err != nil {
		return nil, err
	}

	var queries []Query
	for _, inequality := range inequalities {
		query, err := parseInequality(inequality)
		if err != nil {
			return nil, err
		}
		queries = append(queries, query)
	}
	if withDefault {
		queries = append(queries, DefaultFilter)
	}
	return queries, nil
}
