package models

import "time"

type VerifiableStatus struct {
	Description          string                       `json:"description,omitempty"`
	Id                   string                       `json:"id"`
	Proof                string                       `json:"issued,omitempty"`
	VerifiableCredential []VerifiableCredentialStatus `json:"verifiableCredential"`
}

type VerifiableCredentialStatus struct {
	Claim   *Claim     `json:"claim"`
	Issued  *time.Time `json:"issued,omitempty"`
	Issuer  string     `json:"issuer,omitempty"`
	Proof   string     `json:"proof,omitempty"`
	Updated *time.Time `json:"updated,omitempty"`
}

type Claim struct {
	CurrentStatus string `json:"currentStatus"`
	Id            string `json:"id"`
	StatusReason  string `json:"statusReason,omitempty"`
}
