package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria/types"
)

type Catalog struct {
	Products []string
}

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
