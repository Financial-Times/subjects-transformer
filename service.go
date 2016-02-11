package main

import (
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type subjectService interface {
	getSubjects() ([]subjectLink, bool)
	getSubjectByUUID(uuid string) (subject, bool)
}

type subjectServiceImpl struct {
	repository   repository
	transformer  subjectTransformer
	baseURL      string
	subjectsMap  map[string]subject
	subjectLinks []subjectLink
}

func newSubjectService(repo repository, transformer subjectTransformer, baseURL string) (subjectService, error) {

	s := &subjectServiceImpl{repository: repo, transformer: transformer, baseURL: baseURL}
	err := s.init()
	if err != nil {
		return &subjectServiceImpl{}, err
	}
	return s, nil
}

func (s *subjectServiceImpl) init() error {
	s.subjectsMap = make(map[string]subject)
	tax, err := s.repository.getSubjectsTaxonomy()
	if err != nil {
		return err
	}
	s.initSubjectsMap(tax.Terms)
	return nil
}

func (s *subjectServiceImpl) getSubjects() ([]subjectLink, bool) {
	if len(s.subjectLinks) > 0 {
		return s.subjectLinks, true
	}
	return s.subjectLinks, false
}

func (s *subjectServiceImpl) getSubjectByUUID(uuid string) (subject, bool) {
	subject, found := s.subjectsMap[uuid]
	return subject, found
}

func (s *subjectServiceImpl) initSubjectsMap(terms []term) {
	for _, t := range terms {
		sub := s.transformer.transform(t)
		s.subjectsMap[sub.UUID] = sub
		s.subjectLinks = append(s.subjectLinks, subjectLink{APIURL: s.baseURL + sub.UUID})
		s.initSubjectsMap(t.Children.Terms)
	}
}
