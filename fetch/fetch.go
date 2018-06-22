package fetch

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/xoriath/alexandria-go/index"
	"github.com/xoriath/alexandria-go/types"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// MainIndex Fetch the main index file. Can be local on online.
func MainIndex(index string) (*types.Books, error) {

	if u, err := url.Parse(index); err != nil || u.Scheme == "" {
		reader, err := fetchMainIndexLocally(index)
		if err != nil {
			return nil, err
		}

		return parseMainIndex(reader)
	}

	fmt.Println("Fetching main index file...", index)
	resp, err := http.Get(index)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return parseMainIndex(resp.Body)
}

func fetchMainIndexLocally(index string) (io.Reader, error) {
	file, err := os.Open(index)
	if err != nil {
		return nil, err
	}

	fmt.Println("Opening main index file...", index)
	return bufio.NewReader(file), nil
}

func parseMainIndex(reader io.Reader) (*types.Books, error) {
	decoder := xml.NewDecoder(reader)
	books := new(types.Books)

	err := decoder.Decode(books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// F1Indexes Fetch the F1 indexes that correspond to the content of the books
func F1Indexes(books *types.Books, index index.Store) index.Store {
	var wg sync.WaitGroup
	wg.Add(len(books.Books))

	progressBar := pb.New(len(books.Books)).Prefix("Fetching F1 indexes ").Start()

	for _, book := range books.Books {
		index.FetchIndex(&book, &wg, progressBar)
	}

	wg.Wait()
	progressBar.Finish()

	return index
}
