package handlers

import (
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria/types"
)

type Root struct {
	books *types.Books

	//template *template.Template
	templatePath string
}

// NewRootHandler create the handler for the root page
func NewRootHandler(books *types.Books) *Root {

	//t := template.Must(template.ParseFiles("./templates/root.html"))

	return &Root{books: books, templatePath: "./templates/root.html"}
}

// CatalogHandler handles the
func (r *Root) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	t := template.Must(template.ParseFiles(r.templatePath))
	err := t.Execute(w, r.books)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
