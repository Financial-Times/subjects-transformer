package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubjects(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		baseURL  string
		terms    []term
		subjects []subjectLink
		found    bool
		err      error
	}{
		{"Success", "localhost:8080/transformers/subjects/",
			[]term{term{CanonicalName: "SubjectZ", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "mNGQ2MWQ0NDMtMDc5Mi00NWExLTlkMGQtNWZhZjk0NGExOWU2-Z2VucVz"}},
			[]subjectLink{subjectLink{APIURL: "localhost:8080/transformers/subjects/83c29673-27a5-3d63-b801-750f024a1ef2"},
				subjectLink{APIURL: "localhost:8080/transformers/subjects/d75cbca8-aff7-3ac9-8171-604ee8ad6daa"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/subjects/", []term{}, []subjectLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSubjectService(&repo, test.baseURL, "Subjects", 10000)
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
		terms   []term
		uuid    string
		subject subject
		found   bool
		err     error
	}{
		{"Success", []term{term{CanonicalName: "SubjectZ", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "TkdRMk1XUTBORE10TURjNU1pMDBOV0V4TFRsa01HUXROV1poWmprME5HRXhPV1UyLVoyVnVjbVZ6-U2VjdGlvbnM=]"}},
			"83c29673-27a5-3d63-b801-750f024a1ef2", getDummySubject("83c29673-27a5-3d63-b801-750f024a1ef2", "SubjectZ", "YjgzMzc1NTktYWMwOC0zNDA0LTkwMjUtYmFkNTFlYmUyZmM3-U3ViamVjdHM="), true, nil},
		{"Not found", []term{term{CanonicalName: "SubjectZ", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Feature", RawID: "NGQ2MWdefsdfsfcmVz"}},
			"some uuid", subject{}, false, nil},
		{"Error on init", []term{}, "some uuid", subject{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSubjectService(&repo, "", "Subjects", 10000)
		expectedSubject, found := service.getSubjectByUUID(test.uuid)
		assert.Equal(test.subject, expectedSubject, fmt.Sprintf("%s: Expected subject incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	terms []term
	err   error
}

func (d *dummyRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	if startRecord > 0 {
		return nil, d.err
	}
	var interfaces = make([]interface{}, len(d.terms))
	for i, data := range d.terms {
		interfaces[i] = data
	}
	return interfaces, d.err
}
func (d *dummyRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return d.terms[0], d.err
}

func getDummySubject(uuid string, prefLabel string, tmeID string) subject {
	return subject{
		UUID:                   uuid,
		PrefLabel:              prefLabel,
		PrimaryType:            primaryType,
		TypeHierarchy:          subjectTypes,
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{tmeID}, Uuids: []string{uuid}}}
}
