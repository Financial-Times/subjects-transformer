package main

import (
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type subjectService interface {
	getSubjects() ([]SubjectLink, bool)
	getSubjectByUuid(uuid string) (Subject, bool)
}

type subjectServiceImpl struct {
	repository   Repository
	transformer  subjectTransformer
	baseUrl      string
	subjectsMap  map[string]Subject
	subjectLinks []SubjectLink
}

func newSubjectService(repo Repository, transformer subjectTransformer, baseUrl string) (subjectService, error) {

	s := &subjectServiceImpl{repository: repo, transformer: transformer, baseUrl: baseUrl}
	err := s.init()
	if err != nil {
		return &subjectServiceImpl{}, err
	}
	return s, nil
}

func (s *subjectServiceImpl) init() error {
	s.subjectsMap = make(map[string]Subject)
	tax, err := s.repository.getSubjectsTaxonomy()
	if err != nil {
		return err
	}
	s.initSubjectsMap(tax.Terms)
	return nil
}

func (s *subjectServiceImpl) getSubjects() ([]SubjectLink, bool) {
	if len(s.subjectLinks) > 0 {
		return s.subjectLinks, true
	}
	return s.subjectLinks, false
}

func (s *subjectServiceImpl) getSubjectByUuid(uuid string) (Subject, bool) {
	subject, found := s.subjectsMap[uuid]
	return subject, found
}

func (s *subjectServiceImpl) initSubjectsMap(terms []Term) {
	for _, t := range terms {
		sub := s.transformer.transform(t)
		s.subjectsMap[sub.UUID] = sub
		s.subjectLinks = append(s.subjectLinks, SubjectLink{ApiUrl: s.baseUrl + sub.UUID})
		s.initSubjectsMap(t.Children.Terms)
	}
}
