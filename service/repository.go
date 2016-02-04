package service

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Financial-Times/subjects-transformer/model"
	"io/ioutil"
	"net/http"
)

type Repository interface {
	getSubjectsTaxonomy() (model.Taxonomy, error)
}

type TmeRepository struct {
	httpClient              Client
	principalHeader         string
	structureServiceBaseUrl string
}

func NewTmeRepository(client Client, structureServiceBaseUrl string, principalHeader string) Repository {
	return &TmeRepository{httpClient: client, principalHeader: principalHeader, structureServiceBaseUrl: structureServiceBaseUrl}
}

func (t *TmeRepository) getSubjectsTaxonomy() (model.Taxonomy, error) {
	req, err := http.NewRequest("GET", t.structureServiceBaseUrl+"/metadata-services/structure/1.0/taxonomies/subjects/terms?includeDisabledTerms=true", nil)
	if err != nil {
		return model.Taxonomy{}, err
	}
	req.Header.Set("Clientuserprincipal", t.principalHeader)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return model.Taxonomy{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Taxonomy{}, errors.New(fmt.Sprintf("Structure service returned %d", resp.StatusCode))
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Taxonomy{}, err
	}

	tax := model.Taxonomy{}
	err = xml.Unmarshal(contents, &tax)
	if err != nil {
		return model.Taxonomy{}, err
	}
	return tax, nil
}
