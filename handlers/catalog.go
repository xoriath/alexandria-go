package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria/types"
)

type Catalog struct {
	Products []string
	template *template.Template
}

func NewCatalogHandler(books *types.Books) *Catalog {

	products := books.Products()
	t := template.Must(template.ParseFiles("./templates/catalog.html"))

	return &Catalog{Products: products, template: t}
}

// CatalogHandler handles the
func (c *Catalog) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := c.template.Execute(w, c.Products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
