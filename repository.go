package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type repository interface {
	getSubjectsTaxonomy() (taxonomy, error)
}

type tmeRepository struct {
	httpClient              httpClient
	principalHeader         string
	structureServiceBaseUrl string
}

func newTmeRepository(client httpClient, structureServiceBaseUrl string, principalHeader string) repository {
	return &tmeRepository{httpClient: client, principalHeader: principalHeader, structureServiceBaseUrl: structureServiceBaseUrl}
}

func (t *tmeRepository) getSubjectsTaxonomy() (taxonomy, error) {
	req, err := http.NewRequest("GET", t.structureServiceBaseUrl+"/metadata-services/structure/1.0/taxonomies/subjects/terms?includeDisabledTerms=true", nil)
	if err != nil {
		return taxonomy{}, err
	}
	req.Header.Set("ClientUserPrincipal", t.principalHeader)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return taxonomy{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return taxonomy{}, errors.New(fmt.Sprintf("Structure service returned %d", resp.StatusCode))
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return taxonomy{}, err
	}

	tax := taxonomy{}
	err = xml.Unmarshal(contents, &tax)
	if err != nil {
		return taxonomy{}, err
	}
	return tax, nil
}
