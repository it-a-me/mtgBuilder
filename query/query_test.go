package query_test

import (
	"testing"

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
