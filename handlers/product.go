package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria-go/types"
)

type Product struct {
	books *types.Books
}

type productInfo struct {
}

func NewProductHandler(books *types.Books) *Product {
	return &Product{books: books}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// vars := mux.Vars(req)
	// product := vars["product"]
	// locale := vars["locale"]

	t := template.Must(template.ParseFiles("./templates/product.html"))
	err := t.Execute(w, p.books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
