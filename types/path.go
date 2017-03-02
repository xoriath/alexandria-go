package types

// Path is a placement in a tree. Used for displaying.
type Path struct {
	Path string `xml:"path,attr"`
	Priority int `xml:"priority,attr"`
}