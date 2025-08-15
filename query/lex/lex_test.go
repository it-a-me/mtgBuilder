package lex

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLexer(t *testing.T) {
	testCases := map[string][]Token{
		`o:beaver`: {
			{Literal, "o"},
			{Colon, ":"},
			{Literal, "beaver"},
		},
		`(o:beaver)`: {
			{LParen, "("},
			{Literal, "o"},
			{Colon, ":"},
			{Literal, "beaver"},
			{RParen, ")"},
		},
		`cmc<1`: {
			{Literal, "cmc"},
			{Less, "<"},
			{Literal, "1"},
		},
		`cmc<=1`: {
			{Literal, "cmc"},
			{LessEqual, "<="},
			{Literal, "1"},
		},
		`cmc=1`: {
			{Literal, "cmc"},
			{Equal, "="},
			{Literal, "1"},
		},
		`cmc!=1`: {
			{Literal, "cmc"},
			{BangEqual, "!="},
			{Literal, "1"},
		},
		`cmc>1`: {
			{Literal, "cmc"},
			{Greater, ">"},
			{Literal, "1"},
		},
		`cmc>=1`: {
			{Literal, "cmc"},
			{GreaterEqual, ">="},
			{Literal, "1"},
		},
		`cmc<=1 or o:beaver`: {
			{Literal, "cmc"},
			{LessEqual, "<="},
			{Literal, "1"},
			{Literal, "or"},
			{Literal, "o"},
			{Colon, ":"},
			{Literal, "beaver"},
		},
		`/ creature that saddled it this turn /`: {
			{String, "/ creature that saddled it this turn /"},
		},
		`   / creature that saddled it this turn / `: {
			{String, "/ creature that saddled it this turn /"},
		},
		`!"giant beaver"`: {
			{Bang, "!"},
			{String, `"giant beaver"`},
		},
		`-name:"giant beaver"`: {
			{Minus, "-"},
			{Literal, "name"},
			{Colon, ":"},
			{String, `"giant beaver"`},
		},
		`-(name:"giant beaver" or cmc>3) type:beaver`: {
			{Minus, "-"},
			{LParen, "("},
			{Literal, "name"},
			{Colon, ":"},
			{String, `"giant beaver"`},
			{Literal, "or"},
			{Literal, "cmc"},
			{Greater, ">"},
			{Literal, "3"},
			{RParen, ")"},
			{Literal, "type"},
			{Colon, ":"},
			{Literal, "beaver"},
		},
	}
	var i int
	for query, expected := range testCases {
		i++
		name := fmt.Sprintf("textLexerCase%d", i)
		t.Run(name, func(t *testing.T) { textLexer(t, query, expected) })
	}
}

func textLexer(t *testing.T, query string, expected []Token) {
	l := NewLexer(query)
	tokens, err := l.LexAll()
	if err != nil {
		t.Fatalf("failed to lex `%s`: %s", query, err)
	}
	if diff := cmp.Diff(tokens, expected); diff != "" {
		t.Fatalf("expected\n%+v\n\ngot\n%+v\n%s", expected, tokens, diff)
	}
}
