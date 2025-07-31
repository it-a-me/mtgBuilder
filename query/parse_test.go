package query

import (
	"reflect"
	"regexp"
	"slices"
	"testing"
)

func TestGrouping(t *testing.T) {
	cases := map[string][]Inequality{
		"f:c":  {{"f", Colon, "c"}},
		"!cow": {{"!name", Equal, "cow"}},
		"!cow f:c": {
			{"!name", Equal, "cow"},
			{"f", Colon, "c"},
		},
		`cmc<=12 cmc>=12 cmc>11 cmc<13 name:/sword .f/ name=Excalibur !"Excalibur, Sword of Eden" t:artifact`: {
			{"cmc", LessEqual, "12"},
			{"cmc", GreaterEqual, "12"},
			{"cmc", Greater, "11"},
			{"cmc", Less, "13"},
			{"name", Colon, `/sword .f/`},
			{"name", Equal, "Excalibur"},
			{"!name", Equal, `"Excalibur, Sword of Eden"`},
			{"t", Colon, "artifact"},
		},
	}
	for c, expected := range cases {
		tokens, err := scan(c)
		if err != nil {
			t.Fatalf("failed to scan '%s', %s", c, err)
		}
		parsed, err := groupTokens([]rune(c), tokens)
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

func TestParser(t *testing.T) {
	cases := map[string]Query{
		// TODO"f:c":  {{"f", Colon, "c"}},
		"!cow": Intersection{[]Query{NameExact{"cow"}}},
		// "!cow": {
		// 	ExactName{"cow"},
		// 	{"f", Colon, "c"},
		// },
		`cmc<=12 manavalue>=12 cmc>11 cmc<13 name:/sword .f/ name=Excalibur !"Excalibur, Sword of Eden" t:artifact`: Intersection{[]Query{
			Manavalue{LessEqual, 12},
			Manavalue{GreaterEqual, 12},
			Manavalue{Greater, 11},
			Manavalue{Less, 13},
			NameRegex{regexp.MustCompile(`(?im)sword .f`)},
			Name{"Excalibur"},
			NameExact{"Excalibur, Sword of Eden"},
			Type{"artifact"},
		}},
	}
	for c, expected := range cases {
		parsed, err := Parse(c, false)
		if err != nil {
			t.Fatalf("failed to parse '%s', %s", c, err)
		}
		if !reflect.DeepEqual(parsed, expected) {
			t.Errorf("expected %#v, got %#v when parsing '%s'", expected, parsed, c)
			parsed := parsed.(Intersection)
			expected := expected.(Intersection)
			for i := range min(len(parsed.Queries), len(expected.Queries)) {
				if parsed.Queries[i] != expected.Queries[i] {
					t.Errorf("Difference at index %d: %+v != %+v", i, expected.Queries[i], parsed.Queries[i])
				}
			}
		}
	}
}
