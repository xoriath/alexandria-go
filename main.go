package main

import (
	"alexandria/types"

	"encoding/xml"

	"net/http"
	"fmt"
)

func main() {

	resp, err := http.Get("http://content.alexandria.atmel.com/meta/index.xml")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)

	books := new(types.Books)

	err = decoder.Decode(books)
	if err != nil {
		panic(err)
	}

	for _, book := range books.Books {
		fmt.Printf("%v (%v:%v:%v)\n", book.Title, book.ID, book.Version, book.Language)
	}
}