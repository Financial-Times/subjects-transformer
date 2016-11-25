package main

type subject struct {
	UUID                   string                 `json:"uuid"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	PrefLabel              string                 `json:"prefLabel"`
	PrimaryType            string                 `json:"type"`
	TypeHierarchy          []string               `json:"types"`
}

type alternativeIdentifiers struct {
	TME   []string `json:"TME,omitempty"`
	Uuids []string `json:"uuids,omitempty"`
}

type subjectLink struct {
	APIURL string `json:"apiUrl"`
}

var primaryType = "Subject"
var subjectTypes = []string{"Thing", "Concept", "Classification", "Subject"}
