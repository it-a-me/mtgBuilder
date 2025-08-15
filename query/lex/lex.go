package lex

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

//go:generate go tool stringer -type TokenType
type TokenType int

const (
	Error TokenType = iota
	Bang
	Minus
	Less
	LessEqual
	Equal
	BangEqual
	GreaterEqual
	Greater
	Colon
	LParen
	RParen
	String
	Literal
)

type Token struct {
	Type    TokenType
	Content string
}

func (t Token) String() string {
	if len(t.Content) < 10 {
		return fmt.Sprintf("%s('%s')", t.Type, t.Content)
	}
	return fmt.Sprintf("%s('%s'...)", t.Type, t.Content[:10])
}

type Lexer struct {
	Start   int
	Pos     int
	Content string
	Tokens  []Token
}

func NewLexer(content string) Lexer {
	return Lexer{Content: content}
}

var ErrInvalidUtf8 = errors.New("invalid utf8")

func (l *Lexer) LexAll() ([]Token, error) {
	for {
		if l.peek() == -1 {
			return l.Tokens, nil
		}
		switch {
		// Two character
		case l.matchString("!="):
			l.Tokens = append(l.Tokens, Token{BangEqual, "!="})
		case l.matchString("<="):
			l.Tokens = append(l.Tokens, Token{LessEqual, "<="})
		case l.matchString(">="):
			l.Tokens = append(l.Tokens, Token{GreaterEqual, ">="})

		// Single character
		case l.matchRune('!'):
			l.Tokens = append(l.Tokens, Token{Bang, "!"})
		case l.matchRune('-'):
			l.Tokens = append(l.Tokens, Token{Minus, "-"})
		case l.matchRune('<'):
			l.Tokens = append(l.Tokens, Token{Less, "<"})
		case l.matchRune('='):
			l.Tokens = append(l.Tokens, Token{Equal, "="})
		case l.matchRune('>'):
			l.Tokens = append(l.Tokens, Token{Greater, ">"})
		case l.matchRune(':'):
			l.Tokens = append(l.Tokens, Token{Colon, ":"})
		case l.matchRune('('):
			l.Tokens = append(l.Tokens, Token{LParen, "("})
		case l.matchRune(')'):
			l.Tokens = append(l.Tokens, Token{RParen, ")"})
		case l.matchRune(' '): // Ignore whitespace
		case l.startsWith(`"`):
			if err := l.takeString('"'); err != nil {
				return nil, err
			}
		case l.startsWith(`/`):
			if err := l.takeString('/'); err != nil {
				return nil, err
			}
		case isLiteralChar(l.peek()):
			if err := l.parseLiteral(); err != nil {
				return nil, err
			}
		default:
			panic(fmt.Sprintf("invalid character '%c' at %d in '%s'", l.peek(), l.Pos, l.Content))
		}
		l.Start = l.Pos
	}
}

func isLiteralChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '.'
}

var ErrInvalidLiteralRune = errors.New("invalid character in literal string")

func (l *Lexer) parseLiteral() error {
	end := -1
	for _, c := range l.Content[l.Pos:] {
		if c == ' ' {
			end = l.Pos
			l.Pos += utf8.RuneLen(c)
			break
		} else if strings.ContainsRune(":<=>)", c) || // vaild trailing characters
			strings.HasPrefix(l.Content[l.Pos:], "!=") {
			end = l.Pos
			break
		} else if !isLiteralChar(c) {
			return fmt.Errorf("%w: %c is not a valid literal character", ErrInvalidLiteralRune, c)
		}

		l.Pos += utf8.RuneLen(c)
	}
	if end == -1 {
		end = len(l.Content)
	}
	t := Token{Literal, l.Content[l.Start:end]}
	l.Tokens = append(l.Tokens, t)
	return nil
}

var ErrMissingDelimiter = errors.New("missing delimiter")

func (l *Lexer) takeString(delimiter rune) error {
	start := l.Pos
	if !l.matchRune(delimiter) {
		panic("strings must start with their delimiter")
	}
	for _, c := range l.Content[l.Pos:] {
		l.Pos += utf8.RuneLen(c)
		if c == delimiter {
			t := Token{String, l.Content[start:l.Pos]}
			l.Tokens = append(l.Tokens, t)
			return nil
		}
	}
	return fmt.Errorf("%w: expected %c", ErrMissingDelimiter, delimiter)
}

// peek returns the next utf8 rune if it exists or -1 if it does not
func (l *Lexer) peek() rune {
	r, _ := utf8.DecodeRuneInString(l.Content[l.Pos:])
	if r == utf8.RuneError {
		return -1
	}
	return r
}

// next attempts to advance by one rune and return the next utf8 rune. Otherwise it returns -1
func (l *Lexer) next() rune {
	r := l.peek()
	if r != -1 {
		l.Pos += utf8.RuneLen(r)
	}
	return r
}

// matchRune advances if r is the next rune and returns if it was successful
func (l *Lexer) matchRune(r rune) bool {
	if l.peek() == r {
		l.next()
		return true
	}
	return false
}

// startsWith returns whether the the remaining content starts with s
func (l *Lexer) startsWith(s string) bool {
	return strings.HasPrefix(l.Content[l.Pos:], s)
}

func (l *Lexer) matchString(s string) bool {
	if l.startsWith(s) {
		l.Pos += utf8.RuneCountInString(s)
		return true
	}
	return false
}
