package models

import (
	"bytes"
	"encoding/json"
)

type TenantConfig struct {
	//Tenant
	TenantId              string                  `json:"tenantid" example:"my-tenant" description:"Tenant unique identifier"`
	DID                   string                  `json:"did" example:"did:example:xxxxxxxxxxxx" description:"Verifier DID used by the tenant"`
	Domain                string                  `json:"domain" example:"https://host.domain.com" description:"Verifier domain to receive the presentation Response"`
	Credentials           []CredentialRequest     `json:"credentials" description:"Credentials and constraint to be specified on the tenant's Presentation Requests"`
	Security              []SecMechanism          `json:"security" description:"List of security mechanisms to be enforced by the Presentation Response"`
	Callback              string                  `json:"callback" description:"Endpoint from the relying party used to notify when the user makes changes on the data"`
	DataAgreementTemplate *DataAgreement          `json:"dataAgreementTemplate" description:"Template of the Data Agreement associated with this service"`
	ServicePurpose        string                  `json:"service" description:"Description of the service that is being provided with this QR"`
	AdvancedDefinition    *PresentationDefinition `json:"advancedDefinition" description:"Presentation exchange definition at an advanced level for expert admin users"`
}

type CredentialRequest struct {
	Mandatory  bool   `json:"mandatory" example:"true" description:"Mark if this credential must be present on the Presentation Response."`
	Purpose    string `json:"purpose" example:"Authentication" description:"Stated Consent of usage of this particular credential to be agreed upon."`
	TrustLevel int    `json:"trustLevel" example:"1" description:"Minimal level of trust required by this credential"`
	Type       string `json:"type" example:"emailCredential" description:"Credential Type defining the schema and the information about the subject required"`
}

type SecMechanism struct {
	Accepted []string    `json:"accepted,omitempty"`
	Type     MechanismId `json:"type,omitempty"`
}

type MechanismId int

const (
	AuthNFactor MechanismId = iota + 1
	AppAuth
	Credential
)

func (d MechanismId) String() string {
	return [...]string{"AuthNFactor", "AppAuth", "Credential"}[d-1]
}

var toID = map[string]MechanismId{
	"AuthNFactor": AuthNFactor,
	"AppAuth":     AppAuth,
	"Credential":  Credential,
}

// MarshalJSON marshals the enum as a quoted json string
func (s MechanismId) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *MechanismId) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toID[j]
	return nil
}

type AuthNFactorCredential int

const (
	Silent AuthNFactorCredential = iota + 1
	Biometric
	Email
	SMS
	FaceSDK
	RemoteFaceID
)

func (d AuthNFactorCredential) String() string {
	return [...]string{"silent", "biometric", "email", "sms", "faceSDK", "remoteFaceId"}[d-1]
}

var atoID = map[string]AuthNFactorCredential{
	"silent":       Silent,
	"biometric":    Biometric,
	"email":        Email,
	"sms":          SMS,
	"faceSDK":      FaceSDK,
	"remoteFaceId": RemoteFaceID,
}

// MarshalJSON marshals the enum as a quoted json string
func (s AuthNFactorCredential) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func ToAuthNFactorCredential(code string) AuthNFactorCredential {
	return atoID[code]
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *AuthNFactorCredential) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = atoID[j]
	return nil
}
