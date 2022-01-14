package service

import (
	"encoding/json"
	"testing"

	"github.com/gataca-io/vui-core/log"
	"github.com/gataca-io/vui-core/models"
	"github.com/gataca-io/vui-core/testdata"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var difValidator *ValidatorServiceDIF
var mockedSSIs SSIService
var mockedJVal JSONValidator

type mockSSIService struct{}

type mockJSONValidator struct{}

func (ms *mockSSIService) ValidateLdContext(ctx echo.Context, ldv models.LdContext) (string, error) {
	return "", nil
}
func (ms *mockSSIService) SignCredential(ctx echo.Context, vc *models.VerifiableCredential, vm, proofType string) error {
	return nil
}
func (ms *mockSSIService) SignQualifiedCredential(ctx echo.Context, vc *models.VerifiableCredential, vm string) error {
	return nil
}
func (ms *mockSSIService) VerifyCredential(ctx echo.Context, vc *models.VerifiableCredential, requester string, sbx bool) (int, error) {
	return 0, nil
}
func (ms *mockSSIService) SignPresentation(ctx echo.Context, vc *models.VerifiablePresentation, did string) error {
	return nil
}
func (ms *mockSSIService) VerifyPresentation(ctx echo.Context, fc *models.VerifiablePresentation, requester string) error {
	return nil
}
func (ms *mockSSIService) VerifyDIDDocument(ctx echo.Context, fc *models.DIDDocument, vmethods []*models.PublicKey) error {
	return nil
}
func (ms *mockSSIService) SignDIDDocument(ctx echo.Context, fc *models.DIDDocument, vmethod string) error {
	return nil
}
func (ms *mockSSIService) SignDataAgreement(ctx echo.Context, da *models.DataAgreement, vmethod string) error {
	return nil
}
func (ms *mockSSIService) VerifyDataAgreement(ctx echo.Context, da *models.DataAgreement, requester string) error {
	return nil
}
func (ms *mockSSIService) SignPresentationDefinition(ctx echo.Context, pd *models.PresentationDefinition, vmethod string) error {
	return nil
}
func (ms *mockSSIService) VerifyPresentationDefinition(ctx echo.Context, pd *models.PresentationDefinition, requester string) error {
	return nil
}

func (mj *mockJSONValidator) Validate(document models.JSONSchema) error {
	return nil
}
func (mj *mockJSONValidator) ValidateWithRef(document models.JSONSchema, ref string) error {
	return nil
}
func (mj *mockJSONValidator) ValidateStrings(schema, document string) error {
	return nil
}

func init() {
	mockedSSIs = &mockSSIService{}
	mockedJVal = &mockJSONValidator{}
	difValidator = &ValidatorServiceDIF{
		ssiS: mockedSSIs,
		jVal: mockedJVal,
	}
}

//PRIVATE DATA Models

func createPresentationDefinition(t *testing.T) *models.PresentationDefinition {
	var presDef models.PresentationDefinition
	presDefBytes := []byte(testdata.MultiGroupPresentationDefinition)
	err := json.Unmarshal(presDefBytes, &presDef)
	assert.NoError(t, err)
	return &presDef
}

func createVerifiablePresentation(t *testing.T) *models.VerifiablePresentation {
	var pres models.VerifiablePresentation
	presBytes := []byte(testdata.SampleVerifiablePresentation)
	err := json.Unmarshal(presBytes, &pres)
	assert.NoError(t, err)
	return &pres
}

func createEmptyVerificationResult() *models.VerificationResult {
	return &models.VerificationResult{
		Checks:   []string{},
		Errors:   []string{},
		Warnings: []string{},
	}
}

func TestDIFValidatorService_ValidateSameSubject(t *testing.T) {
	vp := createVerifiablePresentation(t)
	res := createEmptyVerificationResult()

	subj, err := difValidator.getCredentialSubject(nil, res, vp.VerifiableCredential)
	assert.NoError(t, err)
	assert.Equal(t, "did:example:ebfeb1f712ebc6f1c276e12ec21", subj)
	assert.NotEmpty(t, res.Checks)
	assert.Equal(t, 1, len(res.Checks))
	assert.Empty(t, res.Errors)
}

func TestDIFValidatorService_ValidateDifferentSubject(t *testing.T) {
	vp := createVerifiablePresentation(t)
	(*vp.VerifiableCredential[0].CredentialSubject)["id"] = "did:example:212444"
	res := createEmptyVerificationResult()

	subj, err := difValidator.getCredentialSubject(nil, res, vp.VerifiableCredential)
	assert.Error(t, err)
	assert.Equal(t, "", subj)
	assert.Empty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 1, len(res.Errors))
}

func TestDIFValidatorService_ValidateIds(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	res := createEmptyVerificationResult()

	err := difValidator.validateIds(nil, res, presentationDefinition, vp.PresentationSubmission)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Equal(t, 1, len(res.Checks))
	assert.Empty(t, res.Errors)
}
func TestDIFValidatorService_ValidateIdsNonMatching(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	vp.PresentationSubmission.DefinitionID = vp.PresentationSubmission.ID
	res := createEmptyVerificationResult()

	err := difValidator.validateIds(nil, res, presentationDefinition, vp.PresentationSubmission)
	assert.Error(t, err)
	assert.Empty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 1, len(res.Errors))
}

func TestDIFValidatorService_ValidateSubmission(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	res := createEmptyVerificationResult()

	err := difValidator.validateSubmission(nil, res, presentationDefinition.InputDescriptors, vp, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Equal(t, 16, len(res.Checks))
	assert.Empty(t, res.Errors)
}

func TestDIFValidatorService_ValidateNoSubmission(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	vp.PresentationSubmission = nil

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrInvalidFormat, err)
	assert.Nil(t, res)
}

func TestDIFValidatorService_ValidateIssuerSigning(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	log.Debugf("%+v", vp.VerifiableCredential[0].Proof)
	vp.VerifiableCredential[0].Proof.Value.VerificationMethod = "did:example:987#keys-1"

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 3, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))

}

func TestDIFValidatorService_ValidateIssuerPatternConstraint(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	vp.VerifiableCredential[0].Issuer = "did:example:987"

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 3, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidateAccountIdMinLengthConstraint(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	subject := *vp.VerifiableCredential[0].CredentialSubject
	accounts := subject["account"].([]interface{})
	account1 := accounts[0].(map[string]interface{})
	account1["id"] = "123456789" //Min 10

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 6, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidateAccountIdMaxLengthConstraint(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	subject := *vp.VerifiableCredential[0].CredentialSubject
	accounts := subject["account"].([]interface{})
	account1 := accounts[0].(map[string]interface{})
	account1["id"] = "1234567890123" //Max 12

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 6, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidateAccountIdPatternConstraintValid(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	subject := *vp.VerifiableCredential[0].CredentialSubject
	accounts := subject["account"].([]interface{})
	accounts = accounts[0:1] //remove last
	account1 := accounts[0].(map[string]interface{})
	account1["route"] = "JP-1234567890" //Max 12

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Empty(t, res.Errors)
	assert.Equal(t, 10, len(res.Checks))
}

func TestDIFValidatorService_ValidateAccountIdPatternConstraintInvalid(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	subject := *vp.VerifiableCredential[0].CredentialSubject
	accounts := subject["account"].([]interface{})
	account1 := accounts[0].(map[string]interface{})
	account1["route"] = "ES-1234567890" //Must be DE or US or JP

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 6, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidatePreferred(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	req := models.Required
	presentationDefinition.InputDescriptors[2].Constraints.Fields[0].Predicate = &req //That predicate would be failing

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 7, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidateAccountBooleanPattern(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	req := models.Required
	presentationDefinition.InputDescriptors[2].Constraints.Fields[0].Predicate = &req
	subject := *vp.VerifiableCredential[1].CredentialSubject
	jobs := subject["jobs"].([]interface{})
	jobs1 := jobs[0].(map[string]interface{})
	jobs1["active"] = true //Max 12

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Empty(t, res.Errors)
	assert.Equal(t, 10, len(res.Checks))
}

func TestDIFValidatorService_ValidateMissingField(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	paths := presentationDefinition.InputDescriptors[3].Constraints.Fields[1].Path
	paths = append(paths[0:1], paths[2:]...)
	presentationDefinition.InputDescriptors[3].Constraints.Fields[1].Path = paths

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 7, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidateDateConstraint(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	filter := presentationDefinition.InputDescriptors[3].Constraints.Fields[1].Filter
	filter.Minimum = filter.Maximum
	filter.Maximum = nil

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingConstraint, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 7, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidateNoSchema(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	vp.VerifiableCredential[0].CredentialSchema = nil
	vp.VerifiableCredential[0].Context = nil
	//Should pass even if not stated explicitely. The mock json validator assumes it is a valid schema

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Empty(t, res.Errors)
	assert.Equal(t, 10, len(res.Checks))
}

func TestDIFValidatorService_ValidateContextSchema(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	vp.VerifiableCredential[0].CredentialSchema = nil
	vp.VerifiableCredential[0].Context = &models.SSIContext{
		Context: "https://bank-schemas.org/1.0.0/accounts.json",
	}
	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Empty(t, res.Errors)
	assert.Equal(t, 10, len(res.Checks))
}

func TestDIFValidatorService_ValidateWrongSchema(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	vp.VerifiableCredential[0].CredentialSchema = &models.CredentialSchema{
		Id: "https://schema.org/nonexisting",
	}

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrInvalidFormat, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 2, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidatePossibleSchema(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	vp.VerifiableCredential[0].CredentialSchema = &models.CredentialSchema{
		Id: "https://bank-schemas.org/1.0.0/accounts.json",
	}

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Empty(t, res.Errors)
	assert.Equal(t, 10, len(res.Checks))
}

func TestDIFValidatorService_ValidateRequiredSchema(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	presentationDefinition.InputDescriptors[1].Schema[0].Required = true //The credential hasn't got the required but the optional

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrInvalidFormat, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 2, len(res.Checks))
	assert.Equal(t, 2, len(res.Errors))
}

func TestDIFValidatorService_ValidateSubmissionRequirementsCount(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	count := 2
	presentationDefinition.SubmissionRequirements[0].Count = &count

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingRequirement, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 8, len(res.Checks))
	assert.Equal(t, 1, len(res.Errors))
}

func TestDIFValidatorService_ValidateSubmissionRequirementsMin(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	max := 2
	presentationDefinition.SubmissionRequirements[0].Minimum = &max

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingRequirement, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 8, len(res.Checks))
	assert.Equal(t, 1, len(res.Errors))
}

func TestDIFValidatorService_ValidateSubmissionRequirementsAll(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)
	presentationDefinition.SubmissionRequirements[0].Rule = models.All

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.Error(t, err)
	assert.Equal(t, models.ErrMissingRequirement, err)
	assert.NotEmpty(t, res.Checks)
	assert.NotEmpty(t, res.Errors)
	assert.Equal(t, 8, len(res.Checks))
	assert.Equal(t, 1, len(res.Errors))
}

func TestDIFValidatorService_ValidateComplete(t *testing.T) {
	presentationDefinition := createPresentationDefinition(t)
	vp := createVerifiablePresentation(t)

	res, err := difValidator.ValidatePresentationResponse(nil, presentationDefinition, vp, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Checks)
	assert.Empty(t, res.Errors)
	assert.Equal(t, 10, len(res.Checks))
}
