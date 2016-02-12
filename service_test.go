package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSubjects(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		baseURL  string
		tax      taxonomy
		subjects []subjectLink
		found    bool
		err      error
	}{
		{"Success", "localhost:8080/transformers/subjects/",
			taxonomy{Terms: []term{term{CanonicalName: "Company News", ID: "MQ==-U3ViamVjdHM=", Children: children{[]term{term{CanonicalName: "Bankruptcy & Receivership", ID: "Mg==-U3ViamVjdHM="}}}}}},
			[]subjectLink{subjectLink{APIURL: "localhost:8080/transformers/subjects/29b56d8f-3528-37ae-9551-c50a0d37d4bb"},
				subjectLink{APIURL: "localhost:8080/transformers/subjects/6725e13a-276d-3096-91fe-bf7db924ff03"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/subjects/", taxonomy{}, []subjectLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newSubjectService(&repo, test.baseURL)
		expectedSubjects, found := service.getSubjects()
		assert.Equal(test.subjects, expectedSubjects, fmt.Sprintf("%s: Expected subjects link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetSubjectByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		tax     taxonomy
		uuid    string
		subject subject
		found   bool
		err     error
	}{
		{"Success", taxonomy{Terms: []term{term{CanonicalName: "Company News", ID: "MQ==-U3ViamVjdHM=", Children: children{[]term{term{CanonicalName: "Bankruptcy & Receivership", ID: "Mg==-U3ViamVjdHM="}}}}}},
			"29b56d8f-3528-37ae-9551-c50a0d37d4bb", subject{UUID: "29b56d8f-3528-37ae-9551-c50a0d37d4bb", CanonicalName: "Company News", TmeIdentifier: "MQ==-U3ViamVjdHM=", Type: "Subject"}, true, nil},
		{"Not found", taxonomy{Terms: []term{term{CanonicalName: "Company News", ID: "MQ==-U3ViamVjdHM=", Children: children{[]term{term{CanonicalName: "Bankruptcy & Receivership", ID: "Mg==-U3ViamVjdHM="}}}}}},
			"some uuid", subject{}, false, nil},
		{"Error on init", taxonomy{}, "some uuid", subject{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newSubjectService(&repo, "")
		expectedSubject, found := service.getSubjectByUUID(test.uuid)
		assert.Equal(test.subject, expectedSubject, fmt.Sprintf("%s: Expected subject incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	tax taxonomy
	err error
}

func (d *dummyRepo) getSubjectsTaxonomy() (taxonomy, error) {
	return d.tax, d.err
}
