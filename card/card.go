package card

import (
	"github.com/google/uuid"
)

type Legalities map[string]string

type CoreFields struct {
	// This card’s Arena ID, if any. A large percentage of cards are not available on Arena and do not have this ID.
	//
	// This value may be nil
	ArenaID *int `json:"arena_id,omitempty"`

	// A unique ID for this card in Scryfall’s database.
	ID uuid.UUID `json:"id"`

	// A language code for this printing.
	Lang string `json:"lang"`

	// This card’s Magic Online ID (also known as the Catalog ID), if any.
	// A large percentage of cards are not available on Magic Online and do not have this ID.
	//
	// This value may be nil
	MtgoID *int `json:"mtgo_id,omitempty"`

	// This card’s foil Magic Online ID (also known as the Catalog ID), if any. A large percentage of cards are not available on Magic Online and do not have this ID.
	//
	// This value may be nil
	MtgoFoilID int `json:"mtgo_foil_id,omitempty"`

	// This card’s multiverse IDs on Gatherer, if any, as an array of integers. Note that Scryfall includes many promo cards, tokens, and other esoteric objects that do not have these identifiers.
	//
	// This value may be nil
	MultiverseIDs []int `json:"multiverse_ids"`

	// This card’s ID on TCGplayer’s API, also known as the productId.
	//
	// This value may be nil
	TcgplayerID *int `json:"tcgplayer_id,omitempty"`

	// This card’s ID on TCGplayer’s API, for its etched version if that version is a separate product.
	//
	// This value may be nil
	TcgplayerEtchedID *int `json:"tcgplayer_etched_id,omitempty"`

	// This card’s ID on Cardmarket’s API, also known as the idProduct."
	//
	// This value may be nil
	CardmarketID *int `json:"cardmarket_id,omitempty"`

	// A content type for this object, always card.
	Object string `json:"object"`

	// A code for this card’s layout.
	Layout string `json:"layout"`

	// A unique ID for this card’s oracle identity.
	// This value is consistent across reprinted card editions, and unique among different cards with the same name (tokens, Unstable variants, etc).
	// Always present except for the reversibleCard layout where it will be absent; oracleId will be found on each face instead.
	//
	// This value may be nil
	OracleID *uuid.UUID `json:"oracle_id,omitempty"`

	// A link to where you can begin paginating all re/prints for this card on Scryfall’s API.
	PrintsSearchURI string `json:"prints_search_uri"`

	// A link to this card’s rulings list on Scryfall’s API.
	RulingsURI string `json:"rulings_uri"`

	// A link to this card’s permapage on Scryfall’s website.
	ScryfallURI string `json:"scryfall_uri"`

	// A link to this card object on Scryfall’s API.
	URI string `json:"uri"`
}

type Card struct {
	CoreFields
	PrintFields

	// If this card is closely related to other cards, this property will be an array with Related Card Objects.
	//
	// This value may be nil
	AllParts []RelatedCard `json:"all_parts,omitempty"`

	// An array of Card Face objects, if this card is multifaced.
	//
	// This value may be nil
	CardFaces []CardFace `json:"card_faces,omitempty"`

	// The card’s mana value. Note that some funny cards have fractional mana costs.
	//
	// This value may be nil but is expected in the json output (Bad docs)
	Cmc *float32 `json:"cmc,omitempty"`

	// This card’s color identity.
	ColorIdentity *Colors `json:"color_identity"`

	// The colors in this card’s color indicator, if any.
	// A null value for this field indicates the card does not have one.
	//
	// This value may be nil
	ColorIndicator []string `json:"color_indicator,omitempty"`

	// This card’s colors, if the overall card has colors defined by the rules. Otherwise the colors will be on the card_faces objects, see below.
	//
	// This value may be nil
	Colors *Colors `json:"colors,omitempty"`

	// This face’s defense, if any.
	//
	// This value may be nil
	Defense *string `json:"defense,omitempty"`

	// This card’s overall rank/popularity on EDHREC. Not all cards are ranked.
	//
	// This value may be nil
	EdhrecRank *int `json:"edhrec_rank,omitempty"`

	// True if this card is on the Commander Game Changer list.
	//
	// This value may be nil
	GameChanger *bool `json:"game_changer,omitempty"`

	// This card’s hand modifier, if it is Vanguard card. This value will contain a delta, such as -1.
	//
	// This value may be nil
	HandModifier *string `json:"hand_modifier,omitempty"`

	// An array of keywords that this card uses, such as 'Flying' and 'Cumulative upkeep'.
	Keywords []string `json:"keywords"`

	// An object describing the legality of this card across play formats. Possible legalities are legal, not_legal, restricted, and banned.
	Legalities Legalities `json:"legalities"`

	// This card’s life modifier, if it is Vanguard card. This value will contain a delta, such as +2.
	//
	// This value may be nil
	LifeModifier *string `json:"life_modifier,omitempty"`

	// This loyalty if any. Note that some cards have loyalties that are not numeric, such as X.
	//
	// This value may be nil
	Loyalty *string `json:"loyalty,omitempty"`

	// The mana cost for this card. This value will be any empty string "" if the cost is absent. Remember that per the game rules, a missing mana cost and a mana cost of {0} are different values. Multi-faced cards will report this value in card faces.
	//
	// This value may be nil
	ManaCost *string `json:"mana_cost,omitempty"`

	// The name of this card. If this card has multiple faces, this field will contain both names separated by ␣//␣.
	Name string `json:"name"`
	// The Oracle text for this card, if any.
	//
	// This value may be nil
	OracleText *string `json:"oracle_text,omitempty"`

	// This card’s rank/popularity on Penny Dreadful. Not all cards are ranked.
	//
	// This value may be nil
	PennyRank *int `json:"penny_rank,omitempty"`

	// This card’s power, if any. Note that some cards have powers that are not numeric, such as *.
	//
	// This value may be nil
	Power *string `json:"power,omitempty"`

	// Colors of mana that this card could produce.
	//
	// This value may be nil
	ProducedMana []string `json:"produced_mana,omitempty"`

	// True if this card is on the Reserved List.
	Reserved bool `json:"reserved"`

	// This card’s toughness, if any. Note that some cards have toughnesses that are not numeric, such as *.
	//
	// This value may be nil
	Toughness *string `json:"toughness,omitempty"`

	// The type line of this card.
	TypeLine string `json:"type_line,omitzero"`
}

type PrintFields struct {
	// The name of the illustrator of this card. Newly spoiled cards may not have this field yet.
	//
	// This value may be nil
	Artist *string `json:"artist,omitempty"`

	// The IDs of the artists that illustrated this card. Newly spoiled cards may not have this field yet.
	//
	// This value may be nil
	ArtistIds []string `json:"artist_ids,omitempty"`

	// The lit Unfinity attractions lights on this card, if any.
	//
	// This value may be nil
	AttractionLights []uint `json:"attraction_lights,omitempty"`

	// Whether this card is found in boosters.
	Booster bool `json:"booster"`

	// This card’s border color: black, white, borderless, yellow, silver, or gold.
	BorderColor string `json:"border_color"`

	// The Scryfall ID for the card back design present on this card.
	//
	// This value may be nil (even if scryfall doesn't say so)
	CardBackID *uuid.UUID `json:"card_back_id,omitempty"`

	// This card’s collector number. Note that collector numbers can contain non-numeric characters, such as letters or ★.
	CollectorNumber string `json:"collector_number"`

	// True if you should consider avoiding use of this print downstream.
	//
	// This value may be nil
	ContentWarning *bool `json:"content_warning,omitempty"`

	// True if this card was only released in a video game.
	Digital bool `json:"digital"`

	// An array of computer-readable flags that indicate if this card can come in foil, nonfoil, or etched finishes.
	Finishes []string `json:"finishes,omitempty"`

	// The just-for-fun name printed on the card (such as for Godzilla series cards).
	//
	// This value may be nil
	FlavorName *string `json:"flavor_name,omitempty"`

	// The flavor text, if any.
	//
	// This value may be nil
	FlavorText *string `json:"flavor_text,omitempty"`

	// This card’s frame effects, if any.
	//
	// This value may be nil
	FrameEffects []string `json:"frame_effects,omitempty"`

	// This card’s frame layout.
	Frame string `json:"frame"`

	// True if this card’s artwork is larger than normal.
	FullArt bool `json:"full_art"`

	// A list of games that this card print is available in, paper, arena, and/or mtgo.
	Games []string `json:"games,omitempty"`

	// True if this card’s imagery is high resolution.
	HighresImage bool `json:"highres_image"`

	// A unique identifier for the card artwork that remains consistent across reprints. Newly spoiled cards may not have this field yet.
	//
	// This value may be nil
	IllustrationID *uuid.UUID `json:"illustration_id,omitempty"`

	// A computer-readable indicator for the state of this card’s image, one of missing, placeholder, lowres, or highres_scan.
	ImageStatus string `json:"image_status"`

	// An object listing available imagery for this card. See the Card Imagery article for more information.
	//
	// This value may be nil
	ImageUris map[string]string `json:"image_uris,omitempty"`

	// True if this card is oversized.
	Oversized bool `json:"oversized"`

	// An object containing daily price information for this card, including usd, usd_foil, usd_etched, eur, eur_foil, eur_etched, and tix prices, as strings.
	Prices Prices `json:"prices"`

	// The localized name printed on this card, if any.
	//
	// This value may be nil
	PrintedName *string `json:"printed_name,omitempty"`

	// The localized text printed on this card, if any.
	//
	// This value may be nil
	PrintedText *string `json:"printed_text,omitempty"`

	// The localized type line printed on this card, if any.
	//
	// This value may be nil
	PrintedTypeLine *string `json:"printed_type_line,omitempty"`

	// True if this card is a promotional print.
	Promo bool `json:"promo"`

	// An array of strings describing what categories of promo cards this card falls into.
	//
	// This value may be nil
	PromoTypes []string `json:"promo_types,omitempty"`

	// An object providing URIs to this card’s listing on major marketplaces. Omitted if the card is unpurchaseable.
	//
	// This value may be nil
	PurchaseURIs map[string]string `json:"purchase_uris,omitempty"`

	// This card’s rarity. One of common, uncommon, rare, special, mythic, or bonus.
	Rarity string `json:"rarity"`

	// An object providing URIs to this card’s listing on other Magic: The Gathering online resources.
	RelatedURIs map[string]string `json:"related_uris"`

	// The date this card was first released.
	ReleasedAt string `json:"released_at"`

	// True if this card is a reprint.
	Reprint bool `json:"reprint"`

	// A link to this card’s set on Scryfall’s website.
	ScryfallSetURI string `json:"scryfall_set_uri"`

	// This card’s full set name.
	SetName string `json:"set_name"`

	// A link to where you can begin paginating this card’s set on the Scryfall API.
	SetSearchURI string `json:"set_search_uri"`

	// The type of set this printing is in.
	SetType string `json:"set_type"`

	// A link to this card’s set object on Scryfall’s API.
	SetURI string `json:"set_uri"`

	// This card’s set code.
	Set string `json:"set"`

	// This card’s Set object uuid.UUID.
	SetID uuid.UUID `json:"set_id"`

	// True if this card is a Story Spotlight.
	StorySpotlight bool `json:"story_spotlight"`

	// True if the card is printed without text.
	Textless bool `json:"textless"`

	// Whether this card is a variation of another printing.
	Variation bool `json:"variation"`

	// The printing ID of the printing this card is a variation of.
	//
	// This value may be nil
	VariationOf *uuid.UUID `json:"variation_of,omitempty"`

	// The security stamp on this card, if any. One of oval, triangle, acorn, circle, arena, or heart.
	//
	// This value may be nil
	SecurityStamp *string `json:"security_stamp,omitempty"`

	// This card’s watermark, if any.
	//
	// This value may be nil
	Watermark *string `json:"watermark,omitempty"`

	// Undocumented field
	//
	// This value may be nil
	Foil *bool `json:"foil,omitempty"`

	// Undocumented field
	//
	// This value may be nil
	NonFoil *bool `json:"nonfoil,omitempty"`

	// Information about when this card was first previewed
	//
	// This value may be nil
	Preview *Preview `json:"preview,omitempty"`
}
type RelatedCard struct {
	// An unique ID for this card in Scryfall’s database.
	ID uuid.UUID `json:"id"`

	// A content type for this object, always related_card.
	Object string `json:"object"`

	// A field explaining what role this card plays in this relationship, one of token, meld_part, meld_result, or combo_piece.
	Component string `json:"component"`

	// The name of this particular related card.
	Name string `json:"name"`

	// The type line of this card.
	TypeLine string `json:"type_line"`

	// A URI where you can retrieve a full object describing this card on Scryfall’s API.
	URI string `json:"uri"`
}

type CardFace struct {
	// The name of the illustrator of this card face. Newly spoiled cards may not have this field yet.
	//
	// This value may be nil
	Artist *string `json:"artist,omitempty"`

	// The ID of the illustrator of this card face. Newly spoiled cards may not have this field yet.
	//
	// This value may be nil
	ArtistID *uuid.UUID `json:"artist_id,omitempty"`

	// The mana value of this particular face, if the card is reversible.
	//
	// This value may be nil (Bad docs)
	Cmc *float32 `json:"cmc,omitempty"`

	// The colors in this face’s color indicator, if any.
	//
	// This value may be nil
	ColorIndicator *Colors `json:"color_indicator,omitempty"`

	// This face’s colors, if the game defines colors for the individual face of this card.
	//
	// This value may be nil
	Colors *Colors `json:"colors,omitempty"`

	// This face’s defense, if any.
	//
	// This value may be nil
	Defense *string `json:"defense,omitempty"`

	// The flavor text printed on this face, if any.
	//
	// This value may be nil
	FlavorText *string `json:"flavor_text,omitempty"`

	// A unique identifier for the card face artwork that remains consistent across reprints. Newly spoiled cards may not have this field yet.
	//
	// This value may be nil
	IllustrationID *string `json:"illustration_id,omitempty"`

	// An object providing URIs to imagery for this face, if this is a double-sided card. If this card is not double-sided, then the imageUris property will be part of the parent object instead.
	//
	// This value may be nil
	ImageUris map[string]string `json:"image_uris,omitempty"`

	// The layout of this card face, if the card is reversible.
	//
	// This value may be nil
	Layout *string `json:"layout,omitempty"`

	// This face’s loyalty, if any.
	//
	// This value may be nil
	Loyalty *string `json:"loyalty,omitempty"`

	// The mana cost for this face. This value will be any empty string "" if the cost is absent. Remember that per the game rules, a missing mana cost and a mana cost of {0} are different values.
	ManaCost string `json:"mana_cost"`

	// The name of this particular face.
	Name string `json:"name"`

	// A content type for this object, always cardFace.
	Object string `json:"object"`

	// The Oracle ID of this particular face, if the card is reversible.
	//
	// This value may be nil
	OracleID *uuid.UUID `json:"oracle_id,omitempty"`

	// The Oracle text for this face, if any.
	//
	// This value may be nil
	OracleText *string `json:"oracle_text,omitempty"`

	// This face’s power, if any. Note that some cards have powers that are not numeric, such as *.
	//
	// This value may be nil
	Power *string `json:"power,omitempty"`

	// The localized name printed on this face, if any.
	//
	// This value may be nil
	PrintedName *string `json:"printed_name,omitempty"`

	// The localized text printed on this face, if any.
	//
	// This value may be nil
	PrintedText *string `json:"printed_text,omitempty"`

	// The localized type line printed on this face, if any.
	//
	// This value may be nil
	PrintedTypeLine *string `json:"printed_type_line,omitempty"`

	// This face’s toughness, if any.
	//
	// This value may be nil
	Toughness *string `json:"toughness,omitempty"`

	// The type line of this particular face, if the card is reversible.
	//
	// This value may be nil
	TypeLine *string `json:"type_line,omitempty"`

	// The watermark on this particulary card face, if any.
	//
	// This value may be nil
	Watermark *string `json:"watermark,omitempty"`
}

type Preview struct {
	// The date this card was previewed.
	//
	// This value may be nil
	PreviewedAt *string `json:"previewed_at,omitempty"`

	// A link to the preview for this card.
	//
	// This value may be nil
	SourceURI *string `json:"source_uri,omitempty"`

	// The name of the source that previewed this card.
	//
	// This value may be nil
	Source *string `json:"source,omitempty"`
}

type Prices struct {
	Usd       *string `json:"usd"`
	UsdFoil   *string `json:"usd_foil"`
	UsdEtched *string `json:"usd_etched"`
	Eur       *string `json:"eur"`
	EurFoil   *string `json:"eur_foil"`
	// Hidden / Removed (Bad Docs)
	EurEtched *string `json:"eur_etched,omitempty"`
	Tix       *string `json:"tix"`
}
