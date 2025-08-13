package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"syscall/js"
	"time"

	"github.com/goccy/go-json"

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

func NewPromise(global js.Value, f func() (any, error)) js.Value {
	typePromise := global.Get("Promise")
	return typePromise.New(js.FuncOf(func(g js.Value, args []js.Value) any {
		log.Printf("NewPromise invoked with %+v", args)
		resolve := args[0]
		reject := args[1]
		go func() {
			v, err := f()
			_ = v
			if err != nil {
				reject.Invoke(NewError(err))
			}
			resolve.Invoke(nil)
		}()
		return nil
	}))
}

func WrapAsync(global js.Value, f func(global js.Value, args []js.Value) (any, error)) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		return NewPromise(global, func() (any, error) {
			return f(global, args)
		})
	})
}

var cards []card.Card

func feedCards(g js.Value, args []js.Value) (any, error) {
	log.Println("feeding cards")
	if err := CheckArgs(args, []js.Type{js.TypeString}); err != nil {
		return nil, err
	}
	url := args[0].String()
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	start := time.Now()

	var c []card.Card
	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}
	log.Printf("parsed cards.json in %s", time.Since(start).String())
	cards = c
	return nil, nil
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
		"feedCards":  WrapAsync(g, feedCards),
		"parseQuery": js.FuncOf(parseQuery),
		"queryCards": js.FuncOf(queryCards),
		"getCard":    js.FuncOf(getCard),
	}
	g.Set(exportName, exports)
	log.Printf("exported:\n%+v", exports)
	<-make(chan struct{})
}
