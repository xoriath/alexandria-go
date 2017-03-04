package main

import (
	"github.com/xoriath/alexandria/handlers"
	"github.com/xoriath/alexandria/index"
	"github.com/xoriath/alexandria/types"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"encoding/xml"

	"fmt"
	"net/http"
	"sync"
)

func fetchMain() (*types.Books, error) {

	resp, err := http.Get("http://content.alexandria.atmel.com/meta/index.xml")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	books := new(types.Books)

	err = decoder.Decode(books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func fetchIndexes(books *types.Books, index index.IndexStore) index.IndexStore {
	var wg sync.WaitGroup
	for _, book := range books.Books {
		wg.Add(1)
		index.FetchIndex(book.ID, book.Version, book.Language, &wg)
	}
	wg.Wait()

	return index
}

func main() {

	books, err := fetchMain()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Fetched", len(books.Books), "books")
	}

	//index := fetchIndexes(books, index.New("keywords", ".db"))

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.Handle("/catalogs", handlers.NewCatalogHandler(books)).Methods("GET")
	router.Handle("/catalogs/{product}", handlers.NewCatalogLocalesHandler(books)).Methods("GET")
	router.Handle("/catalogs/{product}/{locale}", handlers.NewProductHandler(books)).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(router)

	http.ListenAndServe(":3001", n)
}
