package card_test

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"mtgBuilder/card"
	"mtgBuilder/fetch"
)

func init() {
	// HACK to prevent large binary blob
	if _, err := os.Stat("testdata/cards.json"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cards, err := fetch.GetOracleCardsJSON()
			if err != nil {
				panic(err)
			}
			err = os.WriteFile("testdata/cards.json", cards, 0o600)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func ReflectLayout(t *testing.T) {
	cards := []string{"double_faced.json", "flip.json", "mdf.json", "nissa.json", "reversible.json", "split.json"}
	for _, c := range cards {
		content, err := os.ReadFile("testdata/" + c)
		if err != nil {
			t.Fatalf("failed to read %s: %s", c, err)
		}
		var parsed card.Card
		err = json.Unmarshal(content, &parsed)
		if err != nil {
			t.Fatalf("failed to unmarshal %s into card: %s", c, err)
		}
		marshaled, err := json.Marshal(parsed)
		if err != nil {
			t.Fatalf("failed to marshal %s: %s", c, err)
		}
		var expected, got map[string]any
		if err := json.Unmarshal(content, &expected); err != nil {
			t.Fatalf("failed to unmarshal %s: %s", c, err)
		}
		if err := json.Unmarshal(marshaled, &got); err != nil {
			t.Fatalf("failed to unmarshal parsed card %s: %s", c, err)
		}
		if !reflect.DeepEqual(got, expected) {
			a, _ := json.MarshalIndent(expected, "", "  ")
			b, _ := json.MarshalIndent(got, "", "  ")
			t.Fatalf("expected != got\nexpected:\n%s\ngot\n%s", string(a), string(b))
		}

	}
}

func TestUnmarshalAll(t *testing.T) {
	p := "testdata/cards.json"
	content, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read file %s:%s", p, err)
	}
	var cards []card.Card
	err = json.Unmarshal(content, &cards)
	if err != nil {
		t.Fatalf("failed to unmarshal all cards: %s", err)
	}
}

func TestReflectAll(t *testing.T) {
	if !t.Run("Layout", ReflectLayout) {
		return
	}
	t.Run("Layout", ReflectAll)
}

func ReflectAll(t *testing.T) {
	p := "testdata/cards.json"
	content, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read file %s:%s", p, err)
	}
	var roundtrip []any
	var got []byte
	{
		var cards []card.Card
		err = json.Unmarshal(content, &cards)
		if err != nil {
			t.Fatalf("failed to unmarshal all cards to []card.Card: %s", err)
		}
		got, err = json.MarshalIndent(cards, "", "  ")
		if err != nil {
			t.Fatalf("failed to unmarshal all cards from []card.Card to map[string]string: %s", err)
		}
		err = json.Unmarshal(got, &roundtrip)
		if err != nil {
			t.Fatalf("failed to unmarshal all cards from []card.Card to map[string]string: %s", err)
		}
	}
	var expected []any
	{
		err = json.Unmarshal(content, &expected)
		if err != nil {
			t.Fatalf("failed to unmarshal all cards to map[string]string: %s", err)
		}
	}
	if !reflect.DeepEqual(roundtrip, expected) {
		t.Fatalf("expected != got\ngot:\n%s", string(got))
	}
}
