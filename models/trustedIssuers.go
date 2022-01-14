package models

type TrustedIssuerList struct {
	Proof          *Proof          `json:"proof,omitempty"`
	TrustedIssuers []TrustedIssuer `json:"trustedIssuers,omitempty"`
}

type TrustedIssuer struct {
	Context          string            `json:"@context,omitempty"`
	Accreditations   *Accreditations   `json:"accreditations,omitempty"`
	EidasCert        *EidasCertificate `json:"eidasCertificate,omitempty"`
	Dids             []string          `json:"dids,omitempty"`
	Domain           string            `json:"domain,omitempty"`
	Id               string            `json:"id,omitempty"`
	OrgInfo          *OrganizationInfo `json:"organizationInfo,omitempty"`
	Proof            *Proof            `json:"proof,omitempty"`
	ServiceEndpoints string            `json:"serviceEndpoints,omitempty"`
}

type EidasCertificate struct {
	IssuerNumber string `json:"eidasCertificateIssuerNumber,omitempty"`
	Pem          string `json:"eidasCertificatePem,omitempty"`
	SerialNumber string `json:"eidasCertificateSerialNumber,omitempty"`
}

type ServiceEndpoint struct {
	Id           string `json:"id,omitempty"`
	ServEndpoint string `json:"serviceEndpoint,omitempty"`
	Type         string `json:"type,omitempty"`
}

type OrganizationInfo struct {
	AreaGroup               string            `json:"areaGroup,omitempty"`
	EORI                    string            `json:"EORI,omitempty"`
	DiscoveryURL            string            `json:"discoveryURL,omitempty"`
	DomainName              string            `json:"domainName,omitempty"`
	IdentifierBag           string            `json:"identifierBag,omitempty"`
	LegalAddress            string            `json:"legalAddress,omitempty"`
	LegalName               string            `json:"legalName,omitempty"`
	LegalPersonalIdentifier string            `json:"legalPersonalIdentifier,omitempty"`
	LEI                     string            `json:"LEI,omitempty"`
	PrefDisplay             *PreferredDisplay `json:"preferredDisplay,omitempty"`
	SEED                    string            `json:"SEED,omitempty"`
	SIC                     string            `json:"SIC,omitempty"`
	TaxReference            string            `json:"taxReference,omitempty"`
	VATRegistration         string            `json:"VATRegistration,omitempty"`
}

type PreferredDisplay struct {
	Background    *Background `json:"background,omitempty"`
	Logo          string      `json:"logo,omitempty"`
	PreferredName string      `json:"preferredName,omitempty"`
	Style         *Style      `json:"style,omitempty"`
}

type Background struct {
	Color string `json:"color,omitempty"`
	URI   string `json:"uri,omitempty"`
}

type Style struct {
	Color      string `json:"color,omitempty"`
	FontFamily string `json:"fontFamily,omitempty"`
}

type Accreditations struct {
	Accreditor       string              `json:"accreditor,omitempty"`
	CreatedAt        *TimeWithFormat     `json:"createdAt,omitempty"`
	CredentialSchema string              `json:"credentialSchema,omitempty"`
	Evidence         *CredentialEvidence `json:"evidence,omitempty"`
	ExpirationDate   *TimeWithFormat     `json:"expirationDate,omitempty"`
	LevelOfTrust     int                 `json:"levelOfTrust,omitempty"`
	Proof            *Proof              `json:"proof,omitempty"`
	ValidFrom        *TimeWithFormat     `json:"validFrom,omitempty"`
}

type CredentialEvidence struct {
	EvidenceDocs     []string `json:"evidenceDocuments,omitempty"`
	DocumentPresence string   `json:"documentPresence,omitempty"`
	Id               string   `json:"id,omitempty"`
	SubjectPresence  string   `json:"subjectPresence,omitempty"`
	Type             []string `json:"type,omitempty"`
}
