package main

import (
	"net/http"
)

type Client interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type SubjectService interface {
	GetSubjects() ([]SubjectLink, bool)
	GetSubjectByUuid(uuid string) (Subject, bool)
}

type subjectServiceImpl struct {
	repository   Repository
	transformer  SubjectTransformer
	baseUrl      string
	subjectsMap  map[string]Subject
	subjectLinks []SubjectLink
}

func NewSubjectService(repo Repository, transformer SubjectTransformer, baseUrl string) (SubjectService, error) {

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

func (s *subjectServiceImpl) GetSubjects() ([]SubjectLink, bool) {
	if len(s.subjectLinks) > 0 {
		return s.subjectLinks, true
	}
	return s.subjectLinks, false
}

func (s *subjectServiceImpl) GetSubjectByUuid(uuid string) (Subject, bool) {
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
