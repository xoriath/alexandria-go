package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/index"
)

type Keyword struct {
	store    *index.Store
	redirect bool
}

// NewKeywordHandler create a new KeywordHandler
func NewKeywordHandler(store *index.Store) *Keyword {
	return &Keyword{store: store, redirect: false}
}

// Redirect set the keyword handler to redirect via HTTP Temporary redirect to the result page
func (k *Keyword) Redirect() *Keyword {
	k.redirect = true
	return k
}

// NoRedirect sets the keyword handler to not redirect to the result page, but return the data to the client
func (k *Keyword) NoRedirect() *Keyword {
	k.redirect = false
	return k
}

type urlResponseType struct {
	URL string `json:"url"`
}

func (k *Keyword) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	keyword := vars["keyword"]

	keywordResults := k.store.LookupKeyword(keyword)

	if len(keywordResults) == 0 {
		http.Error(w, fmt.Sprintf("No results for query '%v'", keyword), http.StatusNotFound)
	} else {
		result := keywordResults[0]
		url := fmt.Sprintf("http://content.alexandria.atmel.com/webhelp/%v/index.html?%v", result.BookID, strings.TrimSuffix(result.Filename, filepath.Ext(result.Filename)))

		if k.redirect {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			urlResponse := &urlResponseType{URL: url}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(urlResponse)
		}
	}
}
