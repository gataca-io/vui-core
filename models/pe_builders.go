package models

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gopkg.in/go-playground/validator.v9"
)

//Duplicated function to avoid import cycles
func Validate(t interface{}) error {
	if err := validator.New().Struct(t); err != nil {
		return err
	}
	return nil
}

type PresentationDefinitionBuilder struct {
	Definition PresentationDefinition
}

func NewPresentationDefinitionBuilder() *PresentationDefinitionBuilder {
	uuid, _ := uuid.NewV4()
	return &PresentationDefinitionBuilder{
		Definition: PresentationDefinition{
			DIFPresentationDefinition: DIFPresentationDefinition{
				ID: uuid.String(),
			},
			Proof:         nil,
			DataAgreement: nil,
		},
	}
}

func (p *PresentationDefinitionBuilder) Build() (*PresentationDefinitionHolder, error) {
	return &PresentationDefinitionHolder{
		PresentationDefinition: p.Definition,
	}, Validate(p.Definition)
}

func (p *PresentationDefinitionBuilder) SetName(name string) {
	p.Definition.Name = name
}

func (p *PresentationDefinitionBuilder) SetID(id string) {
	p.Definition.ID = id
}

func (p *PresentationDefinitionBuilder) SetPurpose(purpose string) {
	p.Definition.Purpose = purpose
}

func (p *PresentationDefinitionBuilder) SetLocale(locale string) {
	p.Definition.Locale = locale
}

// Submission Requirement Builders //

func (p *PresentationDefinitionBuilder) SetJWTFormat(format JWTFormat, algs []string) error {
	if len(algs) < 1 {
		return fmt.Errorf("must set one or more algs for the jwt type<%s>", format)
	}
	if p.Definition.Format == nil {
		p.Definition.Format = &Format{}
	}
	switch format {
	case JWT:
		p.Definition.Format.JWT = &JWTType{Alg: algs}
	case JWTVC:
		p.Definition.Format.JWTVC = &JWTType{Alg: algs}
	case JWTVP:
		p.Definition.Format.JWTVP = &JWTType{Alg: algs}
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
	return nil
}

func (p *PresentationDefinitionBuilder) SetLDPFormat(format LDPFormat, proofTypes []string) error {
	if len(proofTypes) < 1 {
		return fmt.Errorf("must set one or more proof types for the ldp type<%s>", format)
	}
	if p.Definition.Format == nil {
		p.Definition.Format = &Format{}
	}
	switch format {
	case LDP:
		p.Definition.Format.LDP = &LDPType{ProofType: proofTypes}
	case LDPVC:
		p.Definition.Format.LDPVC = &LDPType{ProofType: proofTypes}
	case LDPVP:
		p.Definition.Format.LDPVP = &LDPType{ProofType: proofTypes}
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
	return nil
}

func (p *PresentationDefinitionBuilder) AddSubmissionRequirements(srs ...SubmissionRequirement) error {
	for _, sr := range srs {
		if err := validateSubmissionRequirement(sr); err != nil {
			return err
		}
	}
	p.Definition.SubmissionRequirements = srs
	return nil
}

func validateSubmissionRequirement(sr SubmissionRequirement) error {
	if sr.From == "" && len(sr.FromNested) == 0 {
		return errors.New("from must have a value")
	}
	if sr.From != "" && len(sr.FromNested) > 0 {
		return errors.New("from and from_nested are mutually exclusive fields")
	}
	return Validate(sr)
}

// Input Descriptor Builders //

func (p *PresentationDefinitionBuilder) AddInputDescriptor(i InputDescriptor) error {
	if err := Validate(i); err != nil {
		return err
	}
	p.Definition.InputDescriptors = append(p.Definition.InputDescriptors, i)
	return nil
}

func NewInputDescriptor(id, name, purpose, metadata string) *InputDescriptor {
	return &InputDescriptor{
		ID:       id,
		Name:     name,
		Purpose:  purpose,
		Metadata: metadata,
	}
}

func (i *InputDescriptor) AddSchema(s Schema) error {
	if err := Validate(s); err != nil {
		return err
	}
	i.Schema = append(i.Schema, s)
	return nil
}

func (i *InputDescriptor) SetConstraints(fields ...Field) error {
	for _, f := range fields {
		if f.Predicate != nil && f.Filter == nil {
			return fmt.Errorf("field cannot have a predicate preference without a filter: %+v", f)
		}
	}
	if i.Constraints == nil {
		i.Constraints = &Constraints{
			Fields: fields,
		}
	} else {
		i.Constraints.Fields = fields
	}
	return Validate(i.Constraints)
}

func (i *InputDescriptor) SetSubjectIsIssuer(preference Preference) error {
	if i.Constraints == nil {
		i.Constraints = &Constraints{}
	}
	i.Constraints.SubjectIsIssuer = &preference
	return nil
}

func (i *InputDescriptor) SetSubjectIsHolder(preference Preference) error {
	if i.Constraints == nil {
		i.Constraints = &Constraints{}
	}
	i.Constraints.SubjectIsHolder = &preference
	return nil
}

func (i *InputDescriptor) SetConstraintsLimitDisclosure(limitDisclosure bool) {
	if i.Constraints == nil {
		i.Constraints = &Constraints{
			LimitDisclosure: limitDisclosure,
		}
	}
	i.Constraints.LimitDisclosure = limitDisclosure
}

func NewConstraintsField(path []string) *Field {
	return &Field{Path: path}
}

func (f *Field) SetPurpose(purpose string) {
	f.Purpose = purpose
}

func (f *Field) SetFilter(filter Filter) error {
	if err := Validate(filter); err != nil {
		return err
	}
	f.Filter = &filter
	return nil
}

type PresentationSubmissionBuilder struct {
	Submission PresentationSubmission
}

func NewPresentationSubmissionBuilder(definitionID string) *PresentationSubmissionBuilder {
	uuid, _ := uuid.NewV4()
	return &PresentationSubmissionBuilder{
		Submission: PresentationSubmission{
			ID:           uuid.String(),
			DefinitionID: definitionID,
		},
	}
}

func (p *PresentationSubmissionBuilder) Build() (*PresentationSubmissionHolder, error) {
	return &PresentationSubmissionHolder{
		PresentationSubmission: p.Submission,
	}, Validate(p.Submission)
}

func (p *PresentationSubmissionBuilder) AddDescriptor(d Descriptor) error {
	if err := Validate(d); err != nil {
		return err
	}
	p.Submission.DescriptorMap = append(p.Submission.DescriptorMap, d)
	return nil
}

func (p *PresentationSubmissionBuilder) SetID(id string) {
	p.Submission.ID = id
}

func (p *PresentationSubmissionBuilder) SetDefinitionID(definitionID string) {
	p.Submission.DefinitionID = definitionID
}

func (p *PresentationSubmissionBuilder) SetLocale(locale string) {
	p.Submission.Locale = locale
}
