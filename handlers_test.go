package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getSubjectsResponse = `[{"apiUrl":"http://localhost:8080/transformers/subjects/bba39990-c78d-3629-ae83-808c333c6dbc"}]`
const getSubjectByUUIDResponse = `{"uuid":"bba39990-c78d-3629-ae83-808c333c6dbc","alternativeIdentifiers":{"TME":["MTE3-U3ViamVjdHM="],"uuids":["bba39990-c78d-3629-ae83-808c333c6dbc"]},"prefLabel":"Global Subjects","type":"Subject"}`

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
		{"Success - get subject by uuid", newRequest("GET", fmt.Sprintf("/transformers/subjects/%s", testUUID)), &dummyService{found: true, subjects: []subject{getDummySubject(testUUID, "Global Subjects", "MTE3-U3ViamVjdHM=")}}, http.StatusOK, "application/json", getSubjectByUUIDResponse},
		{"Not found - get subject by uuid", newRequest("GET", fmt.Sprintf("/transformers/subjects/%s", testUUID)), &dummyService{found: false, subjects: []subject{subject{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get subjects", newRequest("GET", "/transformers/subjects"), &dummyService{found: true, subjects: []subject{subject{UUID: testUUID}}}, http.StatusOK, "application/json", getSubjectsResponse},
		{"Not found - get subjects", newRequest("GET", "/transformers/subjects"), &dummyService{found: false, subjects: []subject{}}, http.StatusNotFound, "application/json", ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(strings.TrimSpace(test.body), strings.TrimSpace(rec.Body.String()), fmt.Sprintf("%s: Wrong body", test.name))
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
	m.HandleFunc("/transformers/subjects/{uuid}", h.getSubjectByUUID).Methods("GET")
	return m
}

type dummyService struct {
	found    bool
	subjects []subject
}

func (s *dummyService) getSubjects() ([]subjectLink, bool) {
	var subjectLinks []subjectLink
	for _, sub := range s.subjects {
		subjectLinks = append(subjectLinks, subjectLink{APIURL: "http://localhost:8080/transformers/subjects/" + sub.UUID})
	}
	return subjectLinks, s.found
}

func (s *dummyService) getSubjectByUUID(uuid string) (subject, bool) {
	return s.subjects[0], s.found
}

func (s *dummyService) checkConnectivity() error {
	return nil
}
