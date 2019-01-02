package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/xoriath/alexandria-go/index"
)

// Query is a HTTP handler for the MTPS query endpoint
type Query struct {
	store                   *index.Store
	contentRedirectTemplate *template.Template
}

// NewQueryHandler create a new QueryHandler
func NewQueryHandler(store *index.Store, redirectPattern string) *Query {
	return &Query{
		store: store,
		contentRedirectTemplate: template.Must(template.New("").Parse(redirectPattern))}
}

func (q *Query) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// query string is on form:
	//   appId=AS70&l=EN-US&k=k(VS.SolutionExplorer.Selection);k(VS.SolutionExplorer);k(VS.SolutionExplorer.Solutions)
	//           &rd=true

	//appId := r.FormValue("appId")
	//language := r.FormValue("l")
	redirect := r.FormValue("rd")

	keywords := collectKeywords(r)
	if len(keywords) == 0 {
		http.Error(w, "Failed to parse queries", http.StatusBadRequest)
		return
	}

	for _, keyword := range keywords {
		keywordResults := q.store.LookupKeyword(keyword)

		if len(keywordResults) == 0 {
			continue
		} else {
			result := keywordResults[0]
			parts := map[string]string{"Book": result.BookID, "Topic": strings.TrimSuffix(result.Filename, filepath.Ext(result.Filename))}

			var url bytes.Buffer
			if err := q.contentRedirectTemplate.Execute(&url, parts); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if redirect == "true" {
				http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
			} else { 
				http.Error(w, fmt.Sprintf("Redirect parameter is not 'true', but '%s'", redirect), http.StatusBadRequest)
			}

			return
		}
	}

	http.Error(w, fmt.Sprintf("No results for %v", keywords), http.StatusNotFound)
	return
}

func collectKeywords(r *http.Request) []string {
	var keywords []string
	for k, v := range r.Form {
		if k == "k" {
			keywords = append(keywords, trimKeywordPart(v[0]))
		} else if strings.HasPrefix(k, "k(") {
			keywords = append(keywords, trimKeywordPart(k))
		}
	}

	return keywords
}

func trimKeywordPart(k string) string {
	withoutPrefix := strings.TrimPrefix(k, "k(")
	withoutPrefixAndSuffix := strings.TrimSuffix(withoutPrefix, ")")
	return withoutPrefixAndSuffix
}
