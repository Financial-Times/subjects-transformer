package main

import (
	"flag"
	"fmt"
	digest "github.com/FeNoMeNa/goha"
	"github.com/Financial-Times/http-handlers-go"
	"github.com/Financial-Times/subjects-transformer/handlers"
	"github.com/Financial-Times/subjects-transformer/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"time"
)

var port = flag.Int("port", 8080, "Port to listen on")
var username = flag.String("structureServiceUsername", "", "Structure service username used for http digest authentication")
var password = flag.String("structureServicePassword", "", "Structure service password used for http digest authentication")
var principalHeader = flag.String("structureServicePrincipalHeader", "", "Structure service principal header used for authentication")
var baseUrl = flag.String("baseUrl", "http://localhost:8080/transformers/subjects/", "Base url")
var structureServiceBaseUrl = flag.String("structureServiceBaseUrl", "http://metadata.internal.ft.com:83", "Structure service base url")

func init() {
	log.SetFormatter(new(log.JSONFormatter))
}

func main() {
	flag.Parse()

	c := digest.NewClient(*username, *password)
	c.Timeout(10 * time.Second)
	s, err := service.NewSubjectService(service.NewTmeRepository(c, *structureServiceBaseUrl, *principalHeader), service.SubjectTransformer{}, *baseUrl)
	if err != nil {
		log.Errorf("Error while creating SubjectsService: [%v]", err.Error())
	}
	h := handlers.NewSubjectsHandler(s)
	m := mux.NewRouter()
	m.HandleFunc("/transformers/subjects", h.GetSubjects).Methods("GET")
	m.HandleFunc("/transformers/subjects/{uuid}", h.GetSubjectByUuid).Methods("GET")
	http.Handle("/", m)

	log.Printf("listening on %d", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port),
		httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry,
			httphandlers.TransactionAwareRequestLoggingHandler(log.StandardLogger(), m)))
}
