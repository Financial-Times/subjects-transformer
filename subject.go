package main

type subject struct {
	UUID          string `json:"uuid"`
	CanonicalName string `json:"canonicalName"`
	TmeIdentifier string `json:"tmeIdentifier,omitempty"`
	Type          string `json:"type"`
}

type subjectLink struct {
	ApiUrl string `json:"apiUrl"`
}
