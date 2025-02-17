package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Applicant struct {
	ID               uuid.UUID         `json:"id" db:"id"`
	Name             string            `json:"name" db:"name"`
	EmploymentStatus string            `json:"employment_status" db:"employment_status"`
	MaritalStatus    string            `json:"marital_status" db:"marital_status"`
	Sex              string            `json:"sex" db:"sex"`
	DateOfBirth      time.Time         `json:"date_of_birth" db:"date_of_birth"`
	HouseholdMembers []HouseholdMember `json:"household,omitempty"`
}

type HouseholdMember struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	EmploymentStatus string    `json:"employment_status" db:"employment_status"`
	Sex              string    `json:"sex" db:"sex"`
	DateOfBirth      time.Time `json:"date_of_birth" db:"date_of_birth"`
	Relation         string    `json:"relation" db:"relation"`
	SchoolLevel      string    `json:"school_level" db:"school_level"`
	ApplicantID      uuid.UUID `json:"applicant_id" db:"applicant_id"`
}

func (a *Applicant) UnmarshalJSON(data []byte) error {
	type Alias Applicant
	aux := &struct {
		DateOfBirth string `json:"date_of_birth"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var parsedTime time.Time
	var err error

	parsedTime, err = time.Parse(time.RFC3339, aux.DateOfBirth)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02", aux.DateOfBirth)
	}

	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	a.DateOfBirth = parsedTime

	return nil
}
