package models

type CredentialDerivationRequest struct {
	VerifiableCredential *VerifiableCredential `json:"verifiableCredential" description:"Original credential to be derived" `
	Frame                interface{}           `json:"frame" description:"JSON-LD frame used for selective disclosure."`
	Options              *DerivationOptions    `json:"options" description:"Configuration options to include in the derived proof"`
}

type CredentialIssueRequest struct {
	VerifiableCredential *VerifiableCredential `json:"verifiableCredential" description:"Original credential to be Issued" `
	Options              *IssueOptions         `json:"options" description:"Configuration options to include in the issuance proof"`
}

type CredentialStatusUpdateRequest struct {
	CredentialID     *string           `json:"credentialId" description:"Credential unique ID" example:"cred:gatc:exampleABC123" `
	CredentialStatus *CredentialStatus `json:"credentialStatus" description:"Configuration of the credential status RL"`
}

type CredentialVerificationRequest struct {
	VerifiableCredential *VerifiableCredential `json:"verifiableCredential" description:"Original credential to be verified" `
	Options              *Proof                `json:"options" description:"Configuration options to verify the proof"`
}

type PresentationVerificationRequest struct {
	Presentation  *VerifiablePresentation `json:"presentation" description:"Proofless Presentation to be verified" `             //either only a non-verifiable presentation (proofless
	VPresentation *VerifiablePresentation `json:"verifiablePresentation" description:"Verifiable Presentation to be validated" ` //or a verifiable presentation (with options)
	Options       *Proof                  `json:"options" description:"Configuration options to verify the proof"`
}

type PresentationProveRequest struct {
	VPresentation *VerifiablePresentation `json:"presentation" description:"Presentation to be proved" `
	Options       *Proof                  `json:"options" description:"Configuration options to include in the proof"`
}

type IssueOptions struct {
	Proof
	CredentialStatus *CredentialStatus `json:"credentialStatus" description:"Configuration of the credential status RL"`
}

type DerivationOptions struct {
	Nonce string `json:"nonce" description:"Nonce to be included in the proof"`
}

type VerificationResult struct {
	Checks   []string `json:"checks" description:"Security checks performed" example:"['proof']"`
	Warnings []string `json:"warnings" description:"Warning messages to include about the validation" example:"['Context not verified']"`
	Errors   []string `json:"errors" description:"Resulting errors on the validation. Should be empty if the validation is successful." example:"[]"`
}

func (v *VerificationResult) Valid() bool {
	return len(v.Errors) == 0 && len(v.Checks) > 0
}
