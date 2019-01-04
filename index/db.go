package index

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3" // use sqlite
	"github.com/xoriath/alexandria-go/types"

	"errors"
	"os"
)

// Store is the store for keywords
type Store struct {
	FileName string
	handle   *sql.DB

	prefix string
	ext    string

	indexWriteChan     chan *types.Indexes
	f1FragmentTemplate *template.Template
}

// NewStore creates a index store
func NewStore(prefix, ext, f1FragmentPattern string) *Store {
	store := &Store{
		prefix:             prefix,
		ext:                ext,
		f1FragmentTemplate: template.Must(template.New("").Parse(f1FragmentPattern))}

	store.prepareDb()

	store.indexWriteChan = store.insertIndexes()

	return store
}

// OldStore open a already existing index store
func OldStore(file, f1FragmentPattern string) *Store {
	store := &Store{
		FileName:           file,
		f1FragmentTemplate: template.Must(template.New("").Parse(f1FragmentPattern))}

	store.prepareDb()

	return store
}

func (i *Store) getDbFile(prefix, ext string) (string, error) {

	if i.FileName == "" {

		for j := 0; j < 10000; j++ {
			filename := "./" + prefix + "-" + strconv.Itoa(j) + "." + ext
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				i.FileName = filename
				return filename, nil
			}
		}

		return "", errors.New("Failed to find free file for db")
	}

	return i.FileName, nil
}

func (i *Store) getDb() *sql.DB {

	if i.handle == nil {

		if i.prefix != "" && i.ext != "" {
			fileName, err := i.getDbFile(i.prefix, i.ext)
			if err != nil {
				panic(err)
			}

			i.FileName = fileName
		}

		handle, err := sql.Open("sqlite3", i.FileName)
		if err != nil {
			panic(err)
		}

		i.handle = handle
	}

	return i.handle
}

func (i *Store) prepareDb() {
	db := i.getDb()

	createTableStmt := `
		CREATE TABLE IF NOT EXISTS files (
			file		INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			book 		VARCHAR(49) NOT NULL,
			filename	VARCHAR(46) NOT NULL,

			CONSTRAINT file_constraint UNIQUE (book, filename)
		);

		CREATE TABLE IF NOT EXISTS keywords ( 
			keyword 	TEXT NOT NULL COLLATE NOCASE, 
			file		INTEGER NOT NULL,
			
			FOREIGN KEY(file) REFERENCES files(file)

			-- CONSTRAINT keywords_constraint UNIQUE(keyword, file)
		);
		
		CREATE INDEX IF NOT EXISTS keywords_index ON keywords(keyword);
		`

	_, err := db.Exec(createTableStmt)

	if err != nil {
		panic(err)
	}
}

func (i *Store) findID(tx *sql.Tx, bookID, file string) (id int64) {

	stmt, _ := tx.Prepare("SELECT file FROM files WHERE book = '?' AND filename = '?'")

	rows, err := tx.Stmt(stmt).Query(bookID, file)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)

		if err != nil {
			panic(err)
		}

		break
	}

	log.Println("Fetched existing id", id)

	return
}

func (i *Store) insertIndexes() chan *types.Indexes {

	ch := make(chan *types.Indexes)

	go func() {

		db := i.getDb()

		bookStmt, err := db.Prepare("INSERT INTO files(book, filename) values(?, ?)")
		if err != nil {
			panic(err)
		}

		keywordStmt, err := db.Prepare("INSERT INTO keywords(keyword, file) values(?, ?)")
		if err != nil {
			panic(err)
		}

		idMap := make(map[string]int64)

		for {
			indexes := <-ch

			tx, err := db.Begin()
			if err != nil {
				panic(err)
			}

			for _, index := range indexes.Keywords {

				// Check if we have this id already
				id := idMap[indexes.BookID+index.File]
				if id == 0 {
					// Try to insert, will fail if exists due to unique constraint
					res, err := tx.Stmt(bookStmt).Exec(indexes.BookID, index.File)
					if err != nil {
						// unique failed, time to search
						id = i.findID(tx, indexes.BookID, index.File)
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
					//tx.Rollback()
					log.Fatalln("Panic for", indexes.BookID, id)
					log.Fatalf("%+v\n", err)
					//panic(err)
				}
			}

			tx.Commit()
		}
	}()

	return ch
}

// FetchIndex downloads the index file, decodes it and requests it to be added to the store
func (i *Store) FetchIndex(book *types.Book) {

	parts := map[string]string{"Id": book.ID, "Language": book.Language, "Version": book.Version}
	var url bytes.Buffer
	i.f1FragmentTemplate.Execute(&url, parts)

	resp, err := http.Get(url.String())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	indexes := new(types.Indexes)
	indexes.Etag = resp.Header.Get("Etag")

	err = decoder.Decode(indexes)
	if err != nil {
		log.Panic("Failed to parse", url)
		panic(err)
	}

	if i.indexWriteChan == nil {
		i.indexWriteChan = i.insertIndexes()
	}

	i.indexWriteChan <- indexes
}

// KeywordResult is a single result, pointing to a BookID and a Filename inside this book which has the keyword
type KeywordResult struct {
	BookID   string
	Filename string
}

// LookupKeyword looks up the keyword in the store, and returns a slice of KeywordResult
func (i *Store) LookupKeyword(keyword string) []KeywordResult {
	stmt, err := i.handle.Prepare(`
		SELECT files.book, files.filename
		FROM keywords
		INNER JOIN files
		ON keywords.file = files.file
		WHERE keywords.keyword = ? COLLATE NOCASE`)

	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query(keyword)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var results []KeywordResult

	for rows.Next() {
		result := KeywordResult{}
		rows.Scan(&result.BookID, &result.Filename)

		results = append(results, result)
	}

	return results
}

// DatabaseStatistics contains some small statistic numbers from the amount of content in the system
type DatabaseStatistics struct {
	NumberOfFiles int `json:"number_of_files"`
	KeywordCount  int `json:"number_of_keywords"`
}

// GetStatistics reads the statistics from the db storage and returns them
func (i *Store) GetStatistics() DatabaseStatistics {
	fileCount := i.getFileCount()
	keywordCount := i.getKeywordCount()

	return DatabaseStatistics{NumberOfFiles: fileCount, KeywordCount: keywordCount}
}

func (i *Store) getFileCount() (count int) {
	rows, err := i.handle.Query(`
		SELECT COUNT() 
		as count 
		FROM files`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		} else {
			return
		}
	}

	return
}

func (i *Store) getKeywordCount() (count int) {
	rows, err := i.handle.Query(`
		SELECT COUNT() 
		as count 
		FROM keywords`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		} else {
			return
		}
	}

	return
}

func (i *Store) getKeywords() (keywords []string) {
	rows, err := i.handle.Query(`
		SELECT keyword
		FROM keywords
		ORDER BY keyword`)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var keyword string

		err := rows.Scan(&keyword)
		if err != nil {
			panic(err)
		}

		keywords = append(keywords, keyword)
	}

	return
}
