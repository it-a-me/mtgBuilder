package query

import (
	"errors"
	"fmt"
)

type Relationship int8

const (
	Invalid Relationship = iota
	Less                 = iota
	LessEqual
	Equal
	GreaterEqual
	Greater
)

type Inequality struct {
	Left         string
	Relationship Relationship
	Right        string
}

var ErrInvalidRelationship = errors.New("invalid relationship")

func ParseRelationship(inner string) (Relationship, error) {
	switch inner {
	case "<":
		return Less, nil
	case "<=":
		return LessEqual, nil
	case "=", ":":
		return Equal, nil
	case ">":
		return Greater, nil
	case ">=":
		return GreaterEqual, nil
	}
	return Relationship(-1), fmt.Errorf("%w: '%s' is not an relationship", ErrInvalidRelationship, inner)
}

var ErrInvalidBang = errors.New("invalid '!'")

func ParseInequality(ineq Inequality) (Query, error) {
	panic("")
}

var (
	ErrUnfinishedInequality = errors.New("unfinished inequality")
	ErrUnexpectedTokenType  = errors.New("unexpected token")
)

func Parse(runes []rune, tokens []Token) ([]Inequality, error) {
	var inqualities []Inequality
	var exact bool
	var ineq Inequality
	for _, t := range tokens {
		switch {
		case ineq.Left != "" && ineq.Relationship != Invalid && t.Type == RHS:
			ineq.Right = string(t.Get(runes))
			inqualities = append(inqualities, ineq)
			ineq = Inequality{}
		case ineq.Left != "" && ineq.Relationship != Invalid:
			return nil, fmt.Errorf("%w: Expected RHS, found %s", ErrUnexpectedTokenType, "???")
		case ineq.Left != "" && t.Type == Comparison:
			relationship, err := ParseRelationship(string(t.Get(runes)))
			if err != nil {
				return nil, err
			}
			ineq.Relationship = relationship
		case ineq.Left != "":
			inqualities = append(inqualities, Inequality{"name", Equal, ineq.Left})
		case t.Type == Bang && !exact:
			exact = true
		case t.Type == Bang:
			return nil, ErrInvalidBang
		case t.Type == LHS && exact:
			exact = false
			inqualities = append(inqualities, Inequality{"!name", Equal, string(t.Get(runes))})
		case t.Type == LHS:
			ineq.Left = string(t.Get(runes))
		default:
			return nil, fmt.Errorf("%w: token %s starting at character %d", ErrUnexpectedTokenType, t.Type.String(), t.Start)
		}
	}
	return inqualities, nil
}
