package models

type DataAgreement struct {
	Context              *SSIContext     `json:"@context,omitempty"`
	DataHolder           string          `json:"data_holder,omitempty"`
	DataReceiver         DataReceiver    `json:"data_receiver,omitempty"`
	DataSubject          string          `json:"data_subject,omitempty"`
	Dpia                 *Dpia           `json:"dpia,omitempty"`
	Event                []Event         `json:"event,omitempty"`
	ID                   string          `json:"id"`
	PersonalData         []PersonalDatum `json:"personal_data"`
	Purposes             []Purpose       `json:"purposes"`
	TemplateID           string          `json:"template_id,omitempty"`
	TemplateVersion      string          `json:"template_version,omitempty"`
	TerminationTimestamp int64           `json:"termination_timestamp,omitempty"`
	Version              string          `json:"version,omitempty"`
}

func (da *DataAgreement) GetProofs() *SSIProof {
	last := len(da.Event) - 1
	if last < 0 {
		return nil
	}
	p := da.Event[last].Proof
	return p
}

func (da *DataAgreement) GetProofChain() *[]Proof {
	return nil
}

func (da *DataAgreement) SetProofs(proof *SSIProof) {
	last := len(da.Event) - 1
	if last < 0 {
		return
	}
	da.Event[last].Proof = proof
}

func (da *DataAgreement) SetProofChain(proof *[]Proof) {
	//No support
}

func (da *DataAgreement) GetContext() *SSIContext {
	return da.Context
}

func (da *DataAgreement) SetContexts(context *SSIContext) {
	da.Context = context
}

type DataReceiver struct {
	ID              string `json:"id,omitempty"`
	ConsentDuration int64  `json:"consent_duration,omitempty"`
	FormOfConsent   string `json:"form_of_consent,omitempty"`
	Name            string `json:"name,omitempty"`
	Service         string `json:"service,omitempty"`
	URL             string `json:"url,omitempty"`
}

type Dpia struct {
	DpiaDate       string `json:"dpia_date,omitempty"`
	DpiaSummaryURL string `json:"dpia_summary_url,omitempty"`
}

type Event struct {
	PrincipleDid string    `json:"principle_did,omitempty"`
	Proof        *SSIProof `json:"proof,omitempty"`
	State        string    `json:"state,omitempty"`
	Version      string    `json:"version,omitempty"`
	TimeStamp    int64     `json:"timestamp,omitempty"`
}

type PersonalDatum struct {
	AttributeID        string   `json:"attribute_id,omitempty"`
	AttributeName      string   `json:"attribute_name,omitempty"`
	AttributeSensitive bool     `json:"attribute_sensitive,omitempty"`
	Purposes           []string `json:"purposes,omitempty"`
}

type Purpose struct {
	DataPolicy         DataPolicy `json:"data_policy,omitempty"`
	ID                 string     `json:"id,omitempty"`
	LegalBasis         string     `json:"legal_basis,omitempty"`
	MethodOfUse        string     `json:"method_of_use,omitempty"`
	PurposeCategory    string     `json:"purpose_category,omitempty"`
	PurposeDescription string     `json:"purpose_description,omitempty"`
}

type DataPolicy struct {
	DataRetentionPeriod   int      `json:"data_retention_period,omitempty"`
	GeographicRestriction string   `json:"geographic_restriction,omitempty"`
	IndustryScope         string   `json:"industry_scope,omitempty"`
	Jurisdictions         []string `json:"jurisdictions,omitempty"`
	PolicyURL             string   `json:"policy_URL,omitempty"`
	StorageLocation       string   `json:"storage_location,omitempty"`
}
