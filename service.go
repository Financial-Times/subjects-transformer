package main

import (
	"fmt"
	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type subjectService interface {
	getSubjects() ([]subjectLink, bool)
	getSubjectByUUID(uuid string) (subject, bool)
	checkConnectivity() error
	getSubjectCount() int
	getSubjectIds() []string
	reload() error
}

type subjectServiceImpl struct {
	repository    tmereader.Repository
	baseURL       string
	subjectsMap   map[string]subject
	subjectLinks  []subjectLink
	taxonomyName  string
	maxTmeRecords int
}

func newSubjectService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (subjectService, error) {
	s := &subjectServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.reload()
	if err != nil {
		return &subjectServiceImpl{}, err
	}
	return s, nil
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

func (s *subjectServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back? Maybe a healthcheck or gtg endpoint?
	// TODO: Can we use a count from our responses while actually in use to trigger a healthcheck?
	//	_, err := s.repository.GetTmeTermsFromIndex(1)
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func (s *subjectServiceImpl) initSubjectsMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		top := transformSubject(t, s.taxonomyName)
		s.subjectsMap[top.UUID] = top
		s.subjectLinks = append(s.subjectLinks, subjectLink{APIURL: s.baseURL + top.UUID})
	}
}

func (s *subjectServiceImpl) getSubjectCount() int {
	return len(s.subjectLinks)
}

func (s *subjectServiceImpl) getSubjectIds() []string {
	i := 0
	keys := make([]string, len(s.subjectsMap))

	for k := range s.subjectsMap {
		keys[i] = k
		i++
	}
	return keys
}

func (s *subjectServiceImpl) reload() error {
	s.subjectsMap = make(map[string]subject)
	responseCount := 0
	log.Println("Fetching subjects from TME")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Println("Finished fetching subjects from TME")
			break
		}
		s.initSubjectsMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d subjects links\n", len(s.subjectLinks))

	return nil
}
