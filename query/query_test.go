package query_test

import (
	"bytes"
	"log/slog"
	"testing"

	"mtgBuilder/card"
	"mtgBuilder/query"
)

func TestStripParens(t *testing.T) {
	before := `({T}: Add {R}, {W}, or {B}.)
This land enters tapped.
Cycling {3} ({3}, Discard this card: Draw a card.)`
	expected := `
This land enters tapped.
Cycling {3} `
	got := query.StripParens.ReplaceAllString(before, "")
	if expected != got {
		t.Fatalf("expected %+v, got %+v", expected, got)
	}
}

func TestPower(t *testing.T) {
	cases := []struct {
		card     card.Card
		expected bool
	}{{card.Card{}, false}, {card.Card{Power: &[]string{"1"}[0]}, true}}
	q, err := query.Parse("power=1", false)
	if err != nil {
		t.Fatalf("failed to parse query: %s", err)
	}
	for _, testcase := range cases {
		logBuf := bytes.Buffer{}
		slog.SetDefault(slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug})))

		got := q.Matches(&testcase.card)
		if got != testcase.expected {
			t.Errorf("got %t, expected %t when matching %#v on %#v\n\n%s", got, testcase.expected, q, testcase.card, logBuf.String())
		}
	}
}
