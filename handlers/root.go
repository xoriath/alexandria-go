package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria-go/types"
)

type Root struct {
	books *types.Books
}

// NewRootHandler create the handler for the root page
func NewRootHandler(books *types.Books) *Root {
	return &Root{books: books}
}

// CatalogHandler handles the
func (r *Root) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	t := template.Must(template.ParseFiles("./templates/root.html"))
	err := t.Execute(w, r.books)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
