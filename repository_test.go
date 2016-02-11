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
	subjectsXml, _ := os.Open("sample_subjects.xml")
	tests := []struct {
		name string
		repo Repository
		tax  taxonomy
		err  error
	}{
		{"Success", repo(dummyClient{assert: assert, structureServiceBaseUrl: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(subjectsXml)}}),
			taxonomy{Terms: []Term{Term{CanonicalName: "Company News", Id: "MQ==-U3ViamVjdHM=", Children: Children{[]Term{Term{CanonicalName: "Bankruptcy & Receivership", Id: "Mg==-U3ViamVjdHM="}}}}}}, nil},
		{"Error", repo(dummyClient{assert: assert, structureServiceBaseUrl: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(subjectsXml)}, err: errors.New("Some error")}),
			taxonomy{}, errors.New("Some error")},
		{"Non 200 from structure service", repo(dummyClient{assert: assert, structureServiceBaseUrl: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(subjectsXml)}}),
			taxonomy{}, errors.New("Structure service returned 400")},
		{"Unmarshalling error", repo(dummyClient{assert: assert, structureServiceBaseUrl: "http://metadata.internal.ft.com:83", principalHeader: "someHeader",
			resp: http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte("Non xml")))}}),
			taxonomy{}, errors.New("EOF")},
	}

	for _, test := range tests {
		expectedTax, err := test.repo.getSubjectsTaxonomy()
		assert.Equal(test.tax, expectedTax, fmt.Sprintf("%s: Expected taxonomy incorrect", test.name))
		assert.Equal(test.err, err)
	}

}

func repo(c dummyClient) Repository {
	return &TmeRepository{httpClient: &c, principalHeader: c.principalHeader, structureServiceBaseUrl: c.structureServiceBaseUrl}
}

type dummyClient struct {
	assert                  *assert.Assertions
	resp                    http.Response
	err                     error
	principalHeader         string
	structureServiceBaseUrl string
}

func (d *dummyClient) Do(req *http.Request) (resp *http.Response, err error) {
	d.assert.Equal(d.principalHeader, req.Header.Get("ClientUserPrincipal"), fmt.Sprintf("Expected ClientUserPrincipal header incorrect"))
	d.assert.Equal(d.structureServiceBaseUrl+"/metadata-services/structure/1.0/taxonomies/subjects/terms?includeDisabledTerms=true", req.URL.String(), fmt.Sprintf("Expected url incorrect"))
	return &d.resp, d.err
}
