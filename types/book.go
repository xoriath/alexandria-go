package types

// Book is a single deployed publication
type Book struct {
	ID string `xml:"id,attr"`
	Language string `xml:"lang,attr"`
	Version string `xml:"version,attr"`
	Title string `xml:"title,attr"`

	Webhelp bool `xml:"webhelp,attr"`

	Timestamp string `xml:"timestamp,attr"`

	CabMD5 string `xml:"cab.md5,attr"`
	CabSHA1 string `xml:"cab.sha1,attr"`

	MshcMD5 string `xml:"mshc.md5,attr"`
	MshcSHA1 string `xml:"mshc.sha1,attr"`

	CompressedSize int `xml:"size.compressed,attr"`
	RawSize int `xml:"size.raw,attr"`

	Products []Product `xml:"products>product"`
	Parameters []Param `xml:"params>param"`
}
