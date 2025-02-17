package postgres

import (
	"context"
	"database/sql"
	"financial_assistance/internal/models"

	"github.com/google/uuid"
)

type ApplicantRepo struct {
	db *sql.DB
}

func NewApplicantRepo(db *sql.DB) *ApplicantRepo {
	return &ApplicantRepo{db: db}
}

func (r *ApplicantRepo) CreateApplicant(ctx context.Context, applicant *models.Applicant) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO applicants (id, name, employment_status, marital_status, sex, date_of_birth)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err = tx.ExecContext(ctx, query,
		applicant.ID,
		applicant.Name,
		applicant.EmploymentStatus,
		applicant.MaritalStatus,
		applicant.Sex,
		applicant.DateOfBirth,
	)
	if err != nil {
		return err
	}

	if len(applicant.HouseholdMembers) > 0 {
		memberQuery := `
            INSERT INTO household_members (id, name, employment_status, sex, date_of_birth, relation, school_level, applicant_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `
		for _, member := range applicant.HouseholdMembers {
			_, err = tx.ExecContext(ctx, memberQuery,
				member.ID,
				member.Name,
				member.EmploymentStatus,
				member.Sex,
				member.DateOfBirth,
				member.Relation,
				member.SchoolLevel,
				applicant.ID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *ApplicantRepo) GetAllApplicants(ctx context.Context) ([]models.Applicant, error) {
	query := `
        SELECT id, name, employment_status, marital_status, sex, date_of_birth 
        FROM applicants
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applicants []models.Applicant
	for rows.Next() {
		var app models.Applicant
		if err := rows.Scan(
			&app.ID,
			&app.Name,
			&app.EmploymentStatus,
			&app.MaritalStatus,
			&app.Sex,
			&app.DateOfBirth,
		); err != nil {
			return nil, err
		}

		membersQuery := `
            SELECT id, name, employment_status, sex, date_of_birth, relation, school_level
            FROM household_members
            WHERE applicant_id = $1
        `

		memberRows, err := r.db.QueryContext(ctx, membersQuery, app.ID)
		if err != nil {
			return nil, err
		}
		defer memberRows.Close()

		var members []models.HouseholdMember
		for memberRows.Next() {
			var member models.HouseholdMember
			if err := memberRows.Scan(
				&member.ID,
				&member.Name,
				&member.EmploymentStatus,
				&member.Sex,
				&member.DateOfBirth,
				&member.Relation,
				&member.SchoolLevel,
			); err != nil {
				return nil, err
			}
			member.ApplicantID = app.ID
			members = append(members, member)
		}

		app.HouseholdMembers = members
		applicants = append(applicants, app)
	}

	return applicants, nil
}

func (r *ApplicantRepo) GetApplicant(ctx context.Context, id uuid.UUID) (*models.Applicant, error) {
	query := `
        SELECT id, name, employment_status, marital_status, sex, date_of_birth 
        FROM applicants 
        WHERE id = $1
    `

	var app models.Applicant
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&app.ID,
		&app.Name,
		&app.EmploymentStatus,
		&app.MaritalStatus,
		&app.Sex,
		&app.DateOfBirth,
	)
	if err != nil {
		return nil, err
	}

	membersQuery := `
        SELECT id, name, employment_status, sex, date_of_birth, relation, school_level
        FROM household_members
        WHERE applicant_id = $1
    `

	rows, err := r.db.QueryContext(ctx, membersQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.HouseholdMember
	for rows.Next() {
		var member models.HouseholdMember
		if err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.EmploymentStatus,
			&member.Sex,
			&member.DateOfBirth,
			&member.Relation,
			&member.SchoolLevel,
		); err != nil {
			return nil, err
		}
		member.ApplicantID = id
		members = append(members, member)
	}

	app.HouseholdMembers = members
	return &app, nil
}
