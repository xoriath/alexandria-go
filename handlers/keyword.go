package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/index"
)

// Keyword is a HTTP handler for keyword lookups
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
	URL     string `json:"url"`
	BookID  string `json:"book_id"`
	TopicID string `json:"topic_id"`
}

func (k *Keyword) mapResultToResponse(keywordResults []index.KeywordResult) []urlResponseType {
	var jsonResults []urlResponseType
	for _, keywordResult := range keywordResults {
		topic := strings.TrimSuffix(keywordResult.Filename, filepath.Ext(keywordResult.Filename))
		parts := map[string]string{"Book": keywordResult.BookID, "Topic": topic}

		var url bytes.Buffer
		if err := k.contentRedirectTemplate.Execute(&url, parts); err == nil {
			jsonResults = append(jsonResults, urlResponseType{URL: url.String(), BookID: keywordResult.BookID, TopicID: topic})
		}
	}

	return jsonResults
}

func (k *Keyword) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	keyword := vars["keyword"]
	keywordResults := k.store.LookupKeyword(keyword)
	jsonResults := k.mapResultToResponse(keywordResults)

	if k.redirect {
		http.Redirect(w, r, jsonResults[0].URL, http.StatusTemporaryRedirect)
	} else {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(jsonResults)
	}
}
