package query

import (
	"bytes"
	"log/slog"
	"slices"
	"testing"
)

func TestScan(t *testing.T) {
	cases := map[string][]Token{
		"!":                        {{0, 1, Bang}},
		" ! ":                      {{1, 2, Bang}},
		`!"woolf"`:                 {{0, 1, Bang}, {1, 8, LHS}},
		`o`:                        {{0, 1, LHS}},
		`o:horse`:                  {{0, 1, LHS}, {1, 2, Comparison}, {2, 7, RHS}},
		`oracle:/goblin.*warrior/`: {{0, 6, LHS}, {6, 7, Comparison}, {7, 24, RHS}},
		`f:c oracle:/goblin.*warrior/ !"Wort, the Raidmother"`: {
			{0, 1, LHS},
			{1, 2, Comparison},
			{2, 3, RHS},
			{4, 10, LHS},
			{10, 11, Comparison},
			{11, 28, RHS},
			{29, 30, Bang},
			{30, 52, LHS},
		},
	}
	for line, expected := range cases {
		buf := bytes.Buffer{}
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
		got, err := scan(line)
		if err != nil {
			t.Fatalf("failed to scan %+v: %s\n%s", line, err, buf.String())
		}
		var words []string
		runes := []rune(line)
		for _, t := range got {
			words = append(words, string(t.Get(runes)))
		}
		if !slices.Equal(got, expected) {
			t.Fatalf("expected %+v, got %+v for %+v. words: %#v", expected, got, line, words)
		}
	}
}
