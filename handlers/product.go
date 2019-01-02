package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/types"
)

// Product contains all the books that relates to the product.
type Product struct {
	books *types.Books
}

type productInfo struct {
}

// NewProductHandler creates the Product which is a HTTP handler. 
func NewProductHandler(books *types.Books) *Product {
	return &Product{books: books}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	product := vars["product"]
	locale := vars["locale"]

	data := struct {
		Host        string
		Product     string
		Locale      string
		ContentRoot *types.Books
	}{req.Host, product, locale, p.books}

	t := template.Must(template.ParseFiles("./templates/product.gohtml"))
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
