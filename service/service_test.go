package service

import (
	"errors"
	"fmt"
	"github.com/Financial-Times/subjects-transformer/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSubjects(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		baseUrl  string
		tax      model.Taxonomy
		subjects []model.SubjectLink
		found    bool
		err      error
	}{
		{"Success", "localhost:8080/transformers/subjects/",
			model.Taxonomy{Terms: []model.Term{model.Term{CanonicalName: "Company News", Id: "MQ==-U3ViamVjdHM=", Children: model.Children{[]model.Term{model.Term{CanonicalName: "Bankruptcy & Receivership", Id: "Mg==-U3ViamVjdHM="}}}}}},
			[]model.SubjectLink{model.SubjectLink{ApiUrl: "localhost:8080/transformers/subjects/29b56d8f-3528-37ae-9551-c50a0d37d4bb"},
				model.SubjectLink{ApiUrl: "localhost:8080/transformers/subjects/6725e13a-276d-3096-91fe-bf7db924ff03"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/subjects/", model.Taxonomy{}, []model.SubjectLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := NewSubjectService(&repo, SubjectTransformer{}, test.baseUrl)
		expectedSubjects, found := service.GetSubjects()
		assert.Equal(test.subjects, expectedSubjects, fmt.Sprintf("%s: Expected subjects link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetSubjectByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		tax     model.Taxonomy
		uuid    string
		subject model.Subject
		found   bool
		err     error
	}{
		{"Success", model.Taxonomy{Terms: []model.Term{model.Term{CanonicalName: "Company News", Id: "MQ==-U3ViamVjdHM=", Children: model.Children{[]model.Term{model.Term{CanonicalName: "Bankruptcy & Receivership", Id: "Mg==-U3ViamVjdHM="}}}}}},
			"29b56d8f-3528-37ae-9551-c50a0d37d4bb", model.Subject{UUID: "29b56d8f-3528-37ae-9551-c50a0d37d4bb", CanonicalName: "Company News", TmeIdentifier: "MQ==-U3ViamVjdHM=", Type: "Subject"}, true, nil},
		{"Not found", model.Taxonomy{Terms: []model.Term{model.Term{CanonicalName: "Company News", Id: "MQ==-U3ViamVjdHM=", Children: model.Children{[]model.Term{model.Term{CanonicalName: "Bankruptcy & Receivership", Id: "Mg==-U3ViamVjdHM="}}}}}},
			"some uuid", model.Subject{}, false, nil},
		{"Error on init", model.Taxonomy{}, "some uuid", model.Subject{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := NewSubjectService(&repo, SubjectTransformer{}, "")
		expectedSubject, found := service.GetSubjectByUuid(test.uuid)
		assert.Equal(test.subject, expectedSubject, fmt.Sprintf("%s: Expected subject incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	tax model.Taxonomy
	err error
}

func (d *dummyRepo) getSubjectsTaxonomy() (model.Taxonomy, error) {
	return d.tax, d.err
}
