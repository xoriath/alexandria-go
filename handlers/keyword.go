package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/index"
)

type Keyword struct {
	store                   *index.Store
	redirect                bool
	contentRedirectTemplate *template.Template
}

// NewKeywordHandler create a new KeywordHandler
func NewKeywordHandler(store *index.Store, redirectPattern string) *Keyword {
	return &Keyword{
		store:                   store,
		redirect:                false,
		contentRedirectTemplate: template.Must(template.New("").Parse(redirectPattern))}
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
		parts := map[string]string{"Book": result.BookID, "Topic": strings.TrimSuffix(result.Filename, filepath.Ext(result.Filename))}

		var url bytes.Buffer
		if err := k.contentRedirectTemplate.Execute(&url, parts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if k.redirect {
			http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
		} else {
			urlResponse := &urlResponseType{URL: url.String()}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(urlResponse)
		}
	}
}
