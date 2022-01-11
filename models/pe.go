package models

import (
	"time"

	"github.com/lib/pq"
)

type PExchange struct {
	Id                     string
	PresentationSubmission *VerifiablePresentation
	Validations            *VerificationResult
	PresentationDefinition *PresentationDefinition
	RequestedAt            *time.Time
	CreatedAt              *time.Time
	UpdatedAt              *time.Time
	ExpiredAt              *time.Time
	DeletedAt              *pq.NullTime
}

func (pe *PExchange) Valid() bool {
	return pe.Validations != nil && pe.Validations.Valid()
}
