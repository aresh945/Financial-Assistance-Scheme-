package postgres

import (
	"context"
	"database/sql"
	"financial_assistance/internal/models"

	"github.com/google/uuid"
)

type SchemeRepo struct {
	db *sql.DB
}

func NewSchemeRepo(db *sql.DB) *SchemeRepo {
	return &SchemeRepo{db: db}
}

func (r *SchemeRepo) GetAllSchemes(ctx context.Context) ([]models.Scheme, error) {
	query := `
        SELECT s.id, s.name, 
               c.employment_status, c.marital_status, c.has_children
        FROM schemes s
        LEFT JOIN criteria c ON s.id = c.scheme_id
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemes []models.Scheme
	for rows.Next() {
		var scheme models.Scheme
		var criteria models.Criteria
		var hasChildren sql.NullBool

		err := rows.Scan(
			&scheme.ID,
			&scheme.Name,
			&criteria.EmploymentStatus,
			&criteria.MaritalStatus,
			&hasChildren,
		)
		if err != nil {
			return nil, err
		}

		criteria.HasChildren = hasChildren.Bool
		scheme.Criteria = criteria
		schemes = append(schemes, scheme)
	}

	for i, scheme := range schemes {
		benefitsQuery := `
            SELECT id, name, amount
            FROM benefits
            WHERE scheme_id = $1
        `

		benefitRows, err := r.db.QueryContext(ctx, benefitsQuery, scheme.ID)
		if err != nil {
			return nil, err
		}
		defer benefitRows.Close()

		var benefits []models.Benefit
		for benefitRows.Next() {
			var benefit models.Benefit
			if err := benefitRows.Scan(&benefit.ID, &benefit.Name, &benefit.Amount); err != nil {
				return nil, err
			}
			benefits = append(benefits, benefit)
		}
		schemes[i].Benefits = benefits
	}

	return schemes, nil
}

func (r *SchemeRepo) GetScheme(ctx context.Context, id uuid.UUID) (*models.Scheme, error) {

	query := `
        SELECT s.id, s.name, 
               c.employment_status, c.marital_status,
               b.id, b.name, b.amount
        FROM schemes s
        LEFT JOIN criteria c ON s.id = c.scheme_id
        LEFT JOIN benefits b ON s.id = b.scheme_id
        WHERE s.id = $1
    `

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheme models.Scheme
	var hasSchemeName bool

	for rows.Next() {
		var criteria models.Criteria
		var benefit models.Benefit

		if err := rows.Scan(
			&scheme.ID, &scheme.Name,
			&criteria.EmploymentStatus, &criteria.MaritalStatus,
			&benefit.ID, &benefit.Name, &benefit.Amount,
		); err != nil {
			return nil, err
		}

		if !hasSchemeName {
			scheme.Criteria = criteria
			hasSchemeName = true
		}
		scheme.Benefits = append(scheme.Benefits, benefit)
	}

	if !hasSchemeName {
		return nil, sql.ErrNoRows
	}

	return &scheme, nil
}

func (r *SchemeRepo) GetEligibleSchemes(ctx context.Context, applicantID uuid.UUID) ([]models.Scheme, error) {
	query := `
    SELECT DISTINCT s.id, s.name 
    FROM schemes s
    JOIN criteria c ON s.id = c.scheme_id
    JOIN applicants a ON 
        (c.employment_status IS NULL OR a.employment_status = c.employment_status)
        AND (c.marital_status IS NULL OR a.marital_status = c.marital_status)
        AND (
            (c.has_children = false AND NOT EXISTS (
                SELECT 1 FROM household_members hm 
                WHERE hm.applicant_id = a.id 
                AND (hm.relation = 'son' OR hm.relation = 'daughter')
            )) 
            OR 
            (c.has_children = true AND EXISTS (
                SELECT 1 FROM household_members hm 
                WHERE hm.applicant_id = a.id 
                AND (hm.relation = 'son' OR hm.relation = 'daughter')
            ))
            OR 
            c.has_children IS NULL
        )
    WHERE a.id = $1
`

	rows, err := r.db.QueryContext(ctx, query, applicantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemes []models.Scheme
	for rows.Next() {
		var scheme models.Scheme
		if err := rows.Scan(&scheme.ID, &scheme.Name); err != nil {
			return nil, err
		}
		schemes = append(schemes, scheme)
	}

	return schemes, nil
}

func (r *SchemeRepo) CreateScheme(ctx context.Context, scheme *models.Scheme) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	schemeQuery := `
        INSERT INTO schemes (id, name)
        VALUES ($1, $2)
    `
	_, err = tx.ExecContext(ctx, schemeQuery, scheme.ID, scheme.Name)
	if err != nil {
		return err
	}

	criteriaQuery := `
        INSERT INTO criteria (scheme_id, employment_status, marital_status)
        VALUES ($1, $2, $3)
    `
	_, err = tx.ExecContext(ctx, criteriaQuery,
		scheme.ID,
		scheme.Criteria.EmploymentStatus,
		scheme.Criteria.MaritalStatus,
	)
	if err != nil {
		return err
	}

	benefitQuery := `
        INSERT INTO benefits (id, name, amount)
        VALUES ($1, $2, $3)
    `
	for _, benefit := range scheme.Benefits {
		_, err = tx.ExecContext(ctx, benefitQuery,
			benefit.ID,
			benefit.Name,
			benefit.Amount,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
