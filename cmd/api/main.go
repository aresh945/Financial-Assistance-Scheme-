package main

import (
	"financial_assistance/internal/handler"
	"financial_assistance/internal/repository/postgres"
	"financial_assistance/internal/service"
	"financial_assistance/pkg/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	dbConfig := &database.Config{
		Host:     "localhost",
		Port:     5433,
		User:     "postgres",
		Password: "aresh",
		DBName:   "Tutorial1",
		SSLMode:  "disable",
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Could not initialize database connection: %v", err)
	}
	defer db.Close()

	applicantRepo := postgres.NewApplicantRepo(db)
	schemeRepo := postgres.NewSchemeRepo(db)
	applicationRepo := postgres.NewApplicationRepo(db)

	svc := service.NewService(applicantRepo, schemeRepo, applicationRepo)

	h := handler.NewHandler(svc)

	r := mux.NewRouter()

	r.HandleFunc("/api/applicants", h.GetAllApplicants).Methods("GET")
	r.HandleFunc("/api/applicants", h.CreateApplicant).Methods("POST")
	r.HandleFunc("/api/schemes", h.GetAllSchemes).Methods("GET")
	r.HandleFunc("/api/schemes", h.CreateScheme).Methods("POST")
	r.HandleFunc("/api/schemes/eligible", h.GetEligibleSchemes).Methods("GET")
	r.HandleFunc("/api/applications", h.GetAllApplications).Methods("GET")
	r.HandleFunc("/api/applications", h.CreateApplication).Methods("POST")

	log.Printf("Starting server")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
