package query

import (
	"errors"
	"fmt"
	"log/slog"
	"unicode"
)

type Token struct {
	Start, End int
	Type       TokenType
}

func (t Token) Get(runes []rune) []rune {
	return runes[t.Start:t.End]
}

type TokenType int

const (
	Bang TokenType = iota
	Comparison
	LHS
	RHS
)

func (t TokenType) String() string {
	switch t {
	case Bang:
		return "Bang"
	case Comparison:
		return "Comparison"
	case LHS:
		return "Left Hand Side"
	case RHS:
		return "Right Hand Side"
	default:
		panic(fmt.Sprintf("Invalid TokenType: %d", t))
	}
}

func scan(queryLine string) ([]Token, error) {
	var tokens []Token
	runes := []rune(queryLine)
	var i int
	for i < len(runes) {
		r := queryLine[i]
		switch r {
		case '!':
			tokens = append(tokens, Token{i, i + 1, Bang})
		case '"':
			consumed, err := handleQuote(runes[i+1:], '"')
			if err != nil {
				return nil, err
			}
			t := LHS
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == Comparison {
				t = RHS
			}
			tokens = append(tokens, Token{i, i + consumed, t})
			i += consumed
			continue
		case '/':
			consumed, err := handleQuote(runes[i+1:], '/')
			if err != nil {
				return nil, err
			}
			t := LHS
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == Comparison {
				t = RHS
			}
			tokens = append(tokens, Token{i, i + consumed, t})
			i += consumed
			continue
		case ':':
			tokens = append(tokens, Token{i, i + 1, Comparison})
		case '=':
			tokens = append(tokens, Token{i, i + 1, Comparison})
		case '<':
			var next rune = -1
			if i+1 < len(runes) {
				next = runes[i+1]
			}
			if next == '=' {
				tokens = append(tokens, Token{i, i + 2, Comparison})
				i++
			} else {
				tokens = append(tokens, Token{i, i + 1, Comparison})
			}
		case '>':
			var next rune = -1
			if i+1 < len(runes) {
				next = runes[i+1]
			}
			if next == '=' {
				tokens = append(tokens, Token{i, i + 2, Comparison})
				i++
			} else {
				tokens = append(tokens, Token{i, i + 1, Comparison})
			}
		case ' ':
		default:
			consumed, err := handleUnquotedLiteral(runes[i:])
			if err != nil {
				return nil, err
			}
			t := LHS
			if len(tokens) > 0 && tokens[len(tokens)-1].Type == Comparison {
				t = RHS
			}
			tokens = append(tokens, Token{i, i + consumed, t})
			i += consumed
			continue
		}
		i++
	}
	return tokens, nil
}

var (
	ErrUnexpectedRune       = errors.New("unknown character")
	ErrUnexpectedEndOfInput = errors.New("unexpected end of input")
)

func handleUnquotedLiteral(runes []rune) (consumed int, err error) {
	if !unicode.IsLetter(runes[0]) && !unicode.IsDigit(runes[0]) {
		return -1, fmt.Errorf("%w: character '%c'", ErrUnexpectedRune, runes[0])
	}
	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return i, nil
		}
	}
	return len(runes), nil
}

func handleQuote(runes []rune, delimiter rune) (consumed int, err error) {
	slog.Debug("Starting Quote", "delimiter", delimiter)
	escaped := false
	for i, r := range runes {
		slog.Debug("Handling Quote", "encountered", string(r), "escaped", escaped)
		if escaped {
			escaped = false
			continue
		}
		switch r {
		case '\\':
			slog.Debug("Found Backslash")
			escaped = true
		case delimiter:
			slog.Debug("Found Delimiter", "consumed", i+1)
			return i + 2, nil
		}
	}
	return -1, fmt.Errorf("%w: missing delimiter %c", ErrUnexpectedEndOfInput, delimiter)
}
