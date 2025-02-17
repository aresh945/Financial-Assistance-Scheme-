package service

import (
	"context"
	"financial_assistance/internal/models"
	"financial_assistance/internal/repository"

	"github.com/google/uuid"
)

type Service struct {
	applicantRepo   repository.ApplicantRepository
	schemeRepo      repository.SchemeRepository
	applicationRepo repository.ApplicationRepository
}

func NewService(
	applicantRepo repository.ApplicantRepository,
	schemeRepo repository.SchemeRepository,
	applicationRepo repository.ApplicationRepository,
) *Service {
	return &Service{
		applicantRepo:   applicantRepo,
		schemeRepo:      schemeRepo,
		applicationRepo: applicationRepo,
	}
}

func (s *Service) GetAllApplicants(ctx context.Context) ([]models.Applicant, error) {
	return s.applicantRepo.GetAllApplicants(ctx)
}

func (s *Service) CreateApplicant(ctx context.Context, applicant *models.Applicant) error {
	return s.applicantRepo.CreateApplicant(ctx, applicant)
}

func (s *Service) GetAllSchemes(ctx context.Context) ([]models.Scheme, error) {
	return s.schemeRepo.GetAllSchemes(ctx)
}

func (s *Service) GetEligibleSchemes(ctx context.Context, applicantID uuid.UUID) ([]models.Scheme, error) {
	return s.schemeRepo.GetEligibleSchemes(ctx, applicantID)
}

func (s *Service) CreateApplication(ctx context.Context, application *models.Application) error {
	return s.applicationRepo.CreateApplication(ctx, application)
}

func (s *Service) GetAllApplications(ctx context.Context) ([]models.Application, error) {
	return s.applicationRepo.GetAllApplications(ctx)
}

func (s *Service) CreateScheme(ctx context.Context, scheme *models.Scheme) error {
	return s.schemeRepo.CreateScheme(ctx, scheme)
}
