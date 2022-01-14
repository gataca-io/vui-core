package models

import (
	"encoding/json"
	"errors"
)

type VerifiablePresentation struct {
	Context                *SSIContext             `json:"@context,omitempty" example:"https://www.w3.org/2018/credentials/v1" description:"Context for JSON-LD"`
	DataAgreementId        string                  `json:"data_agreement_id,omitempty" example:"da:gatc:ehgiuwg39487wq9gf7a47af37" description:"Id of the data agreement supporting this exchange"`
	Holder                 *string                 `json:"holder,omitempty" example:"did:gatc:example1234567" description:"DID of the Holder of the credentials"`
	PresentationSubmission *PresentationSubmission `json:"presentation_submission,omitempty" description:"Presentation submission according to DIF PE"`
	Proof                  *SSIProof               `json:"proof,omitempty" description:"Proofs to verify the presentation"`
	Type                   []string                `json:"type,omitempty" example:"VerifiablePresentation" description:"Definition of the format of the presentation"`
	VerifiableCredential   []VerifiableCredential  `json:"verifiableCredential,omitempty" description:"List of Verifiable Credentials included in the presentation"`
}

func (v *VerifiablePresentation) GetProofs() *SSIProof {
	return v.Proof
}

func (v *VerifiablePresentation) GetProofChain() *[]Proof {
	return nil
}

func (v *VerifiablePresentation) SetProofs(proof *SSIProof) {
	v.Proof = proof
}

func (v *VerifiablePresentation) SetProofChain(proof *[]Proof) {
	//No support
}

func (v *VerifiablePresentation) GetContext() *SSIContext {
	return v.Context
}

func (v *VerifiablePresentation) SetContexts(context *SSIContext) {
	v.Context = context
}

func (v *VerifiablePresentation) IsResponse() bool {
	return true
}

func (v *VerifiablePresentation) ToPresentation() *VerifiablePresentation {
	return v
}

type VerifiableCredential struct {
	Context           *SSIContext             `json:"@context,omitempty" example:"https://www.w3.org/2018/credentials/v1" description:"Context for JSON-LD"`
	CredentialSchema  *CredentialSchema       `json:"credentialSchema,omitempty" swaggertype:"object,string" description:"Definition to retrieve the schema of the credential"` //Reusing the CredentialStatusType coz its the same struct here
	CredentialStatus  *CredentialStatus       `json:"credentialStatus,omitempty" swaggertype:"object,string" description:"Definition to retrieve the current status of the credential"`
	CredentialSubject *map[string]interface{} `json:"credentialSubject,omitempty" swaggerignore:"true" example:"id:did_example_xxxxxxxxxxx,email:example@domain.com" description:"Claims in free format stated about the subject. Linked to the credential type."` //swagger ignore because of unexpected problems with current library
	Evidence          *Evidence               `json:"evidence,omitempty" swaggertype:"object,string" description:"Definition of the evidence. Required by eIDas"`
	ExpirationDate    *TimeWithFormat         `json:"expirationDate,omitempty" swaggertype:"string" example:"2019-10-01T12:12:15.999Z" description:"Timestamp of expiration of the credential"`
	Id                string                  `json:"id" example:"cred:example:zzzzzzzzzzzz" description:"Unique identifier of the Verifiable Credential"`
	Iss               string                  `json:"iss,omitempty" description:"jwt fields needed for VA"`
	IssuanceDate      *TimeWithFormat         `json:"issuanceDate,omitempty"  swaggertype:"string" example:"2019-10-01T12:12:05.999Z" description:"Timestamp of issuance of the credential"`
	Issuer            string                  `json:"issuer,omitempty" example:"did:example:yyyyyyyyyyyyyyyy" description:"Issuer of the credential"`
	Proof             *SSIProof               `json:"proof,omitempty"  description:"Proofs to verify the presentation"`
	Type              []string                `json:"type,omitempty" example:"emailCredential" description:"Type definition of this verifiable credential stablishing a specific json schema."`
	ValidFrom         *TimeWithFormat         `json:"validFrom,omitempty" swaggertype:"string" example:"2019-10-01T12:12:05.999Z" description:"Timestamp from which the credential its valid"`
}

func (v *VerifiableCredential) GetProofs() *SSIProof {
	return v.Proof
}

func (v *VerifiableCredential) GetProofChain() *[]Proof {
	return nil
}

func (v *VerifiableCredential) SetProofs(proof *SSIProof) {
	v.Proof = proof
}

func (v *VerifiableCredential) SetProofChain(proof *[]Proof) {
	//No support
}

func (v *VerifiableCredential) GetContext() *SSIContext {
	return v.Context
}

func (v *VerifiableCredential) SetContexts(context *SSIContext) {
	v.Context = context
}

func (v *VerifiableCredential) GetSchemaRef() string {
	if v.CredentialSchema != nil {
		return v.CredentialSchema.Id
	}
	return ""
}
func (v *VerifiableCredential) GetSchema() string {
	return ""
}

func (v *VerifiableCredential) IsRef() bool {
	return true
}

type Evidence struct {
	DocumentPresence string            `json:"documentPresence,omitempty" example:"Physical"`
	EvidenceDocument *EvidenceDocument `json:"evidenceDocument,omitempty"`
	Id               string            `json:"id,omitempty" example:"https://base-uri/evidence/f2ae...5678"`
	SubjectPresence  string            `json:"subjectPresence,omitempty" example:"Physical"`
	Type             []string          `json:"type,omitempty" example:"DocumentVerification,PassportVerification"`
	Verifier         string            `json:"verifier,omitempty" example:"did:ebsi:xxxxxxx"`
}

type EvidenceDocument struct {
	DocumentCode           string          `json:"documentCode,omitempty" example:"P"`
	DocumentExpirationDate *TimeWithFormat `json:"documentExpirationDate,omitempty"` //Check if its RFC3339
	DocumentIssuingState   string          `json:"documentIssuingState,omitempty" example:"NLD"`
	DocumentNumber         string          `json:"documentNumber,omitempty" example:"SPECI2014"`
	Type                   string          `json:"type,omitempty" example:"Passport"`
}

type CredentialStatus struct {
	Id   string `json:"id,omitempty" example:"https://issuer.domain.com/cred:example:zzzzzzzzzz" description:"URI to query the current credential status"`
	Type string `json:"type,omitempty" example:"CredentialStatusList2017" description:"Credential Status Protocol definition"`
}

type CredentialSchema = CredentialStatus

type Proof struct {
	CaDES              string          `json:"cades,omitempty" example:"308204c906092a864886f70d010702...266ad9fee3375d8095" description:"Proof Value for ADes signatures" `
	Challenge          string          `json:"challenge,omitempty" example:"TyYfomXjwPaQoSRzCZk7CxFYR8DwAigt" description:"Challenge enforcement of a nonce to avoid replay attacks."`
	Context            *SSIContext     `json:"@context,omitempty" example:"https://www.w3.org/2018/credentials/v1" description:"Context for JSON-LD"`
	Created            *TimeWithFormat `json:"created,omitempty" swaggertype:"string" example:"2019-10-01T12:12:05.999Z" description:"Timestamp of signature of the proof"`
	Creator            string          `json:"creator,omitempty" example:"did:gatc:yyyyyyyyyyyy#keys-1" description:"URI of the key used to sign the proof."`
	Domain             string          `json:"domain,omitempty" description:""`
	Jws                string          `json:"jws,omitempty" swaggerignore:"true" example:"eyJhbGciOiJQUzUxMiIsImtpZCI6IiJ9.eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZVByZXNlbnRhdGlvbiIsIkNyZWRlbnRpYWxNYW5hZ2VyUHJlc2VudGF0aW9uIl0sInZlcmlmaWFibGVDcmVkZW50aWFsIjpbeyJpZCI6ImRpZDpnYXRjOmU1ZDRkMWFlODg4MGM0NDMwNGJiYzU4MjBkMjZiOWE0OmRpZDpnYXRjOjhjZDkzMTMxYjhjOTk0YTAyZTE4ODU0MTAxYTM4YTk3IiwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIk5hbWVDcmVkZW50aWFsIl0sImlzc3VlciI6Imh0dHA6Ly9nYXRhY2EtYmFja2JvbmUuZ2F0YWNhaWQuY29tOjkwOTAvYXBpL3YxL2RpZHMvZGlkOmdhdGM6ZTVkNGQxYWU4ODgwYzQ0MzA0YmJjNTgyMGQyNmI5YTQiLCJpc3N1YW5jZURhdGUiOiIyMDE3LTAyLTIzVDIwOjQ2OjMwLjE5OTUwOTE4OVoiLCJjcmVkZW50aWFsU3ViamVjdCI6eyJpZCI6ImRpZDpnYXRjOjhjZDkzMTMxYjhjOTk0YTAyZTE4ODU0MTAxYTM4YTk3IiwibmFtZSI6Impvc2UifX0seyJpZCI6ImRpZDpnYXRjOmU1ZDRkMWFlODg4MGM0NDMwNGJiYzU4MjBkMjZiOWE0OmRpZDpnYXRjOjhjZDkzMTMxYjhjOTk0YTAyZTE4ODU0MTAxYTM4YTk3IiwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIkVtYWlsQ3JlZGVudGlhbCJdLCJpc3N1ZXIiOiJodHRwOi8vZ2F0YWNhLWJhY2tib25lLmdhdGFjYWlkLmNvbTo5MDkwL2FwaS92MS9kaWRzL2RpZDpnYXRjOmU1ZDRkMWFlODg4MGM0NDMwNGJiYzU4MjBkMjZiOWE0IiwiaXNzdWFuY2VEYXRlIjoiMjAxNS0wOS0wOFQyMDo0NjozMC4xOTk1MDkxODlaIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiZW1haWwiOiJqb3NlQGdhdGFjYWlkLmNvbSIsImlkIjoiZGlkOmdhdGM6OGNkOTMxMzFiOGM5OTRhMDJlMTg4NTQxMDFhMzhhOTcifX0seyJpZCI6ImRpZDpnYXRjOmU1ZDRkMWFlODg4MGM0NDMwNGJiYzU4MjBkMjZiOWE0OmRpZDpnYXRjOjhjZDkzMTMxYjhjOTk0YTAyZTE4ODU0MTAxYTM4YTk3IiwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIkltYWdlQ3JlZGVudGlhbCJdLCJpc3N1ZXIiOiJodHRwOi8vZ2F0YWNhLWJhY2tib25lLmdhdGFjYWlkLmNvbTo5MDkwL2FwaS92MS9kaWRzL2RpZDpnYXRjOmU1ZDRkMWFlODg4MGM0NDMwNGJiYzU4MjBkMjZiOWE0IiwiaXNzdWFuY2VEYXRlIjoiMjAxOS0wMS0xMVQyMDo0NjozMC4xOTk1MDkxODlaIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiaWQiOiJkaWQ6Z2F0Yzo4Y2Q5MzEzMWI4Yzk5NGEwMmUxODg1NDEwMWEzOGE5NyIsImltYWdlIjoiYmFzZTY0In19LHsiaWQiOiJkaWQ6Z2F0YzplNWQ0ZDFhZTg4ODBjNDQzMDRiYmM1ODIwZDI2YjlhNDpkaWQ6Z2F0Yzo4Y2Q5MzEzMWI4Yzk5NGEwMmUxODg1NDEwMWEzOGE5NyIsInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiLCJDb3VudHJ5Q3JlZGVudGlhbCJdLCJpc3N1ZXIiOiJodHRwOi8vZ2F0YWNhLWJhY2tib25lLmdhdGFjYWlkLmNvbTo5MDkwL2FwaS92MS9kaWRzL2RpZDpnYXRjOmU1ZDRkMWFlODg4MGM0NDMwNGJiYzU4MjBkMjZiOWE0IiwiaXNzdWFuY2VEYXRlIjoiMjAxOS0wMS0wOFQyMDo0NjozMC4xOTk1MDkxODlaIiwiY3JlZGVudGlhbFN1YmplY3QiOnsiY291bnRyeSI6IkVTIiwiaWQiOiJkaWQ6Z2F0Yzo4Y2Q5MzEzMWI4Yzk5NGEwMmUxODg1NDEwMWEzOGE5NyJ9fSx7ImlkIjoiZGlkOmdhdGM6ZTVkNGQxYWU4ODgwYzQ0MzA0YmJjNTgyMGQyNmI5YTQ6ZGlkOmdhdGM6OGNkOTMxMzFiOGM5OTRhMDJlMTg4NTQxMDFhMzhhOTciLCJ0eXBlIjpbIlZlcmlmaWFibGVDcmVkZW50aWFsIiwiUGhvbmVDcmVkZW50aWFsIl0sImlzc3VlciI6Imh0dHA6Ly9nYXRhY2EtYmFja2JvbmUuZ2F0YWNhaWQuY29tOjkwOTAvYXBpL3YxL2RpZHMvZGlkOmdhdGM6ZTVkNGQxYWU4ODgwYzQ0MzA0YmJjNTgyMGQyNmI5YTQiLCJpc3N1YW5jZURhdGUiOiIyMDE4LTAzLTMwVDIwOjQ2OjMwLjE5OTUwOTE4OVoiLCJjcmVkZW50aWFsU3ViamVjdCI6eyJpZCI6ImRpZDpnYXRjOjhjZDkzMTMxYjhjOTk0YTAyZTE4ODU0MTAxYTM4YTk3IiwicGhvbmUiOiI2MTIgMzQgNTYgNzgifX1dLCJwcm9vZiI6eyJ0eXBlIjoiUlNBU1NBLVBTUyB1c2luZyBTSEE1MTIgYW5kIE1HRjEtU0hBNTEyIiwiY3JlYXRlZCI6IjIwMTktMDYtMjNUMjA6NDY6MzAuMTk5NTA5MTg5WiIsInByb29mUHVycG9zZSI6ImF1dGhlbnRpY2F0aW9uIiwiY2hhbGxlbmdlIjoiMzA0ZWI2NjItZDYwMS00NjU2LTljMDYtNTg2YmY1ZjQyYTdmIiwiZG9tYWluIjoiZ2F0YWNhLWNvbm5lY3QuZ2F0YWNhaWQuY29tIn19.u8MLlbZzqGklCq-dnhONhzHaH53LO2cYZPLs7Nn68AyG8kZOas3Yb11rWc4DdgdRFTgE-flBI-_yN9mhgjpfd2OivfsJ3zgV0WhJ_76e9kx0OYP9mOthPq4LdraZXzu6syuGPmqSmhWYgnUgNFds7NXZblADO3b_BWAICbOIUnU" description:"Value of the proof in JWS format"`
	Nonce              string          `json:"nonce,omitempty" example:"TyYfomXjwPaQoSRzCZk7CxFYR8DwAigt" description:"Challenge enforcement of a nonce to avoid replay attacks."`
	ProofPurpose       string          `json:"proofPurpose,omitempty" example:"Authentication" description:"Stated usage of the proof"`
	ProofValue         string          `json:"proofValue,omitempty" example:"bQ5AimlvOv6p5wa9pVlmjWgPMr7j9rKw_yjUL6yHlQNwnKk7HL8VQzIT0Xx" description:"Proof value."`
	SignatureValue     string          `json:"signatureValue,omitempty" example:"bQ5AimlvOv6p5wa9pVlmjWgPMr7j9rKw_yjUL6yHlQNwnKk7HL8VQzIT0Xx" description:"Proof value."`
	Type               string          `json:"type,omitempty" example:"Ed25519Signature2018" description:"Cryptographic suite used for signature"`
	VerificationMethod string          `json:"verificationMethod,omitempty" example:"did:gatc:yyyyyyyyyyyy#keys-1" description:"URI of the key used to validate the proof."`
}

func (p *Proof) GetContext() *SSIContext {
	return p.Context
}

func (p *Proof) SetContexts(context *SSIContext) {
	p.Context = context
}

func (p *Proof) GetCreator() string {
	if p.Creator != "" {
		return p.Creator
	}
	if p.VerificationMethod != "" {
		return p.VerificationMethod
	}
	return ""
}

func (p *Proof) GetValue() string {
	if p.SignatureValue != "" {
		return p.SignatureValue
	}
	if p.ProofValue != "" {
		return p.ProofValue
	}
	if p.CaDES != "" {
		return p.CaDES
	}
	if p.Jws != "" {
		return p.Jws
	}
	return ""
}

type SSIProof struct {
	Value  *Proof
	Values *[]Proof
}

func (sp *SSIProof) GetProof() *[]Proof {
	if sp.Values != nil {
		return sp.Values
	}
	if sp.Value != nil {
		arr := []Proof{(*sp.Value)}
		return &arr
	}
	return nil
}

func (sp *SSIProof) GetCreators() []string {
	creators := []string{}
	if sp.Values != nil {
		for _, p := range *sp.Values {
			creators = append(creators, p.GetCreator())
		}
	}
	if sp.Value != nil {
		creators = append(creators, sp.Value.GetCreator())
	}
	return creators
}

func (sp *SSIProof) UnmarshalJSON(jsonData []byte) error {
	p := &Proof{}
	if errO := json.Unmarshal([]byte(jsonData), p); errO != nil {
		// Could be array
		var array []Proof
		if errA := json.Unmarshal([]byte(jsonData), &array); errA != nil {
			// Could be a string
			var s string
			if errS := json.Unmarshal([]byte(jsonData), &s); errS != nil {
				//Not a string, type unkonwn
				return errors.New("unrecognized entity")
			}
			p.ProofValue = s
			sp.Value = p
			return nil
		}
		sp.Values = &array
	}
	sp.Value = p
	return nil
}

// MarshalJSON TODO -When not breaking code in wallet-> Marshal as single object if only one Proof
func (sp SSIProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(sp.GetProof())
}
