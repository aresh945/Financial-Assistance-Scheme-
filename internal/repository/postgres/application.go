package postgres

import (
	"context"
	"database/sql"
	"financial_assistance/internal/models"

	"github.com/google/uuid"
)

type ApplicationRepo struct {
	db *sql.DB
}

func NewApplicationRepo(db *sql.DB) *ApplicationRepo {
	return &ApplicationRepo{db: db}
}

func (r *ApplicationRepo) CreateApplication(ctx context.Context, application *models.Application) error {
	query := `
        INSERT INTO applications (application_id, applicant_id, scheme_id, status)
        VALUES ($1, $2, $3, $4)
    `

	_, err := r.db.ExecContext(ctx, query,
		application.ID,
		application.ApplicantID,
		application.SchemeID,
		application.Status,
	)
	return err
}

func (r *ApplicationRepo) GetApplication(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	query := `
        SELECT application_id, applicant_id, scheme_id, status
        FROM applications
        WHERE application_id = $1
    `

	var app models.Application
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&app.ID,
		&app.ApplicantID,
		&app.SchemeID,
		&app.Status,
	)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *ApplicationRepo) GetAllApplications(ctx context.Context) ([]models.Application, error) {
	query := `
        SELECT application_id, applicant_id, scheme_id, status
        FROM applications
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []models.Application
	for rows.Next() {
		var app models.Application
		if err := rows.Scan(
			&app.ID,
			&app.ApplicantID,
			&app.SchemeID,
			&app.Status,
		); err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, nil
}
