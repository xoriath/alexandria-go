package types

// Indexes is the set of keyword to file mappings for a given book
type Indexes struct {
	BookID   string  `xml:"id,attr"`
	Keywords []Index `xml:"index"`
}
