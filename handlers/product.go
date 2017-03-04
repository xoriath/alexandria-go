package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria/types"
)

type Product struct {
	template *template.Template
}

type productPage struct {
}

func NewProductHandler(books *types.Books) *Product {

}

func (p *Product) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	product := vars["product"]
	locale := vars["locale"]

	locales := c.Catalogs[product]
	if locales == nil {
		http.Error(w, "404 No locale for "+product, http.StatusNotFound)
		return
	}

	err := p.template.Execute(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
