package models

import (
	"github.com/google/uuid"
)

type Scheme struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Criteria Criteria  `json:"criteria,omitempty"`
	Benefits []Benefit `json:"benefits,omitempty"`
}

type Criteria struct {
	EmploymentStatus string `json:"employment_status" db:"employment_status"`
	MaritalStatus    string `json:"marital_status" db:"marital_status"`
	HasChildren      bool   `json:"has_children" db:"has_children"`
}

type Benefit struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Amount float64   `json:"amount" db:"amount"`
}
