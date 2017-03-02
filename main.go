package main

import (
	"github.com/xoriath/alexandria/types"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"encoding/xml"

	"fmt"
	"net/http"

	"sync"
)

var db *sql.DB

func getDb() (db *sql.DB, fresh bool) {
	if db == nil {
		var err error

		db, err = sql.Open("sqlite3", "./keywords.db")
		if err != nil {
			panic(err)
		}

		fresh = true
	}

	fresh = false
	return
}

func prepareDb() {
	db, fresh := getDb()

	if fresh {
		stmt, err := db.Prepare("CREATE TABLE keywords ( keyword varchar(255) PRIMARY KEY, book varchar(49), file varchar(46) )")

		_, err = stmt.Exec()
		if err != nil {
			panic(err)
		}
	}
}

func insertIndexes(indexes *types.Indexes) {
	db, _ := getDb()

	stmt, err := db.Prepare("INSERT OR REPLACE INTO keywords(keyword, book, file) values(?, ?, ?)")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	for _, index := range indexes.Keywords {
		fmt.Printf("Insert %v (%v:%v)\n", index.Keyword, indexes.BookID, index.File)
		_, err := tx.Stmt(stmt).Exec(index.Keyword, indexes.BookID, index.File)
		if err != nil {
			tx.Rollback()
			fmt.Println("Panic for", indexes.BookID)
		}
	}

	fmt.Println("Commit", indexes.BookID)
	tx.Commit()
}

func fetchIndex(id, version, language string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("http://content.alexandria.atmel.com/meta/f1/%v-%v-%v.xml", id, language, version)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	indexes := new(types.Indexes)

	err = decoder.Decode(indexes)
	if err != nil {
		fmt.Println("Failed to parse", url)
		panic(err)
	}

	insertIndexes(indexes)
}

func main() {
	fmt.Println("Prepare DB")
	prepareDb()

	fmt.Println("Fetch index.xml")
	resp, err := http.Get("http://content.alexandria.atmel.com/meta/index.xml")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	books := new(types.Books)

	err = decoder.Decode(books)
	if err != nil {
		panic(err)
	}

	fmt.Println("Fetched", len(books.Books), "books")

	var wg sync.WaitGroup

	for _, book := range books.Books {
		wg.Add(1)
		fetchIndex(book.ID, book.Version, book.Language, &wg)
	}

	wg.Wait()
	fmt.Println("Done")
}
