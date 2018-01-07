package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria/types"
)

type catalogLocales struct {
	Product string
	Locale  string
}

type CatalogLocale struct {
	ProductLocales map[string][]string
}

type catalogLocalePage struct {
	Product string
	Locales []string
}

func NewCatalogLocalesHandler(books *types.Books) *CatalogLocale {

	productLocales := make(map[string][]string)
	for _, product := range books.Products() {
		productLocales[product] = books.Locales(product)
	}

	return &CatalogLocale{ProductLocales: productLocales}
}

func (c *CatalogLocale) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	product := vars["product"]

	locales := c.ProductLocales[product]
	if locales == nil {
		http.Error(w, "404 No locale for "+product, http.StatusNotFound)
		return
	}

	t := template.Must(template.ParseFiles("./templates/catalogLocales.html"))
	err := t.Execute(w, &catalogLocalePage{Product: product, Locales: locales})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
