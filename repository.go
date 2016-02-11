package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Repository interface {
	getSubjectsTaxonomy() (taxonomy, error)
}

type TmeRepository struct {
	httpClient              httpClient
	principalHeader         string
	structureServiceBaseUrl string
}

func NewTmeRepository(client httpClient, structureServiceBaseUrl string, principalHeader string) Repository {
	return &TmeRepository{httpClient: client, principalHeader: principalHeader, structureServiceBaseUrl: structureServiceBaseUrl}
}

func (t *TmeRepository) getSubjectsTaxonomy() (taxonomy, error) {
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
