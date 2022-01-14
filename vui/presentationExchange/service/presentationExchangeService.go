package service

import (
	"strings"
	"time"

	"github.com/gataca-io/vui-core/log"
	coreModels "github.com/gataca-io/vui-core/models"
	coreServices "github.com/gataca-io/vui-core/service"
	"github.com/gataca-io/vui-core/tools"
	presentationexchange "github.com/gataca-io/vui-core/vui/presentationExchange"
	"github.com/labstack/echo/v4"
)

type peService struct {
	peRepo           presentationexchange.PresExchangeDao
	daService        presentationexchange.DataAgreementService
	validator        coreServices.Validator
	configRepository presentationexchange.TenantDao
	ssiService       coreServices.SSIService
	trustedWallets   []string
	govS             coreServices.GovernanceService
	limitRequests    bool
}

func NewPresentationExchangeService(presentationExchangeRepo presentationexchange.PresExchangeDao, daService presentationexchange.DataAgreementService, validatorService coreServices.Validator, configRepository presentationexchange.TenantDao, ssiService coreServices.SSIService, govS coreServices.GovernanceService, trustedWallets []string, limitRequests bool) presentationexchange.PresExchangeService {
	return &peService{
		peRepo:           presentationExchangeRepo,
		daService:        daService,
		validator:        validatorService,
		configRepository: configRepository,
		ssiService:       ssiService,
		trustedWallets:   trustedWallets,
		govS:             govS,
		limitRequests:    limitRequests,
	}
}

func (pes *peService) CreateFromTenant(c echo.Context, tenant string, credentialsRequested []string) (*coreModels.PresentationDefinition, error) {
	config, err := pes.configRepository.GetTenantConfig(c, tenant)
	if err != nil {
		log.CDebugf(c, "Could not get tenant config,", err)
		return nil, err
	}

	dataAgreement := config.DataAgreementTemplate
	pes.filterPresentationDefinitionAndDataAgreement(c, credentialsRequested, config, dataAgreement)
	dataAgreement.ID = tools.RandSeq(32)
	dataAgreement.Version = "0"
	dataAgreement.Event = []coreModels.Event{
		{
			PrincipleDid: config.DID,
			State:        "Preparation",
			TimeStamp:    time.Now().UnixMilli(),
			Version:      dataAgreement.Version,
		},
	}
	err = pes.ssiService.SignDataAgreement(c, dataAgreement, config.DID)
	if err != nil {
		log.CError(c, "Cannot sign data agreement", err)
		return nil, err
	}

	difPE, err := pes.createDefinitionFromTenantConfig(c, config)
	if err != nil {
		log.CError(c, "Cannot map input descriptor requirements", err)
		return nil, err
	}
	definition := &coreModels.PresentationDefinition{
		DIFPresentationDefinition: *difPE,
		DataAgreement: &coreModels.DataAgreementRef{
			DataAgreement: dataAgreement,
		},
		Proof: nil,
	}
	err = pes.ssiService.SignPresentationDefinition(c, definition, config.DID)
	if err != nil {
		log.CError(c, "Cannot sign presentation definition", err)
		return nil, err
	}
	_, err = pes.Create(c, definition)
	if err != nil {
		log.CError(c, "Cannot store presentation exchange for validation", err)
		return nil, err
	}
	return definition, nil
}

func (pes *peService) Create(c echo.Context, pe *coreModels.PresentationDefinition) (*coreModels.PExchange, error) {
	t := time.Now()
	pex := &coreModels.PExchange{
		Id:                     pe.ID,
		PresentationDefinition: pe,
		PresentationSubmission: nil,
		CreatedAt:              &t,
		UpdatedAt:              &t,
	}
	err := pes.peRepo.Create(c, pex)
	if err != nil {
		log.CError(c, "Presentation exchange couldn't be created", err)
		return nil, err
	}
	return pex, nil
}

func (pes *peService) GetDefinition(c echo.Context, id string, dataAgreementOnly bool) (*coreModels.PresentationDefinition, error) {
	pe, err := pes.GetExchange(c, id)
	if err != nil {
		return nil, err
	}
	if pe.RequestedAt != nil {
		err := coreModels.ErrSessionRequested
		log.CError(c, "This session is already ongoing", pe.Id, err)
		return nil, err
	}

	if pes.limitRequests {
		t := time.Now()
		pe.RequestedAt = &t

		err = pes.peRepo.Update(c, pe)
		if err != nil {
			log.CError(c, "Cannot update presentation exchange in db", err)
			return nil, err
		}
	}

	return pe.PresentationDefinition, nil
}

func (pes *peService) Submit(c echo.Context, id string, verifiablePresentation *coreModels.VerifiablePresentation) (*coreModels.VerificationResult, error) {
	pe, err := pes.GetExchange(c, id)
	if err != nil {
		return nil, err
	}
	pe.PresentationSubmission = verifiablePresentation
	verificationResult, err := pes.verify(c, pe)

	t := time.Now()
	pe.Validations = verificationResult
	pe.UpdatedAt = &t

	err2 := pes.peRepo.Update(c, pe)
	if err2 != nil {
		log.CError(c, "Cannot update presentation exchange in db", err2)
		return nil, err2
	}

	return verificationResult, err
}

func (pes *peService) GetVerification(c echo.Context, id string) (*coreModels.VerificationResult, error) {
	pe, err := pes.GetExchange(c, id)
	if err != nil {
		return nil, err
	}
	return pe.Validations, nil
}

func (pes *peService) GetExchange(c echo.Context, id string) (*coreModels.PExchange, error) {
	pe, err := pes.peRepo.GetByID(c, id)
	if err != nil {
		log.CErrorf(c, "Presentation exchange with id %s not found: %s", id, err.Error())
		return nil, err
	}
	return pe, nil
}

func (pes *peService) GetSubmittedData(c echo.Context, id string) (map[string]interface{}, error) {
	mergedData := map[string]interface{}{}
	return mergedData, nil
}

func (pes *peService) Delete(c echo.Context, id string) error {
	err := pes.peRepo.Delete(id)
	if err != nil {
		log.CError(c, "Unable to delete presentation exchange with id ", id, err)
	}
	return err
}

// ############
// ## PRIVATE
// ############

func (pes *peService) verify(c echo.Context, pe *coreModels.PExchange) (*coreModels.VerificationResult, error) {
	verificationResult, err := pes.validator.ValidatePresentationResponse(c, pe.PresentationDefinition, pe.PresentationSubmission, pe.PresentationDefinition.DataAgreement.DataAgreement.DataReceiver.ID)
	if err != nil {
		return verificationResult, err
	}
	verificationResult.Checks = append(verificationResult.Checks, "consent")
	presentationCreators := pe.PresentationSubmission.GetProofs().GetCreators()
	dataSubject := (*pe.PresentationSubmission.VerifiableCredential[0].CredentialSubject)["id"].(string)
	if pe.PresentationSubmission.DataAgreementId == "" {
		log.CWarn(c, "Presentation Submission not enforcing data agreement, it is not from Gataca. Check if it is signed by the user")
		if !presentInArray(presentationCreators, dataSubject) {
			log.CError(c, "Cannot trust source of the credential")
			errConsent := coreModels.ErrConsentValidation
			verificationResult.Errors = append(verificationResult.Errors, errConsent.Error())
			return verificationResult, coreModels.ErrConsentValidation
		}
	}
	dataAgreement, err := pes.daService.GetDataAgreement(c, pe.PresentationSubmission.DataAgreementId, -1)
	if err != nil {
		log.CError(c, "Cannot retrieved associated data agreement")
		errConsent := coreModels.ErrConsentValidation
		verificationResult.Errors = append(verificationResult.Errors, errConsent.Error())
		return verificationResult, coreModels.ErrConsentValidation
	}
	if dataAgreement.DataSubject != dataSubject {
		log.CError(c, "Cannot data agreement doesn't allow the usage of this subject")
		errConsent := coreModels.ErrConsentValidation
		verificationResult.Errors = append(verificationResult.Errors, errConsent.Error())
		return verificationResult, coreModels.ErrConsentValidation
	}
	if !presentInArray(presentationCreators, dataAgreement.DataHolder) {
		log.CError(c, "Cannot data agreement doesn't allow to trust the holder of these credentials")
		errConsent := coreModels.ErrConsentValidation
		verificationResult.Errors = append(verificationResult.Errors, errConsent.Error())
		return verificationResult, coreModels.ErrConsentValidation
	}
	for _, cred := range pe.PresentationSubmission.VerifiableCredential {
		found := false
		for _, pid := range dataAgreement.PersonalData {
			if pid.AttributeID == cred.Id {
				found = true
				break
			}
		}
		if !found {
			log.CErrorf(c, "Cannot data agreement doesn't allow to use this credential %s", cred.Id)
			errConsent := coreModels.ErrConsentValidation
			verificationResult.Errors = append(verificationResult.Errors, errConsent.Error())
			return verificationResult, coreModels.ErrConsentValidation
		}
	}
	return verificationResult, nil
}

func (pes *peService) filterPresentationDefinitionAndDataAgreement(c echo.Context, credentialsRequested []string, config *coreModels.TenantConfig, dataAgreement *coreModels.DataAgreement) {
	if len(credentialsRequested) > 0 {
		if config.AdvancedDefinition != nil {
			pes.filterPresentationDefinitionAndDataAgreementFromSet(c, credentialsRequested, config.AdvancedDefinition, dataAgreement)
		} else {
			pes.filterPresentationDefinitionAndDataAgreementFromConfig(c, credentialsRequested, config, dataAgreement)
		}
	}
}

func (pes *peService) filterPresentationDefinitionAndDataAgreementFromSet(c echo.Context, credentialsRequested []string, definition *coreModels.PresentationDefinition, dataAgreement *coreModels.DataAgreement) {
	submissionSet := map[string]bool{}
	filteredSubs := []coreModels.SubmissionRequirement{}
	filteredInputDescriptors := []coreModels.InputDescriptor{}
	purposesSet := map[string]bool{}
	filteredPurposes := []coreModels.Purpose{}
	filteredData := []coreModels.PersonalDatum{}
	for _, requested := range credentialsRequested {
		for _, idesc := range definition.InputDescriptors {
			if idesc.ID == requested {
				filteredInputDescriptors = append(filteredInputDescriptors, idesc)
				for _, sub := range definition.SubmissionRequirements {
					if tools.Contains(idesc.Group, sub.From) && !submissionSet[sub.From] {
						submissionSet[sub.From] = true
						filteredSubs = append(filteredSubs, sub)
					}
				}
				for _, pDatum := range dataAgreement.PersonalData {
					if pDatum.AttributeName == idesc.ID {
						filteredData = append(filteredData, pDatum)
						for _, purpose := range dataAgreement.Purposes {
							if tools.Contains(pDatum.Purposes, purpose.ID) && !purposesSet[purpose.ID] {
								purposesSet[purpose.ID] = true
								filteredPurposes = append(filteredPurposes, purpose)
							}
						}
					}
				}
			}
		}
	}
	definition.InputDescriptors = filteredInputDescriptors
	definition.SubmissionRequirements = filteredSubs
	dataAgreement.PersonalData = filteredData
	dataAgreement.Purposes = filteredPurposes
}

func (pes *peService) filterPresentationDefinitionAndDataAgreementFromConfig(c echo.Context, credentialsRequested []string, config *coreModels.TenantConfig, dataAgreement *coreModels.DataAgreement) {
	newPurposes := []coreModels.Purpose{}
	purposesSet := map[string]bool{}
	newPersonalData := []coreModels.PersonalDatum{}
	newCreds := []coreModels.CredentialRequest{}
	for _, cred := range credentialsRequested {
		for _, credConfig := range config.Credentials {
			if cred == credConfig.Type {
				newCreds = append(newCreds, credConfig)
				for _, pd := range dataAgreement.PersonalData {
					if tools.Contains(pd.Purposes, credConfig.Purpose) {
						newPersonalData = append(newPersonalData, pd)
						if !purposesSet[credConfig.Purpose] {
							for _, p := range dataAgreement.Purposes {
								if p.ID == credConfig.Purpose {
									newPurposes = append(newPurposes, p)
									purposesSet[credConfig.Purpose] = true
								}
							}
						}
					}
				}
			}
		}
	}
	dataAgreement.PersonalData = newPersonalData
	dataAgreement.Purposes = newPurposes
	config.Credentials = newCreds
}

func (pes *peService) createDefinitionFromTenantConfig(c echo.Context, config *coreModels.TenantConfig) (*coreModels.DIFPresentationDefinition, error) {
	if config.AdvancedDefinition != nil {
		if config.AdvancedDefinition.ID == "" {
			config.AdvancedDefinition.ID = tools.RandSeq(32)
		}
		if config.AdvancedDefinition.Name == "" {
			config.AdvancedDefinition.Name = config.TenantId
		}
		if config.AdvancedDefinition.Purpose == "" {
			config.AdvancedDefinition.Purpose = config.ServicePurpose
		}
		if config.AdvancedDefinition.Format == nil {
			config.AdvancedDefinition.Format = &coreModels.Format{
				LDPVP: &coreModels.LDPType{
					ProofType: []string{
						"JsonWebSignature2020",
						"Ed25519Signature2018",
						"EcdsaSecp256k1Signature2019",
						"RsaSignature2018",
					},
				},
				LDP: &coreModels.LDPType{
					ProofType: []string{
						"JsonWebSignature2020",
						"Ed25519Signature2018",
						"EcdsaSecp256k1Signature2019",
						"RsaSignature2018",
					},
				},
			}
		}
		inputDescriptors, subRequirements := pes.addSecurityConfig(c, config, config.AdvancedDefinition.InputDescriptors, config.AdvancedDefinition.SubmissionRequirements)
		config.AdvancedDefinition.InputDescriptors = inputDescriptors
		config.AdvancedDefinition.SubmissionRequirements = subRequirements
		return &config.AdvancedDefinition.DIFPresentationDefinition, nil
	}
	min := 0
	subRequirements := []coreModels.SubmissionRequirement{
		{
			Name:    "Mandatory data",
			Purpose: "Basic data to provide the service",
			Rule:    "all",
			FromOption: coreModels.FromOption{
				From: "mandatory",
			},
		},
		{
			Name:    "Optional data",
			Purpose: "Additional data to enrich the service",
			Rule:    "pick",
			Minimum: &min,
			FromOption: coreModels.FromOption{
				From: "optional",
			},
		},
	}
	inputDescriptors := []coreModels.InputDescriptor{}
	for _, cred := range config.Credentials {
		id, err := pes.buildCredentialInputDescriptor(c, cred, config.DID)
		if err != nil {
			log.CError(c, "Cannot build input descriptors for credential config", cred.Type, err)
			return nil, err
		}
		inputDescriptors = append(inputDescriptors, id)
	}
	inputDescriptors, subRequirements = pes.addSecurityConfig(c, config, inputDescriptors, subRequirements)
	def := &coreModels.DIFPresentationDefinition{
		ID:                     tools.RandSeq(32),
		Name:                   config.TenantId,
		Purpose:                config.ServicePurpose,
		SubmissionRequirements: subRequirements,
		InputDescriptors:       inputDescriptors,
		Format: &coreModels.Format{
			LDPVP: &coreModels.LDPType{
				ProofType: []string{
					"JsonWebSignature2020",
					"Ed25519Signature2018",
					"EcdsaSecp256k1Signature2019",
					"RsaSignature2018",
				},
			},
			LDP: &coreModels.LDPType{
				ProofType: []string{
					"JsonWebSignature2020",
					"Ed25519Signature2018",
					"EcdsaSecp256k1Signature2019",
					"RsaSignature2018",
				},
			},
		},
	}
	return def, nil
}

func (pes *peService) addSecurityConfig(c echo.Context, config *coreModels.TenantConfig, inputDescriptors []coreModels.InputDescriptor, subRequirements []coreModels.SubmissionRequirement) ([]coreModels.InputDescriptor, []coreModels.SubmissionRequirement) {
	if len(config.Security) > 0 {
		sr := coreModels.SubmissionRequirement{
			Name:    "Identity verification",
			Purpose: "Additional mechanisms of identity verification",
			Rule:    "all",
			FromOption: coreModels.FromOption{
				From: "identity",
			},
		}
		subRequirements = append(subRequirements, sr)

		authFactors := []string{}
		for _, sec := range config.Security {
			switch sec.Type {
			case coreModels.AppAuth:
				trustedWallets := strings.Join(pes.trustedWallets, "|")
				id := pes.buildGenericInputDescriptor("wallet_auth", "Wallet Application Authentication", "We need to assert the source of the credentials as a trusted wallet", "identity", "vp.proof.creator", "The presentation must be signed with one of the trusted app dids", "^("+trustedWallets+")#appAuth$")
				inputDescriptors = append(inputDescriptors, id)
			case coreModels.AuthNFactor:
				authFactors = append(authFactors, sec.Accepted...)
			}
		}
		if len(authFactors) > 0 {
			authFactorsParam := strings.Join(authFactors, "|")
			id := pes.buildGenericInputDescriptor("authnfactors", "Identity verification enforced", "We need to assert the holder's identity matching the subject", "identity", "vp.proof.creator", "A second factor must be enforced", "^did:(.*)#("+authFactorsParam+")$")
			inputDescriptors = append(inputDescriptors, id)
		}
	}
	return inputDescriptors, subRequirements
}

func (pes *peService) buildCredentialInputDescriptor(ctx echo.Context, cred coreModels.CredentialRequest, configDID string) (coreModels.InputDescriptor, error) {
	var id coreModels.InputDescriptor
	if cred.Mandatory {
		id = pes.buildGenericInputDescriptor(cred.Type, cred.Type, cred.Purpose, "mandatory", "$.type[-1:]", "The claim must be of a specific type", cred.Type)
	} else {
		id = pes.buildGenericInputDescriptor(cred.Type, cred.Type, cred.Purpose, "optional", "$.type[-1:]", "The claim must be of a specific type", cred.Type)
	}
	err := pes.getIssuersForTrustedLevel(ctx, cred, &id, configDID)
	if err != nil {
		log.CErrorf(ctx, "Error finding issuer requirements", err)
		return id, err
	}
	return id, nil
}

func (pes *peService) getIssuersForTrustedLevel(ctx echo.Context, cred coreModels.CredentialRequest, id *coreModels.InputDescriptor, configDID string) error {
	if cred.TrustLevel == 0 {
		return nil
	}
	mappedSchemas := []string{}
	for _, schema := range id.Schema {
		mappedSchemas = append(mappedSchemas, schema.URI)
	}
	trustedIssuers, err := pes.govS.GetTrustedIssuersForSchemas(ctx, cred.Type, mappedSchemas)
	if err != nil {
		log.CError(ctx, "Error retrieving catalogs", err)
		return err
	}
	issuers := []string{}
	for _, ti := range trustedIssuers {
		issuers = append(issuers, ti.Dids...)
	}
	issPattern := "(" + strings.Join(issuers, "|") + ")"
	issuerConstraint := coreModels.Field{
		Path: []string{"$.issuer",
			"$.vc.issuer",
			"$.iss",
		},
		Purpose: "Assert the trust of the issuer",
		Filter: &coreModels.Filter{
			Type:    "string",
			Pattern: issPattern,
		},
	}
	id.Constraints.Fields = append(id.Constraints.Fields, issuerConstraint)
	return nil
}

func (pes *peService) buildGenericInputDescriptor(id string, name string, purpose string, group string, path string, pathpurpose string, pattern string) coreModels.InputDescriptor {
	return coreModels.InputDescriptor{
		ID:      id,
		Name:    name,
		Purpose: purpose,
		Schema: []coreModels.Schema{{
			URI: "https://www.w3.org/2018/credentials/v1",
		},
		},
		Group: []string{group},
		Constraints: &coreModels.Constraints{
			Fields: []coreModels.Field{
				{
					Path:    []string{path},
					Purpose: pathpurpose,
					Filter: &coreModels.Filter{
						Type:    "string",
						Pattern: pattern,
					},
				},
			},
		},
	}
}

func presentInArray(arr []string, s string) bool {
	for _, arrS := range arr {
		if strings.Contains(arrS, s) {
			return true
		}
	}
	return false
}
