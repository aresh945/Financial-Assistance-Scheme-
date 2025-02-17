package repository

import (
	"context"
	"financial_assistance/internal/models"

	"github.com/google/uuid"
)

type ApplicantRepository interface {
	CreateApplicant(ctx context.Context, applicant *models.Applicant) error
	GetApplicant(ctx context.Context, id uuid.UUID) (*models.Applicant, error)
	GetAllApplicants(ctx context.Context) ([]models.Applicant, error)
}

type SchemeRepository interface {
	GetScheme(ctx context.Context, id uuid.UUID) (*models.Scheme, error)
	GetAllSchemes(ctx context.Context) ([]models.Scheme, error)
	GetEligibleSchemes(ctx context.Context, applicantID uuid.UUID) ([]models.Scheme, error)
	CreateScheme(ctx context.Context, scheme *models.Scheme) error
}

type ApplicationRepository interface {
	CreateApplication(ctx context.Context, application *models.Application) error
	GetApplication(ctx context.Context, id uuid.UUID) (*models.Application, error)
	GetAllApplications(ctx context.Context) ([]models.Application, error)
}
