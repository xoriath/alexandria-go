package types

// Product is a target product for a book
type Product struct {
	Name string `xml:"name,attr"`
	SkuID int `xml:"sku.id,attr"`
	SkuName string `xml:"sku.name,attr"`
	Paths []Path `xml:"paths>path"`
}