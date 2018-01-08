package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/xoriath/alexandria-go/index"
)

type Query struct {
	regexp *regexp.Regexp
}

// NewQueryHandler create a new QueryHandler
func NewQueryHandler(store *index.Store) *Query {
	return &Query{regexp: regexp.MustCompile("k\\(([^\\)]+)\\)")}
}

func (q *Query) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// query string is on form:
	//   appId=AS70&l=EN-US&k=k(atmel;device%3AATmega128A;register%3AEIMSK);k(DevLang-c)
	//           &rd=true

	// and we want a list of all k(<query>) parts of the string. k\(([^\)]+)\)

	appId := r.FormValue("appId")
	language := r.FormValue("language")
	rawKeywords := r.FormValue("keywords")
	redirect := r.FormValue("redirect")

	keywords := q.regexp.FindAllString(rawKeywords, -1)
	if keywords == nil {
		http.Error(w, fmt.Sprintf("Failed to parse '%v' into queries", rawKeywords), http.StatusBadRequest)
		return
	}

}
