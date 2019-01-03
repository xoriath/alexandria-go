package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/xoriath/alexandria-go/types"
)

// Root is the HTTP handler for the root page. It contains an overview of all content available on the server.
type Root struct {
	books             *types.Books
	templateFunctions template.FuncMap
}

// NewRootHandler create the handler for the root page
func NewRootHandler(books *types.Books) *Root {
	functions := template.FuncMap{
		"prettySize": prettySize,
	}

	return &Root{books: books, templateFunctions: functions}
}

// CatalogHandler handles the
func (r *Root) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	t := template.Must(template.New("root.html").Funcs(r.templateFunctions).ParseFiles("./templates/root.gohtml"))
	err := t.Execute(w, r.books)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func prettySize(size int) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}
