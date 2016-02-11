package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const uuid = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getSubjectsResponse = "[{\"apiUrl\":\"http://localhost:8080/transformers/subjects/bba39990-c78d-3629-ae83-808c333c6dbc\"}]\n"
const getSubjectByUuidResponse = "{\"uuid\":\"bba39990-c78d-3629-ae83-808c333c6dbc\",\"canonicalName\":\"Metals Markets\",\"tmeIdentifier\":\"MTE3-U3ViamVjdHM=\",\"type\":\"Subject\"}\n"

func TestHandlers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		req          *http.Request
		dummyService subjectService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get subject by uuid", newRequest("GET", fmt.Sprintf("/transformers/subjects/%s", uuid)), &dummyService{found: true, subjects: []Subject{Subject{UUID: uuid, CanonicalName: "Metals Markets", TmeIdentifier: "MTE3-U3ViamVjdHM=", Type: "Subject"}}}, http.StatusOK, "application/json", getSubjectByUuidResponse},
		{"Not found - get subject by uuid", newRequest("GET", fmt.Sprintf("/transformers/subjects/%s", uuid)), &dummyService{found: false, subjects: []Subject{Subject{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get subjects", newRequest("GET", "/transformers/subjects"), &dummyService{found: true, subjects: []Subject{Subject{UUID: uuid}}}, http.StatusOK, "application/json", getSubjectsResponse},
		{"Not found - get subjects", newRequest("GET", "/transformers/subjects"), &dummyService{found: false, subjects: []Subject{}}, http.StatusNotFound, "application/json", ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(test.body, rec.Body.String(), fmt.Sprintf("%s: Wrong body", test.name))
	}
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func router(s subjectService) *mux.Router {
	m := mux.NewRouter()
	h := newSubjectsHandler(s)
	m.HandleFunc("/transformers/subjects", h.getSubjects).Methods("GET")
	m.HandleFunc("/transformers/subjects/{uuid}", h.getSubjectByUuid).Methods("GET")
	return m
}

type dummyService struct {
	found    bool
	subjects []Subject
}

func (s *dummyService) getSubjects() ([]SubjectLink, bool) {
	var subjectLinks []SubjectLink
	for _, sub := range s.subjects {
		subjectLinks = append(subjectLinks, SubjectLink{ApiUrl: "http://localhost:8080/transformers/subjects/" + sub.UUID})
	}
	return subjectLinks, s.found
}

func (s *dummyService) getSubjectByUuid(uuid string) (Subject, bool) {
	return s.subjects[0], s.found
}
