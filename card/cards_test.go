package card_test

import (
	_ "embed"
	// "encoding/json"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/goccy/go-json"

	"mtgBuilder/card"
	"mtgBuilder/fetch"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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

func TestUnmarshalNissa(t *testing.T) {
	want := card.Card{
		CoreFields: card.CoreFields{
			ArenaID:           nil,
			ID:                uuid.UUID{0xa4, 0x71, 0xb3, 0x06, 0x49, 0x41, 0x4e, 0x46, 0xa0, 0xcb, 0xd9, 0x28, 0x95, 0xc1, 0x6f, 0x8a},
			Lang:              "en",
			MtgoID:            &[]int{137223}[0],
			MtgoFoilID:        0,
			MultiverseIDs:     []int{692174},
			TcgplayerID:       &[]int{615195}[0],
			TcgplayerEtchedID: nil,
			CardmarketID:      &[]int{807933}[0],
			Object:            "card",
			Layout:            "normal",
			OracleID:          &[]uuid.UUID{uuid.MustParse("00037840-6089-42ec-8c5c-281f9f474504")}[0],
			PrintsSearchURI:   "https://api.scryfall.com/cards/search?order=released&q=oracleid%3A00037840-6089-42ec-8c5c-281f9f474504&unique=prints",
			RulingsURI:        "https://api.scryfall.com/cards/a471b306-4941-4e46-a0cb-d92895c16f8a/rulings",
			ScryfallURI:       "https://scryfall.com/card/drc/13/nissa-worldsoul-speaker?utm_source=api",
			URI:               "https://api.scryfall.com/cards/a471b306-4941-4e46-a0cb-d92895c16f8a",
		},
		PrintFields: card.PrintFields{
			Artist:           &[]string{"Magali Villeneuve"}[0],
			ArtistIds:        []string{"9e6a55ae-be4d-4c23-a2a5-135737ffd879"},
			AttractionLights: nil,
			Booster:          false,
			BorderColor:      "black",
			CardBackID:       &[]uuid.UUID{uuid.MustParse("0aeebaf5-8c7d-4636-9e82-8c27447861f7")}[0],
			CollectorNumber:  "13",
			ContentWarning:   nil,
			Digital:          false,
			ScryfallSetURI:   "https://scryfall.com/sets/drc?utm_source=api",
			Finishes:         []string{"nonfoil"},
			FlavorName:       nil,
			FlavorText:       &[]string{`"Zendikar still seems so far off, but Chandra is my home."`}[0],
			FrameEffects:     []string{"legendary"},
			Frame:            "2015",
			FullArt:          false,
			Games: []string{
				"paper",
				"mtgo",
			},
			HighresImage:   true,
			IllustrationID: &[]uuid.UUID{uuid.MustParse("7106ab4f-bd3e-4d2a-ba9c-7f223b5a0b7f")}[0],
			ImageStatus:    "highres_scan",
			ImageUris: map[string]string{
				"art_crop":    "https://cards.scryfall.io/art_crop/front/a/4/a471b306-4941-4e46-a0cb-d92895c16f8a.jpg?1738355341",
				"border_crop": "https://cards.scryfall.io/border_crop/front/a/4/a471b306-4941-4e46-a0cb-d92895c16f8a.jpg?1738355341",
				"large":       "https://cards.scryfall.io/large/front/a/4/a471b306-4941-4e46-a0cb-d92895c16f8a.jpg?1738355341",
				"normal":      "https://cards.scryfall.io/normal/front/a/4/a471b306-4941-4e46-a0cb-d92895c16f8a.jpg?1738355341",
				"png":         "https://cards.scryfall.io/png/front/a/4/a471b306-4941-4e46-a0cb-d92895c16f8a.png?1738355341",
				"small":       "https://cards.scryfall.io/small/front/a/4/a471b306-4941-4e46-a0cb-d92895c16f8a.jpg?1738355341",
			},
			Oversized: false,
			Prices: card.Prices{
				Usd:       &[]string{"0.26"}[0],
				UsdFoil:   nil,
				UsdEtched: nil,
				Eur:       &[]string{"0.42"}[0],
				EurFoil:   nil,
				EurEtched: nil,
				Tix:       &[]string{"1.17"}[0],
			},
			PrintedName:     nil,
			PrintedText:     nil,
			PrintedTypeLine: nil,
			Promo:           false,
			PromoTypes:      nil,
			PurchaseURIs: map[string]string{
				"tcgplayer":   "https://partner.tcgplayer.com/c/4931599/1830156/21018?subId1=api&u=https%3A%2F%2Fwww.tcgplayer.com%2Fproduct%2F615195%3Fpage%3D1",
				"cardmarket":  "https://www.cardmarket.com/en/Magic/Products?idProduct=807933&referrer=scryfall&utm_campaign=card_prices&utm_medium=text&utm_source=scryfall",
				"cardhoarder": "https://www.cardhoarder.com/cards/137223?affiliate_id=scryfall&ref=card-profile&utm_campaign=affiliate&utm_medium=card&utm_source=scryfall",
			},
			Rarity: "rare",
			RelatedURIs: map[string]string{
				"edhrec":                      "https://edhrec.com/route/?cc=Nissa%2C+Worldsoul+Speaker",
				"gatherer":                    "https://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=692174&printed=false",
				"tcgplayer_infinite_articles": "https://partner.tcgplayer.com/c/4931599/1830156/21018?subId1=api&trafcat=tcgplayer.com%2Fsearch%2Farticles&u=https%3A%2F%2Fwww.tcgplayer.com%2Fsearch%2Farticles%3FproductLineName%3Dmagic%26q%3DNissa%252C%2BWorldsoul%2BSpeaker",
				"tcgplayer_infinite_decks":    "https://partner.tcgplayer.com/c/4931599/1830156/21018?subId1=api&trafcat=tcgplayer.com%2Fsearch%2Fdecks&u=https%3A%2F%2Fwww.tcgplayer.com%2Fsearch%2Fdecks%3FproductLineName%3Dmagic%26q%3DNissa%252C%2BWorldsoul%2BSpeaker",
			},
			ReleasedAt: "2025-02-14",
			Preview: &card.Preview{
				PreviewedAt: &[]string{"2025-01-23"}[0],
				SourceURI:   &[]string{"https://www.youtube.com/watch?v=km1f1W0Tl6k"}[0],
				Source:      &[]string{"The Command Zone"}[0],
			},
			Reprint:        false,
			SetName:        "Aetherdrift Commander",
			SetSearchURI:   "https://api.scryfall.com/cards/search?order=set&q=e%3Adrc&unique=prints",
			SetType:        "commander",
			SetURI:         "https://api.scryfall.com/sets/d33ef7a4-41bb-4f16-bad3-b3ee13c257e6",
			Set:            "drc",
			SetID:          uuid.UUID{0xd3, 0x3e, 0xf7, 0xa4, 0x41, 0xbb, 0x4f, 0x16, 0xba, 0xd3, 0xb3, 0xee, 0x13, 0xc2, 0x57, 0xe6},
			StorySpotlight: false,
			Textless:       false,
			Variation:      false,
			VariationOf:    nil,
			SecurityStamp:  &[]string{"oval"}[0],
			Watermark:      &[]string{"desparked"}[0],
			Foil:           &[]bool{false}[0],
			NonFoil:        &[]bool{true}[0],
		},
		AllParts: []card.RelatedCard{
			{
				ID:        uuid.MustParse("9aeb44e0-2257-444a-b805-7e939ca5f6fe"),
				Object:    "related_card",
				Component: "combo_piece",
				Name:      "Nissa, Worldsoul Speaker",
				TypeLine:  "Legendary Creature — Elf Druid",
				URI:       "https://api.scryfall.com/cards/9aeb44e0-2257-444a-b805-7e939ca5f6fe",
			},
			{
				ID:        uuid.MustParse("6a2c1fa5-deed-48ba-afe4-6c8ea8d9135e"),
				Object:    "related_card",
				Component: "combo_piece",
				Name:      "Energy Reserve",
				TypeLine:  "Card",
				URI:       "https://api.scryfall.com/cards/6a2c1fa5-deed-48ba-afe4-6c8ea8d9135e",
			},
		},
		CardFaces:      nil,
		Cmc:            &[]float32{4}[0],
		ColorIdentity:  &card.Colors{"G"},
		ColorIndicator: nil,
		Colors:         &card.Colors{"G"},
		Defense:        nil,
		EdhrecRank:     &[]int{8490}[0],
		GameChanger:    &[]bool{false}[0],
		HandModifier:   nil,
		Keywords:       []string{"Landfall"},
		Legalities: card.Legalities{
			"alchemy":         "not_legal",
			"brawl":           "not_legal",
			"commander":       "legal",
			"duel":            "legal",
			"future":          "not_legal",
			"gladiator":       "not_legal",
			"historic":        "not_legal",
			"legacy":          "legal",
			"modern":          "not_legal",
			"oathbreaker":     "legal",
			"oldschool":       "not_legal",
			"pauper":          "not_legal",
			"paupercommander": "not_legal",
			"penny":           "not_legal",
			"pioneer":         "not_legal",
			"predh":           "not_legal",
			"premodern":       "not_legal",
			"standard":        "not_legal",
			"standardbrawl":   "not_legal",
			"timeless":        "not_legal",
			"vintage":         "legal",
		},
		LifeModifier: nil,
		Loyalty:      nil,
		ManaCost:     &[]string{"{3}{G}"}[0],
		Name:         "Nissa, Worldsoul Speaker",
		OracleText:   &[]string{"Landfall — Whenever a land you control enters, you get {E}{E} (two energy counters).\nYou may pay eight {E} rather than pay the mana cost for permanent spells you cast."}[0],
		PennyRank:    nil,
		Power:        &[]string{"3"}[0],
		ProducedMana: nil,
		Reserved:     false,
		Toughness:    &[]string{"3"}[0],
		TypeLine:     "Legendary Creature — Elf Druid",
	}

	content, err := os.ReadFile("testdata/nissa.json")
	if err != nil {
		t.Fatal(err)
	}
	var got card.Card
	err = json.Unmarshal(content, &got)
	if err != nil {
		log.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("want != got; -want +got\n%s", diff)
	}
}

func TestReflectAll(t *testing.T) {
	if !t.Run("Layout", ReflectLayout) {
		return
	}
	t.Run("Layout", ReflectAll)
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
		if diff := cmp.Diff(expected, got); diff != "" {
			t.Errorf("json roundtrip mismatch (-want +got):\n%s", diff)
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

func ReflectAll(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
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
	if diff := cmp.Diff(expected, roundtrip); diff != "" {
		t.Errorf("json roundtrip mismatch (-want +got):\n%s", diff)
	}
}
