package query

import (
	"slices"
	"testing"
)

func TestParser(t *testing.T) {
	cases := map[string][]Inequality{
		"f:c":  {{"f", Equal, "c"}},
		"!cow": {{"!name", Equal, "cow"}},
		"!cow f:c": {
			{"!name", Equal, "cow"},
			{"f", Equal, "c"},
		},
		`cmc<=12 cmc>=12 cmc>11 cmc<13 name:/sword .f/ name=Excalibur !"Excalibur, Sword of Eden" t:artifact`: {
			{"cmc", LessEqual, "12"},
			{"cmc", GreaterEqual, "12"},
			{"cmc", Greater, "11"},
			{"cmc", Less, "13"},
			{"name", Equal, `/sword .f/`},
			{"name", Equal, "Excalibur"},
			{"!name", Equal, `"Excalibur, Sword of Eden"`},
			{"t", Equal, "artifact"},
		},
	}
	for c, expected := range cases {
		tokens, err := scan(c)
		if err != nil {
			t.Fatalf("failed to scan '%s', %s", c, err)
		}
		parsed, err := Parse([]rune(c), tokens)
		if err != nil {
			t.Fatalf("failed to parse '%s', %s", c, err)
		}
		if !slices.Equal(parsed, expected) {
			t.Errorf("expected %+v, got %+v when parsing '%s'", expected, parsed, c)
			for i := range min(len(parsed), len(expected)) {
				if parsed[i] != expected[i] {
					t.Errorf("Difference at index %d: %+v != %+v", i, expected[i], parsed[i])
				}
			}
		}
	}
}
