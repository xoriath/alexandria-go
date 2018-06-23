package types

import (
	"fmt"
)

// Book is a single deployed publication
type Book struct {
	ID       string `xml:"id,attr"`
	Language string `xml:"lang,attr"`
	Version  string `xml:"version,attr"`
	Title    string `xml:"title,attr"`

	Webhelp bool `xml:"webhelp,attr"`

	Timestamp string `xml:"timestamp,attr"`

	CabMD5  string `xml:"cab.md5,attr"`
	CabSHA1 string `xml:"cab.sha1,attr"`

	MshcMD5  string `xml:"mshc.md5,attr"`
	MshcSHA1 string `xml:"mshc.sha1,attr"`

	CompressedSize int `xml:"size.compressed,attr"`
	RawSize        int `xml:"size.raw,attr"`

	Products   []Product `xml:"products>product"`
	Parameters []Param   `xml:"params>param"`
}

func (b *Book) InProduct(name string) bool {
	return b.Product(name) != nil
}

func (b *Book) Product(name string) *Product {
	for _, product := range b.Products {
		if product.Name == name {
			return &product
		}
	}

	return nil
}

func (b *Book) Description() string {
	val, _ := b.findParam("FDESCRIPTION", "logical")
	return val
}

func (b *Book) PublicationType() string {
	val, _ := b.findParam("FISHPUBLICATIONTYPE", "logical")
	return val
}

func (b *Book) findParam(name, level string) (string, error) {
	for _, param := range b.Parameters {
		if param.Name == name && param.Level == level {
			return param.Value, nil
		}
	}

	return "", fmt.Errorf("Couldn't find parameter %s at %s level", name, level)
}
