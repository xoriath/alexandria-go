package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Resource struct {
	resource string
}

func NewResourceHandler(resource string) *Resource {
	return &Resource{resource: resource}
}

func (r *Resource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	guid, foundGUID := vars["guid"]
	language, foundLanguage := vars["language"]
	version, foundVersion := vars["version"]

	if !foundGUID || !foundLanguage || !foundVersion {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
	} else {
		url := fmt.Sprintf("http://content.alexandria.atmel.com/%v/%v-%v-%v", r.resource, guid, language, version)
		http.Redirect(w, req, url, http.StatusMovedPermanently)
	}
}
