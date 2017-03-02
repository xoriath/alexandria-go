package types

// Param is a signgle parameter
type Param struct {
	Level string `xml:"level,attr"`
	Name string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}