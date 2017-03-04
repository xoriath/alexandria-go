package main

import (
	"github.com/xoriath/alexandria/index"
	"github.com/xoriath/alexandria/types"

	"encoding/xml"

	"fmt"
	"net/http"
	"sync"
)

func main() {
	fmt.Println("Fetch index.xml")
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

	fmt.Println("Fetched", len(books.Books), "books")

	index := index.New("keywords", ".db")

	var wg sync.WaitGroup

	for _, book := range books.Books {
		wg.Add(1)
		index.FetchIndex(book.ID, book.Version, book.Language, &wg)
	}

	wg.Wait()
	fmt.Println("Done")

	lookup := index.LookupKeyword("atmel;device:atsaml21e15a;register:intenclr")

	fmt.Printf("Lookup: %+v", lookup)
}
