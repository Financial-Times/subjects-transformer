package model

type Taxonomy struct {
	Terms []Term `xml:"term"`
}

type Term struct {
	CanonicalName string `xml:"canonicalName"`
	Id string `xml:"id,attr"`
	Children Children `xml:"children"`
}

type Children struct {
	Terms []Term `xml:"term"`
}