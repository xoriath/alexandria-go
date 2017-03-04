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
	template *template.Template

	Catalogs map[string][]catalogLocales
}

type catalogLocalePage struct {
	Product string
	Locales []catalogLocales
}

func NewCatalogLocalesHandler(books *types.Books) *CatalogLocale {

	catalogsMap := make(map[string]map[string]struct{})

	for _, book := range books.Books {
		for _, product := range book.Products {
			if catalogsMap[product.Name] == nil {
				catalogsMap[product.Name] = make(map[string]struct{})
			}

			catalogsMap[product.Name][book.Language] = struct{}{}
		}
	}

	m := make(map[string][]catalogLocales)

	for product := range catalogsMap {
		var catalogs []catalogLocales

		for language := range catalogsMap[product] {
			catalogs = append(catalogs, catalogLocales{Product: product, Locale: language})
		}

		m[product] = catalogs
	}

	t := template.Must(template.ParseFiles("./templates/catalogLocales.html"))

	return &CatalogLocale{Catalogs: m, template: t}
}

func (c *CatalogLocale) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	product := vars["product"]

	locales := c.Catalogs[product]
	if locales == nil {
		http.Error(w, "404 No locale for "+product, http.StatusNotFound)
		return
	}

	err := c.template.Execute(w, &catalogLocalePage{Product: product, Locales: locales})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
