package main

import (
	"log"
	"os"

	"mtgBuilder/fetch"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please supply exactly 1 argument")
	}

	f, err := os.OpenFile(os.Args[1], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o0644)
	if err != nil {
		log.Fatal(err)
	}

	js, err := fetch.GetOracleCardsJSON()
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(js)
	if err != nil {
		log.Fatal(err)
	}
}
