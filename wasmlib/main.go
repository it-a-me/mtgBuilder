package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"syscall/js"
	"time"

	"mtgBuilder/card"
	"mtgBuilder/query"
)

func NewError(err error) js.Value {
	return js.Global().Call("Error", err.Error())
}

var errInvalidArgument = errors.New("invalid argument")

func CheckArgs(args []js.Value, expected []js.Type) error {
	var got []js.Type
	for _, arg := range args {
		got = append(got, arg.Type())
	}
	if !slices.Equal(got, expected) {
		return fmt.Errorf("%w: expected %v, got %v", errInvalidArgument, expected, got)
	}
	return nil
}

var cards []card.Card

func feedCards(_ js.Value, args []js.Value) any {
	log.Println("feeding cards")
	if err := CheckArgs(args, []js.Type{js.TypeString}); err != nil {
		return NewError(err)
	}
	
	start := time.Now()

	var c []card.Card
	if err := json.Unmarshal([]byte(args[0].String()), &c); err != nil {
		return NewError(err)
	}
	log.Printf("parsed cards.json in %s", time.Since(start).String())
	cards = c
	return nil
}

func parseQuery(_ js.Value, args []js.Value) any {
	if err := CheckArgs(args, []js.Type{js.TypeString}); err != nil {
		return NewError(err)
	}
	queryLine := args[0].String()
	q, err := query.Parse(queryLine, true)
	if err != nil {
		return NewError(err)
	}
	j, err := json.Marshal(q)
	if err != nil {
		return NewError(err)
	}
	return string(j)
}

func queryCards(_ js.Value, args []js.Value) any {
	if err := CheckArgs(args, []js.Type{js.TypeString}); err != nil {
		return NewError(err)
	}
	query, err := query.Parse(args[0].String(), true)
	if err != nil {
		return NewError(err)
	}
	var matches []any
	for i, c := range cards {
		if query.Matches(&c) {
			matches = append(matches, i)
		}
	}
	return matches
}

var ErrIndexOutOfBounds = errors.New("index out of bounds")

func getCard(_ js.Value, args []js.Value) any {
	if err := CheckArgs(args, []js.Type{js.TypeNumber}); err != nil {
		return NewError(err)
	}
	i := args[0].Int()
	if len(cards) < i {
		return NewError(fmt.Errorf("%w: %d/%d", ErrIndexOutOfBounds, i, len(cards)))
	}
	c := cards[i]
	bytes, err := json.Marshal(c)
	if err != nil {
		return NewError(err)
	}
	return string(bytes)
}

const exportName = "GO_cardQuery"

func main() {
	g := js.Global()
	exports := map[string]any{
		"feedCards":  js.FuncOf(feedCards),
		"parseQuery": js.FuncOf(parseQuery),
		"queryCards": js.FuncOf(queryCards),
		"getCard":    js.FuncOf(getCard),
	}
	g.Set(exportName, exports)
	log.Printf("exported:\n%+v", exports)
	<-make(chan struct{})
}
