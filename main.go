package main

import (
	"flag"
	"log"

	"github.com/urfave/negroni"
	"github.com/xoriath/alexandria-go/fetch"
	"github.com/xoriath/alexandria-go/index"

	"fmt"
	"net/http"
)

var mainIndex = flag.String("main-index", "http://content.alexandria.atmel.com/meta/index.xml",
	"Provide the name of the main index to use. Leave undefined to download from the content server.")
var fetchKeywords = flag.Bool("fetch-keywords", true,
	"Fetch the keyword indexes.")
var preparedKeywordStore = flag.String("prepared-keyword-store", "",
	"Point initally to a already populated keyword store.")
var keywordStorePrefix = flag.String("keyword-store-prefix", "keywords",
	"Prefix of the keyword store data base.")
var keywordStoreExtension = flag.String("keyword-store-extension", "db",
	"Extension of the keyword store data base.")
var redirectPattern = flag.String("content-redirect-pattern", "http://content.alexandria.atmel.com/webhelp/{{.Book}}/index.html?{{.Topic}}",
	"Redirect pattern for content lookups. 2 replacement parameters, first is Book GUID and second is Topic GUID.")
var f1FragmentPattern = flag.String("f1-fragment-pattern", "http://content.alexandria.atmel.com/meta/f1/{{.Id}}-{{.Language}}-{{.Version}}.xml",
	"Pattern for the F1 fragments")

func main() {

	flag.Parse()
	books, err := fetch.MainIndex(*mainIndex)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Read main index,", len(books.Books), "books are available")
	}

	var store index.Store
	if *preparedKeywordStore != "" {
		store = index.OldStore(*preparedKeywordStore, *f1FragmentPattern)
	} else {
		store = index.NewStore(*keywordStorePrefix, *keywordStoreExtension, *f1FragmentPattern)
	}

	if *fetchKeywords {
		store = fetch.F1Indexes(books, store)
	}

	keywordStatistics := store.GetStatistics()
	fmt.Println("Index store using ", store.FileName, "with", keywordStatistics.KeywordCount, "keywords covering", keywordStatistics.NumberOfFiles, "files")

	mux := createRoutes(books, &store, *mainIndex, *redirectPattern)
	n := negroni.Classic()
	logger := negroni.NewLogger()

	n.Use(logger)
	n.UseHandler(mux)

	serverAddress := ":3001"
	fmt.Println("Server running, listening on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, n))
}
