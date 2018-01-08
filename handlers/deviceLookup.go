package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/index"
)

type DeviceLookup struct {
	store *index.Store
}

func NewDeviceLookupHandler(store *index.Store) *DeviceLookup {
	return &DeviceLookup{store: store}
}

// CabHandler handles the
func (d *DeviceLookup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	device := vars["device"]
	component, foundComponent := vars["component"]
	register, foundRegister := vars["register"]
	bitfield, foundBitfield := vars["bitfield"]

	query := fmt.Sprintf("atmel;device:%v", device)
	if foundComponent {
		query += fmt.Sprintf(";comp:%v", component)
	}
	if foundRegister {
		query += fmt.Sprintf(";register:%v", register)
	}
	if foundBitfield {
		query += fmt.Sprintf(";bitfield:%v", bitfield)
	}

	keywordResults := d.store.LookupKeyword(query)

	if len(keywordResults) == 0 {
		http.Error(w, fmt.Sprintf("No results for query '%v'", query), http.StatusNotFound)
	} else {
		result := keywordResults[0]
		url := fmt.Sprintf("http://content.alexandria.atmel.com/webhelp/%v/index.html?%v", result.BookID, strings.TrimSuffix(result.Filename, filepath.Ext(result.Filename)))
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
