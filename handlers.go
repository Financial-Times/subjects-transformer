package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

type subjectsHandler struct {
	service subjectService
}

func newSubjectsHandler(service subjectService) subjectsHandler {
	return subjectsHandler{service: service}
}

func (h *subjectsHandler) getSubjects(writer http.ResponseWriter, req *http.Request) {
	obj, found := h.service.getSubjects()
	writeJsonResponse(obj, found, writer)
}

func (h *subjectsHandler) getSubjectByUuid(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	obj, found := h.service.getSubjectByUuid(uuid)
	writeJsonResponse(obj, found, writer)
}

func writeJsonResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJsonError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJsonError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}
