package main

import (
	"log"

	"github.com/xoriath/alexandria-go/handlers"
	"github.com/xoriath/alexandria-go/index"
	"github.com/xoriath/alexandria-go/types"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gopkg.in/cheggaaa/pb.v1"

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

func fetchIndexes(books *types.Books, index index.Store) index.Store {
	var wg sync.WaitGroup
	wg.Add(len(books.Books))

	progressBar := pb.New(len(books.Books)).Prefix("Fetching indexes ").Start()

	for _, book := range books.Books {
		index.FetchIndex(&book, &wg, progressBar)
	}

	wg.Wait()
	progressBar.Finish()

	return index
}

func main() {

	fmt.Println("Fetching main index file...")
	books, err := fetchMain()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Fetched main index,", len(books.Books), "books are available")
	}

	mux := mux.NewRouter()
	mux.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.Handle("/", handlers.NewRootHandler(books)).Methods("GET")
	mux.Handle("/catalogs", handlers.NewCatalogHandler(books)).Methods("GET")
	mux.Handle("/catalogs/{product}", handlers.NewCatalogLocalesHandler(books)).Methods("GET")
	mux.Handle("/catalogs/{product}/{locale}", handlers.NewProductHandler(books)).Methods("GET")

	n := negroni.Classic()
	logger := negroni.NewLogger()

	n.Use(logger)
	n.UseHandler(mux)

	_ = fetchIndexes(books, index.NewStore("keywords", ".db"))

	serverAddress := ":3001"
	fmt.Println("Server running, listening on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, n))
}
