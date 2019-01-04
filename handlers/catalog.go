package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Catalog describes the set of catalogs that are available.
//
// This is used to separate consuming applications, i.e AtmelStudio70
type Catalog struct {
	Products []string
	template *template.Template
}

// Productslister can return a list of all known products
type Productslister interface {
	Products() []string
}

// NewCatalogHandler creates a new HTTP handler for the catalogs
func NewCatalogHandler(books Productslister, templateFolder string) *Catalog {

	products := books.Products()

	return &Catalog{
		Products: products,
		template: template.Must(template.ParseFiles(filepath.Join(templateFolder, "catalog.gohtml"))),
	}
}

// CatalogHandler handles the
func (c *Catalog) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("[catalog] Serving %d products", len(c.Products))
	if err := c.template.Execute(w, c.Products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
