package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria/types"
)

type Product struct {
	template *template.Template
}

type productInfo struct {
	
}


func NewProductHandler(books *types.Books) *Product {
	return &Product{}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// vars := mux.Vars(req)
	// product := vars["product"]
	// locale := vars["locale"]

	// err := p.template.Execute(w)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}
