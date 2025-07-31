package query

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"mtgBuilder/card"
)

//go:generate stringer -type=relationship
type relationship int8

const (
	Invalid relationship = iota
	Less
	LessEqual
	Equal
	GreaterEqual
	Greater
	Colon
)

func (r relationship) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

type Inequality struct {
	Left         string
	Relationship relationship
	Right        string
}

var ErrInvalidRelationship = errors.New("invalid relationship")

func parseRelationship(inner string) (relationship, error) {
	switch inner {
	case "<":
		return Less, nil
	case "<=":
		return LessEqual, nil
	case "=":
		return Equal, nil
	case ">":
		return Greater, nil
	case ">=":
		return GreaterEqual, nil
	case ":":
		return Colon, nil
	}
	return relationship(-1), fmt.Errorf("%w: '%s' is not an relationship", ErrInvalidRelationship, inner)
}

var ErrInvalidBang = errors.New("invalid '!'")

var fieldAliases = map[string]string{
	"c":     "color",
	"id":    "identity",
	"ci":    "identity",
	"t":     "type",
	"o":     "oracle",
	"fo":    "fulloracle",
	"kw":    "keyword",
	"m":     "mana",
	"mv":    "manavalue",
	"cmc":   "manavalue",
	"st":    "set_type",
	"f":     "format",
	"legal": "format",
	"p":     "power",
	"pow":   "power",
	"tou":   "toughness",
}

var ErrInvalidColor = errors.New("invalid color")

func stripQuotes(s string) (res string, success bool) {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1], true
	}
	return "", false
}

func stripSlash(s string) (res string, success bool) {
	if len(s) >= 2 && s[0] == '/' && s[len(s)-1] == '/' {
		return s[1 : len(s)-1], true
	}
	return "", false
}

func parseColorString(s string) (card.Colors, error) {
	switch s {
	case "white":
		return card.NewColor("white"), nil
	case "blue":
		return card.NewColor("blue"), nil
	case "black":
		return card.NewColor("black"), nil
	case "red":
		return card.NewColor("red"), nil
	case "green":
		return card.NewColor("green"), nil
	case "c", "colorless":
		return card.Colors{}, nil
	}
	var colors card.Colors
	for _, r := range s {
		switch r {
		case 'w':
			colors.Add(card.NewColor("white"))
		case 'u':
			colors.Add(card.NewColor("blue"))
		case 'b':
			colors.Add(card.NewColor("black"))
		case 'r':
			colors.Add(card.NewColor("red"))
		case 'g':
			colors.Add(card.NewColor("green"))
		default:
			return card.Colors{}, fmt.Errorf("%w: invalid color %c", ErrInvalidColor, r)
		}
	}
	return colors, nil
}

func parseColor(ineq Inequality) (Query, error) {
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	val = strings.ToLower(val)
	if val == "m" || val == "multicolor" {
		return Color{Mulicolor: true}, nil
	}
	color, err := parseColorString(val)
	if err != nil {
		return nil, err
	}
	if ineq.Relationship == Colon {
		ineq.Relationship = Equal
	}
	return Color{Colors: color}, nil
}

func parseColorIdentity(ineq Inequality) (Query, error) {
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	val = strings.ToLower(val)
	if val == "m" || val == "multicolor" {
		return Color{Mulicolor: true}, nil
	}
	color, err := parseColorString(val)
	if err != nil {
		return nil, err
	}
	if ineq.Relationship == Colon {
		ineq.Relationship = Equal
	}
	return Color{Colors: color}, nil
}

func parseType(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	return Type{val}, nil
}

func parseOracle(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	} else if stripped, ok := stripSlash(ineq.Right); ok {
		re, err := regexp.Compile("(?im)" + stripped)
		if err != nil {
			return nil, err
		}
		return OracleTextRegex{re}, nil
	}
	return OracleText{val}, nil
}

func parseFullOracle(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	} else if stripped, ok := stripSlash(ineq.Right); ok {
		re, err := regexp.Compile("(?im)" + stripped)
		if err != nil {
			return nil, err
		}
		return FullOracleTextRegex{re}, nil
	}
	return FullOracleText{val}, nil
}

func parseKeyword(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	return Type{val}, nil
}

func parseSet(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	return Set{val}, nil
}

func parseSetType(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	return SetType{val}, nil
}

func parseMana(ineq Inequality) (Query, error) {
	panic("todo")
}

func parseManavalue(ineq Inequality) (Query, error) {
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return nil, err
	}
	return Manavalue{ineq.Relationship, float32(f)}, nil
}

func parseName(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	} else if stripped, ok := stripSlash(ineq.Right); ok {
		re, err := regexp.Compile("(?im)" + stripped)
		if err != nil {
			return nil, err
		}
		return NameRegex{re}, nil
	}
	return Name{val}, nil
}

func parseNameExact(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	return NameExact{val}, nil
}

func parseFormat(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := strings.ToLower(ineq.Right)
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	if expanded, exists := formatAliases[val]; exists {
		val = expanded
	}
	switch ineq.Left {
	case "format":
		return Format{val, "legal"}, nil
	case "banned":
		return Format{val, "banned"}, nil
	case "restricted":
		return Format{val, "restricted"}, nil
	}
	panic(fmt.Sprintf("invalid format: %#v", ineq))
}

func parsePower(ineq Inequality) (Query, error) {
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	var f float32
	if val != "*" && val != "X" {
		f64, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}
		f = float32(f64)
	}
	return Power{ineq.Relationship, f}, nil
}

func parseToughness(ineq Inequality) (Query, error) {
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	var f float32
	if val != "*" && val != "X" {
		f64, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}
		f = float32(f64)
	}
	return Toughness{ineq.Relationship, f}, nil
}

func parseOracleID(ineq Inequality) (Query, error) {
	if ineq.Relationship != Equal && ineq.Relationship != Colon {
		return nil, fmt.Errorf("%w: unable to compare %s with %+v", ErrInvalidRelationship, ineq.Left, ineq.Relationship)
	}
	val := ineq.Right
	if stripped, ok := stripQuotes(ineq.Right); ok {
		val = stripped
	}
	return OracleID{val}, nil
}

var ErrUnknownField = errors.New("unknown field")

func parseInequality(ineq Inequality) (Query, error) {
	if ineq.Relationship == Invalid {
		panic(ineq)
	}
	field := strings.ToLower(ineq.Left)
	if expanded, exists := fieldAliases[field]; exists {
		field = expanded
	}
	switch field {
	case "color":
		return parseColor(ineq)
	case "identity":
		return parseColorIdentity(ineq)
	case "type":
		return parseType(ineq)
	case "oracle":
		return parseOracle(ineq)
	case "fulloracle":
		return parseFullOracle(ineq)
	case "keyword":
		return parseKeyword(ineq)
	case "mana":
		return parseMana(ineq)
	case "manavalue":
		return parseManavalue(ineq)
	case "name":
		return parseName(ineq)
	case "!name":
		return parseNameExact(ineq)
	case "set":
		return parseSet(ineq)
	case "set_type":
		return parseSetType(ineq)
	case "format", "banned", "restricted":
		ineq.Left = field
		return parseFormat(ineq)
	case "power":
		return parsePower(ineq)
	case "toughness":
		return parseToughness(ineq)
	case "oracle_id":
		return parseOracleID(ineq)
	}
	return nil, fmt.Errorf("%w: '%s'", ErrUnknownField, field)
}

var (
	ErrUnfinishedInequality = errors.New("unfinished inequality")
	ErrUnexpectedTokenType  = errors.New("unexpected token")
)

func groupTokens(runes []rune, tokens []Token) ([]Inequality, error) {
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
			relationship, err := parseRelationship(string(t.Get(runes)))
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
