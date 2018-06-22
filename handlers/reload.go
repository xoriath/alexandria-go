package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/xoriath/alexandria-go/fetch"
	"github.com/xoriath/alexandria-go/index"
	"github.com/xoriath/alexandria-go/types"
)

// ReloadBook HTTP handler to reload books
type ReloadBook struct {
	books *types.Books
	index string
}

// NewReloadBookHandler create new HTTP handler to reload the books
func NewReloadBookHandler(books *types.Books, index string) *ReloadBook {
	return &ReloadBook{books: books, index: index}
}

func (rb *ReloadBook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tempBooks, err := fetch.MainIndex(rb.index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		rb.books = tempBooks
		json.NewEncoder(w).Encode(rb.books)
	}
}

// ReloadKeyword
type ReloadKeyword struct {
	books     *types.Books
	store     *index.Store
	f1Pattern string
}

// NewReloadKeywordHandler
func NewReloadKeywordHandler(books *types.Books, store *index.Store, f1Pattern string) *ReloadKeyword {
	return &ReloadKeyword{books: books, store: store, f1Pattern: f1Pattern}
}

func (rk *ReloadKeyword) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tempStore := fetch.F1Indexes(rk.books, index.NewStore("keywords", ".db", rk.f1Pattern))

	rk.store = &tempStore
	stat := rk.store.GetStatistics()
	json.NewEncoder(w).Encode(stat)
}
