package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

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
		handlers.NewResourceHandler(books, "cab", *contentRedirectPattern)).Methods("GET")
	mux.Handle("/package/{guid:GUID-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+-[A-Z0-9]+}/{version:[0-9]+}/{language:[a-zA-Z]+-[a-zA-Z]+}",
		handlers.NewResourceHandler(books, "package", *contentRedirectPattern)).Methods("GET")

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

	// Reverse to the content server
	contentBaseURL, err := url.Parse(*contentServerBase)
	if err != nil {
		panic(err)
	}

	contentProxy := httputil.NewSingleHostReverseProxy(contentBaseURL)
	contentProxy.ModifyResponse = func(r *http.Response) error {
		if r.StatusCode != http.StatusOK {
			replaceResponseBody(r, fmt.Sprintf("<html><header><title>Content not found (%d).</title></header><body>Content not found (%d).</body></html>", r.StatusCode, r.StatusCode))
		}

		return nil
	}

	mux.PathPrefix("/content/").Handler(
		http.StripPrefix("/content/", contentProxy))
	mux.HandleFunc("/content", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/content/index.html", http.StatusTemporaryRedirect)
	})

	return mux
}

func replaceResponseBody(r *http.Response, message string) {
	r.Body.Close()
	body := ioutil.NopCloser(strings.NewReader(message))
	r.Body = body
	r.ContentLength = int64(len(message))
	r.Header.Set("Content-Length", strconv.Itoa(len(message)))
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
}
