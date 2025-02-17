package models

import (
	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `json:"application_id" db:"application_id"`
	ApplicantID uuid.UUID `json:"applicant_id" db:"applicant_id"`
	SchemeID    uuid.UUID `json:"scheme_id" db:"scheme_id"`
	Status      string    `json:"status" db:"status"`
}
