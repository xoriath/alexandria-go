package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria/types"
)

type Catalog struct {
	Catalogs []string
	template *template.Template
}

func NewCatalogHandler(books *types.Books) *Catalog {

	catalogsMap := make(map[string]struct{})

	for _, book := range books.Books {
		for _, product := range book.Products {
			catalogsMap[product.Name] = struct{}{}
		}
	}

	var catalogs []string
	for key := range catalogsMap {
		catalogs = append(catalogs, key)
	}

	t := template.Must(template.ParseFiles("./templates/catalog.html"))

	return &Catalog{Catalogs: catalogs, template: t}
}

// CatalogHandler handles the
func (c *Catalog) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := c.template.Execute(w, c.Catalogs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
