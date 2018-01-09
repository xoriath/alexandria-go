package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/types"
)

type Resource struct {
	books    *types.Books
	resource string
}

func NewResourceHandler(books *types.Books, resource string) *Resource {
	return &Resource{books: books, resource: resource}
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

	url := fmt.Sprintf("http://content.alexandria.atmel.com/%v/%v-%v-%v", r.resource, guid, language, version)

	switch r.resource {
	case "cab":
		url += ".cab"
	}

	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}
