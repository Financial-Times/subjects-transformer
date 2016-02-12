package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestGetSubjectsTaxonomy(t *testing.T) {
	assert := assert.New(t)
	subjectsXML, _ := os.Open("sample_subjects.xml")
	tests := []struct {
		name string
		repo repository
		tax  taxonomy
		err  error
	}{
		{"Success", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(subjectsXML)}}),
			taxonomy{Terms: []term{term{CanonicalName: "Company News", ID: "MQ==-U3ViamVjdHM=", Children: children{[]term{term{CanonicalName: "Bankruptcy & Receivership", ID: "Mg==-U3ViamVjdHM="}}}}}}, nil},
		{"Error", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(subjectsXML)}, err: errors.New("Some error")}),
			taxonomy{}, errors.New("Some error")},
		{"Non 200 from structure service", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(subjectsXML)}}),
			taxonomy{}, errors.New("Structure service returned 400")},
		{"Unmarshalling error", repo(dummyClient{assert: assert, structureServiceBaseURL: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte("Non xml")))}}),
			taxonomy{}, errors.New("EOF")},
	}

	for _, test := range tests {
		expectedTax, err := test.repo.getSubjectsTaxonomy()
		assert.Equal(test.tax, expectedTax, fmt.Sprintf("%s: Expected taxonomy incorrect", test.name))
		assert.Equal(test.err, err)
	}

}

func repo(c dummyClient) repository {
	return &tmeRepository{httpClient: &c, principalHeader: c.principalHeader, structureServiceBaseURL: c.structureServiceBaseURL}
}

type dummyClient struct {
	assert                  *assert.Assertions
	resp                    http.Response
	err                     error
	principalHeader         string
	structureServiceBaseURL string
}

func (d *dummyClient) Do(req *http.Request) (resp *http.Response, err error) {
	d.assert.Equal(d.principalHeader, req.Header.Get("ClientUserPrincipal"), fmt.Sprintf("Expected ClientUserPrincipal header incorrect"))
	d.assert.Equal(d.structureServiceBaseURL+"/metadata-services/structure/1.0/taxonomies/subjects/terms?includeDisabledTerms=true", req.URL.String(), fmt.Sprintf("Expected url incorrect"))
	return &d.resp, d.err
}