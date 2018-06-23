package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/types"
)

type Resource struct {
	books    *types.Books
	resource string

	contentRedirectTemplate *template.Template
}

func NewResourceHandler(books *types.Books, resource, redirectPattern string) *Resource {
	return &Resource{books: books, resource: resource, contentRedirectTemplate: template.Must(template.New("").Parse(redirectPattern))}
}

func (r *Resource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	guid, foundGUID := vars["guid"]
	language, foundLanguage := vars["language"]
	version, foundVersion := vars["version"]

	if !foundGUID || !foundLanguage || !foundVersion {
		http.Error(w, fmt.Sprintf("Missing parameters. Validation: GUID=%v Language=%v Version=%v", foundGUID, foundLanguage, foundVersion), http.StatusBadRequest)
		return
	}

	doesResourceExist := false
	for _, book := range r.books.Books {
		if guid == book.ID && language == book.Language && version == book.Version {
			doesResourceExist = true
			break
		}
	}

	if !doesResourceExist {
		http.Error(w, fmt.Sprintf("No content known for GUID=%v Language=%v Version=%v", guid, language, version), http.StatusNotFound)
		return
	}

	parts := map[string]string{"ResourceType": r.resource, "Id": guid, "Language": language, "Version": version}
	var url bytes.Buffer
	if err := r.contentRedirectTemplate.Execute(&url, parts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	urlString := url.String()

	switch r.resource {
	case "cab":
		urlString += ".cab"
	}

	http.Redirect(w, req, urlString, http.StatusTemporaryRedirect)
}
