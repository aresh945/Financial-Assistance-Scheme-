package handler

import (
	"encoding/json"
	"financial_assistance/internal/models"
	"financial_assistance/internal/service"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateApplicant(w http.ResponseWriter, r *http.Request) {
	var applicant models.Applicant

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received body: %s", string(body))

	err = json.Unmarshal(body, &applicant)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateApplicant(r.Context(), &applicant); err != nil {
		log.Printf("Error creating applicant: %v", err)
		http.Error(w, "Failed to create applicant", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetAllApplicants(w http.ResponseWriter, r *http.Request) {
	applicants, err := h.service.GetAllApplicants(r.Context())
	if err != nil {
		http.Error(w, "Failed to get applicants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicants)
}

func (h *Handler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var application models.Application

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received application JSON: %s", string(body))

	err = json.Unmarshal(body, &application)
	if err != nil {
		log.Printf("Error decoding application JSON: %v", err)
		log.Printf("Problematic JSON: %s", string(body))
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Decoded Application: %+v", application)

	if err := h.service.CreateApplication(r.Context(), &application); err != nil {
		log.Printf("Error creating application: %v", err)
		http.Error(w, "Failed to create application", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetAllSchemes(w http.ResponseWriter, r *http.Request) {
	schemes, err := h.service.GetAllSchemes(r.Context())
	if err != nil {
		log.Printf("Error getting schemes: %v", err)
		http.Error(w, "Failed to get schemes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schemes)
}

func (h *Handler) GetEligibleSchemes(w http.ResponseWriter, r *http.Request) {
	applicantID := r.URL.Query().Get("applicant")
	if applicantID == "" {
		http.Error(w, "Applicant ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(applicantID)
	if err != nil {
		http.Error(w, "Invalid applicant ID format", http.StatusBadRequest)
		return
	}

	schemes, err := h.service.GetEligibleSchemes(r.Context(), id)
	if err != nil {
		log.Printf("Error getting eligible schemes: %v", err)
		http.Error(w, "Failed to get eligible schemes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schemes)
}

func (h *Handler) GetAllApplications(w http.ResponseWriter, r *http.Request) {
	applications, err := h.service.GetAllApplications(r.Context())
	if err != nil {
		log.Printf("Error getting applications: %v", err)
		http.Error(w, "Failed to get applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}

func (h *Handler) CreateScheme(w http.ResponseWriter, r *http.Request) {
	var scheme models.Scheme

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received JSON: %s", string(body))

	err = json.Unmarshal(body, &scheme)
	if err != nil {
		log.Printf("Error decoding scheme: %v", err)
		log.Printf("Scheme JSON: %s", string(body))
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Decoded Scheme: %+v", scheme)
	log.Printf("Criteria Has Children: %v", scheme.Criteria.HasChildren)

	if scheme.Name == "" {
		http.Error(w, "Scheme name is required", http.StatusBadRequest)
		return
	}

	if scheme.ID == uuid.Nil {
		scheme.ID = uuid.New()
	}

	if err := h.service.CreateScheme(r.Context(), &scheme); err != nil {
		log.Printf("Error creating scheme: %v", err)
		http.Error(w, "Failed to create scheme", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(scheme)
}
