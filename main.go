package main

import (
	"github.com/xoriath/alexandria/types"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"encoding/xml"

	"fmt"
	"net/http"
)

func getDb() (db *sql.DB, fresh bool) {
	var err error

	db, err = sql.Open("sqlite3", "./keywords.db")
	if err != nil {
		panic(err)
	}

	fresh = true
	return
}

func prepareDb() {
	db, fresh := getDb()

	if fresh {
		createTableStmt := `
			CREATE TABLE files (
				fileid		INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				book 		VARCHAR(49) NOT NULL,
				file		VARCHAR(46) NOT NULL,

				CONSTRAINT file_constraint UNIQUE (book, file)
			);

			CREATE TABLE keywords ( 
				keyword 	TEXT NOT NULL, 
				file		INTEGER NOT NULL,
				
				FOREIGN KEY(file) REFERENCES files(fileid)
			);
			
			CREATE INDEX keywords_index ON keywords(keyword);
			`

		_, err := db.Exec(createTableStmt)

		if err != nil {
			panic(err)
		}
	}
}

func findID(tx *sql.Tx, bookID, file string) (id int64) {

	query := fmt.Sprintf("SELECT fileid FROM files WHERE book = '%v' AND file = '%v'", bookID, file)

	rows, err := tx.Query(query)

	if err != nil {
		fmt.Println(query)
		panic(err)
	}

	for rows.Next() {
		err := rows.Scan(&id)

		if err != nil {
			panic(err)
		}

		break
	}

	fmt.Println("Fetched existing id", id)

	return
}

func insertIndexes(indexes *types.Indexes) {
	db, _ := getDb()

	bookStmt, err := db.Prepare("INSERT INTO files(book, file) values(?, ?)")
	if err != nil {
		panic(err)
	}

	keywordStmt, err := db.Prepare("INSERT INTO keywords(keyword, file) values(?, ?)")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	idMap := make(map[string]int64)

	for _, index := range indexes.Keywords {
		fmt.Printf("Insert %v (%v:%v)\n", index.Keyword, indexes.BookID, index.File)

		// Check if we have this id already
		id := idMap[indexes.BookID+index.File]
		if id == 0 {
			// Try to insert, will fail if exists due to unique constraint
			res, err := tx.Stmt(bookStmt).Exec(indexes.BookID, index.File)
			if err != nil {
				// unique failed, time to search
				id = findID(tx, indexes.BookID, index.File)
			} else {
				// fetch id that we just inserted
				id, err = res.LastInsertId()
				if err != nil {
					tx.Rollback()
					panic(err)
				}
			}
		}

		idMap[indexes.BookID+index.File] = id

		_, err := tx.Stmt(keywordStmt).Exec(index.Keyword, id)
		if err != nil {
			tx.Rollback()
			fmt.Println("Panic for", indexes.BookID, id)
			panic(err)
		}
	}

	fmt.Println("Commit", indexes.BookID)
	tx.Commit()
}

func fetchIndex(id, version, language string) {

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

	for _, book := range books.Books {
		fetchIndex(book.ID, book.Version, book.Language)
	}

	fmt.Println("Done")
}
