package types

// Index is a single keyword mapping to a single file for a given book
type Index struct {
	File    string `xml:"file,attr"`
	Keyword string `xml:"f1,attr"`
}
