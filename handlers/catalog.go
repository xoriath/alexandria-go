package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria-go/types"
)

// Catalog describes the set of catalogs that are available.
//
// This is used to separate consuming applications, i.e AtmelStudio70
type Catalog struct {
	Products []string
}

// NewCatalogHandler creates a new HTTP handler for the catalogs
func NewCatalogHandler(books *types.Books) *Catalog {

	products := books.Products()

	return &Catalog{Products: products}
}

// CatalogHandler handles the
func (c *Catalog) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	t := template.Must(template.ParseFiles("./templates/catalog.html"))
	err := t.Execute(w, c.Products)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
