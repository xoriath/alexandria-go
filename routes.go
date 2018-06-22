package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xoriath/alexandria-go/handlers"
	"github.com/xoriath/alexandria-go/index"
	"github.com/xoriath/alexandria-go/types"
)

func createRoutes(books *types.Books, indexStore *index.Store, mainIndex, redirectPattern string) *mux.Router {
	mux := mux.NewRouter()

	// Static routes
	mux.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Root handler
	mux.Handle("/", handlers.NewRootHandler(books)).Methods("GET")

	// Microsoft Help Service Endpoints
	mux.Handle("/catalogs", handlers.NewCatalogHandler(books)).Methods("GET")
	mux.Handle("/catalogs/{product}", handlers.NewCatalogLocalesHandler(books)).Methods("GET")
	mux.Handle("/catalogs/{product}/{locale}", handlers.NewProductHandler(books)).Methods("GET")
	mux.Handle("/query/{query}", handlers.NewQueryHandler(indexStore)).Methods("GET").Queries("appId", "{appId}").Queries("l", "{language}").Queries("k", "keywords").Queries("rd", "redirect")

	// Endpoints serving CAB and package data
	mux.Handle("/cab/{guid:GUID-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+}-{language:[a-zA-Z]+-[a-zA-Z]+}-{version:[0-9]+}.cab",
		handlers.NewResourceHandler(books, "cab")).Methods("GET")
	mux.Handle("/package/{guid:GUID-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+}/{version:[0-9]+}/{language:[a-zA-Z]+-[a-zA-Z]+}",
		handlers.NewResourceHandler(books, "package")).Methods("GET")

	// Keyword endpoints
	mux.Handle("/keyword/{keyword}", handlers.NewKeywordHandler(indexStore, redirectPattern).NoRedirect()).Methods("GET")
	mux.Handle("/keyword/{keyword}/redirect", handlers.NewKeywordHandler(indexStore, redirectPattern).Redirect()).Methods("GET")

	// Device specific lookup endpoints
	mux.Handle("/device-lookup/{device}/register/{register}",
		handlers.NewDeviceLookupHandler(indexStore)).Methods("GET")
	mux.Handle("/device-lookup/{device}/register/{register}/bitfield/{bitfield}",
		handlers.NewDeviceLookupHandler(indexStore)).Methods("GET")
	mux.Handle("/device-lookup/{device}/component/{component}",
		handlers.NewDeviceLookupHandler(indexStore)).Methods("GET")
	mux.Handle("/device-lookup/{device}/component/{component}/register/{register}",
		handlers.NewDeviceLookupHandler(indexStore)).Methods("GET")
	mux.Handle("/device-lookup/{device}/component/{component}/register/{register}/bitfield/{bitfield}",
		handlers.NewDeviceLookupHandler(indexStore)).Methods("GET")

	// Reload endpoints
	mux.Handle("/reload/books", handlers.NewReloadBookHandler(books, mainIndex)).Methods("GET")
	mux.Handle("/reload/keywords", handlers.NewReloadKeywordHandler(books, indexStore, *f1FragmentPattern)).Methods("GET")

	return mux
}
