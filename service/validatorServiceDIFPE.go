package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/gataca-io/vui-core/constant"
	"github.com/gataca-io/vui-core/log"
	"github.com/gataca-io/vui-core/models"
	"github.com/gataca-io/vui-core/tools"

	"github.com/ohler55/ojg/jp"

	converter "github.com/gookit/filter"
)

const (
	CheckIds          = "matchingIds"
	CheckSubject      = "uniqueSubject"
	CheckSubmission   = "submission"
	CheckConstraints  = "constraints"
	CheckRequirements = "submissionRequirements"
	CheckPresentation = "presentationProof"
	CheckCredential   = "credentialProof"
	CheckStatus       = "credentialStatus"
	CheckContext      = "context"
	CheckSchema       = "credentialSchema"
	CheckIssuer       = "issuer"
	CheckIdentity     = "identityVerification"

	thresholdVCStatusCheck = 5 * time.Second
)

type ValidatorServiceDIF struct {
	ssiS SSIService
	jVal JSONValidator
}

//TODO Parallelize processment
func NewDIFValidatorService(ssiService SSIService) Validator {
	return &ValidatorServiceDIF{
		ssiS: ssiService,
		jVal: &jsonValidator{},
	}
}

func (vs *ValidatorServiceDIF) ValidatePresentationResponse(ctx echo.Context, preq models.ExchangeRequest, presp models.ExchangeResponse, token string, signedToken string, requesterVMethod string) (*models.VerificationResult, error) {
	pd := preq.ToPresentationDefinition()
	resp := presp.ToPresentation()
	result := &models.VerificationResult{
		Checks:   []string{},
		Errors:   []string{},
		Warnings: []string{},
	}

	if resp.PresentationSubmission == nil {
		log.CError(ctx, "Missing information about presentation submission")
		return nil, models.ErrInvalidFormat
	}

	submission := resp.PresentationSubmission

	err := vs.validateIds(ctx, result, pd, submission)
	if err != nil {
		return normalizeResult(result), err
	}

	subject, err := vs.getCredentialSubject(ctx, result, resp.VerifiableCredential)
	log.CDebug(ctx, "Received credentials from ", subject)
	if err != nil {
		return normalizeResult(result), err
	}

	err = vs.validateSubmission(ctx, result, pd.InputDescriptors, resp, requesterVMethod)
	if err != nil {
		result.Checks = tools.UniqueSlice(result.Checks)
		return normalizeResult(result), err
	}

	err = vs.validateSubmissionRequirements(ctx, result, pd, submission, resp)
	if err != nil {
		result.Checks = tools.UniqueSlice(result.Checks)
		return normalizeResult(result), err
	}

	err = vs.ssiS.VerifyPresentation(ctx, resp, requesterVMethod)
	if err != nil {
		result.Errors = append(result.Errors, "Verifiable presentation not validated")
		return normalizeResult(result), err
	}
	result.Checks = append(result.Checks, CheckPresentation)

	//TODO: Verify existance of Consent - Outside scope of DIFPE?
	return normalizeResult(result), nil
}

// #########
// * Private
// #########

func (vs *ValidatorServiceDIF) validateIds(ctx echo.Context, result *models.VerificationResult, pd *models.PresentationDefinition, sub *models.PresentationSubmission) error {
	if pd.ID != sub.DefinitionID {
		log.CErrorf(ctx, "Id of presentation submission %s and definition %s not matching", sub.DefinitionID, pd.ID)
		result.Errors = append(result.Errors, "Id of presentation submission and definition are not matching")
		return models.ErrNotMatch
	}
	result.Checks = append(result.Checks, CheckIds)
	return nil
}

func (vs *ValidatorServiceDIF) getCredentialSubject(ctx echo.Context, result *models.VerificationResult, vcs []models.VerifiableCredential) (string, error) {
	user := ""
	for _, vc := range vcs {
		subj := *(vc.CredentialSubject)
		if user == "" {
			user = subj["id"].(string)
		}

		if user != subj["id"] {
			log.CError(ctx, "Credential subject is not the same for all the credentials")
			result.Errors = append(result.Errors, "Credential subject is not the same for all the credentials")
			return "", models.ErrNotMatch
		}
	}
	result.Checks = append(result.Checks, CheckSubject)
	return user, nil
}

func (vs *ValidatorServiceDIF) validateSubmission(ctx echo.Context, result *models.VerificationResult, descriptors []models.InputDescriptor, vp *models.VerifiablePresentation, requesterVMethod string) error {
	vcused := 0
	for _, submitted := range vp.PresentationSubmission.DescriptorMap {
		descriptor := findInputDescriptorWithId(descriptors, submitted.ID)
		if descriptor == nil {
			log.CError(ctx, "Received submission outside of definition: ", submitted.ID)
			result.Errors = append(result.Errors, "Received submission outside of definition")
			return models.ErrMissingClaim
		}
		credData, err := jsonPath(ctx, submitted.Path, vp)
		if err != nil {
			log.CError(ctx, "Cannot discover the reference of the submission")
			result.Errors = append(result.Errors, "Cannot discover the reference of the submission")
			return models.ErrInvalidFormat
		}
		cred := &models.VerifiableCredential{}
		err = tools.ToInterface(credData[0].(map[string]interface{}), cred)
		if err != nil {
			log.CError(ctx, "Cannot discover the reference of the submission")
			result.Errors = append(result.Errors, "Cannot discover the reference of the submission")
			return models.ErrInvalidFormat
		}
		err = vs.validateCredentialWithDescriptor(ctx, result, cred, descriptor, requesterVMethod)
		if err != nil {
			log.CErrorf(ctx, "Submitted credential %s doesn't satisfy descriptor %s constraints", cred.Id, descriptor.ID)
			result.Errors = append(result.Errors, "Submitted credentials don't satisfy descriptor requirements")
			return err
		}
		vcused++
	}
	if vcused < len(vp.VerifiableCredential) {
		log.CErrorf(ctx, "Received more credentials %d than required submissions %d ", len(vp.VerifiableCredential), vcused)
		result.Errors = append(result.Errors, "Received more credentials than required")
		return models.ErrUnwantedClaim
	}
	result.Checks = append(result.Checks, CheckSubmission)
	return nil
}

func (vs *ValidatorServiceDIF) validateCredentialWithDescriptor(ctx echo.Context, result *models.VerificationResult, vc *models.VerifiableCredential, descriptor *models.InputDescriptor, requesterVMethod string) error {

	err := vs.validateSchemas(ctx, result, vc, descriptor.Schema)
	if err != nil {
		log.CError(ctx, "Error validating schemas")
		return err
	}

	err = findIssuerInProofs(ctx, vc, vc.Issuer)
	if err != nil {
		log.CErrorf(ctx, "Asserted issuer %s is not proving the credential %s", vc.Issuer, vc.Id)
		result.Errors = append(result.Errors, "Cannot trust issuer of the credential")
		return err
	}
	result.Checks = append(result.Checks, CheckIssuer)

	_, err = vs.ssiS.VerifyCredential(ctx, vc, requesterVMethod, false)
	if err != nil {
		log.CErrorf(ctx, "Credential %s couldn't be cryptographically validated", vc.Id)
		result.Errors = append(result.Errors, fmt.Sprintf("Credential %s couldn't be cryptographically validated", vc.Id))
		return err
	}
	result.Checks = append(result.Checks, CheckCredential)

	err = vs.verifyStatus(ctx, result, vc)
	if err != nil {
		log.CErrorf(ctx, "Credential %s status couldn't be verified", vc.Id)
		result.Errors = append(result.Errors, fmt.Sprintf("Credential %s status couldn't be verified", vc.Id))
		return err
	}
	result.Checks = append(result.Checks, CheckStatus)

	err = vs.validateCredentialConstraints(ctx, result, vc, descriptor.Constraints)
	if err != nil {
		log.CErrorf(ctx, "Credential %s didn't match required constraints", vc.Id)
		return err
	}
	return nil
}

func (vs *ValidatorServiceDIF) validateSchemas(ctx echo.Context, result *models.VerificationResult, vc *models.VerifiableCredential, requestedSchemas []models.Schema) error {
	found := false
	for _, schema := range requestedSchemas {
		if vc.CredentialSchema != nil {
			if schema.URI == vc.CredentialSchema.Id {
				found = true
				err := vs.jVal.Validate(vc)
				if err != nil {
					log.CErrorf(ctx, "Error validating credential %s with its Json schema", vc.Id)
					return err
				}
				break
			} else if schema.Required {
				log.CErrorf(ctx, "Required schema %s is missing", schema.URI)
				result.Errors = append(result.Errors, "Required schema is missing")
				return models.ErrInvalidFormat
			}
		} else if vc.Context != nil {
			if tools.Contains(vc.Context.GetContext(), schema.URI) {
				found = true
				break
			} else if schema.Required {
				log.CErrorf(ctx, "Required schema %s is missing", schema.URI)
				result.Errors = append(result.Errors, "Required schema is missing")
				return models.ErrInvalidFormat
			}
		} else {
			err := vs.jVal.ValidateWithRef(vc, schema.URI) //No schema given, try to see if matching expected schema
			if err == nil {
				found = true
				result.Warnings = append(result.Errors, "Credential Schema is matching but wasn't explicitely stated")
				break
			}
			if err != nil && schema.Required {
				log.CErrorf(ctx, "Required schema %s is missing", schema.URI)
				result.Errors = append(result.Errors, "Required schema is missing")
				return models.ErrInvalidFormat
			}
		}
	}
	if !found {
		log.CErrorf(ctx, "Credential schema does not match requested schemas")
		result.Errors = append(result.Errors, "Credential schema does not match requested schemas")
		return models.ErrInvalidFormat
	}
	result.Checks = append(result.Checks, CheckSchema)
	return nil
}

func (vs *ValidatorServiceDIF) validateCredentialConstraints(ctx echo.Context, result *models.VerificationResult, vc *models.VerifiableCredential, constraints *models.Constraints) error {
	if constraints == nil {
		result.Warnings = append(result.Warnings, "No constraints required validation")
		return nil
	}
	if constraints.SubjectIsHolder != nil && *constraints.SubjectIsHolder == models.Required {
		//TODO
		result.Warnings = append(result.Warnings, "Subject is holder constraint required but not implemented yet")
		log.Debug("Subject is holder constraint not implemented yet")
	}
	if constraints.SubjectIsIssuer != nil && *constraints.SubjectIsIssuer == models.Required {
		subject := (*vc.CredentialSubject)["id"].(string)
		err := findIssuerInProofs(ctx, vc, subject)
		if err != nil {
			log.CError(ctx, "Subject is not issuer of the credential")
			result.Errors = append(result.Errors, "Subject is not issuer of the credential")
			return err
		}
	}
	if constraints.LimitDisclosure {
		//TODO check if more data than required is being checked
		log.Debug("Limit disclosure constraint not implemented yet")
		result.Warnings = append(result.Warnings, "Limit disclosure constraint required but not implemented yet")
	}

	err := vs.validateFieldConstraint(ctx, vc, constraints.Fields)
	if err != nil {
		log.CError(ctx, "Field constraint not validated")
		result.Errors = append(result.Errors, "Field constraint not validated")
		return err
	}
	result.Checks = append(result.Checks, CheckConstraints)
	return nil
}

func (vs *ValidatorServiceDIF) validateFieldConstraint(ctx echo.Context, vc *models.VerifiableCredential, fieldConstraint []models.Field) error {
	if len(fieldConstraint) == 0 {
		return nil
	}
	mappedCred, err := tools.ToMap(vc)
	if err != nil {
		log.CError(ctx, "Cannot convert vc to map to process it")
		return err
	}
	for _, field := range fieldConstraint {
		//TODO implement filtering
		found := false
		datas := findInPaths(ctx, mappedCred, field.Path)
		for _, data := range datas {
			if data != nil {
				found = true
				err = vs.validateFilter(ctx, data, field.Filter)
				if err != nil {
					log.CError(ctx, "Filtering condition not accepted")
					return err
				}
				break
			}
		}
		if !found && (field.Predicate == nil || *field.Predicate == models.Required) {
			log.CErrorf(ctx, "Missing required paths in object %+v, %+v", field.Path, field.Predicate)
			return models.ErrMissingConstraint
		}
	}
	return nil
}

func (vs *ValidatorServiceDIF) validateFilter(ctx echo.Context, data interface{}, filter *models.Filter) error {
	strData, ok := data.(string)
	if filter.Format == "string" && !ok {
		log.CError(ctx, "Cannot convert expected data to string")
		return models.ErrInvalidFormat
	}
	if filter.Pattern != "" {
		if filter.Type != "string" {
			strData = converter.MustString(data)
		}
		matched, err := regexp.Match(filter.Pattern, []byte(strData))
		if err != nil || !matched {
			log.CError(ctx, "Pattern wasn't fullfilled")
			return models.ErrMissingConstraint
		}
	}
	if len(filter.Enum) > 0 {
		found := false
		for _, v := range filter.Enum {
			if data == v {
				found = true
				break
			}
		}
		if !found {
			log.CError(ctx, "Enum condition wasn't fullfilled")
			return models.ErrMissingConstraint
		}
	}
	if filter.Const != nil {
		if data != filter.Const {
			log.CError(ctx, "Value didn't match constant")
			return models.ErrMissingConstraint
		}
	}
	if filter.MinLength > 0 {
		if len(strData) < filter.MinLength {
			log.CError(ctx, "Value didn't have required min length")
			return models.ErrMissingConstraint
		}
	}
	if filter.MaxLength > 0 {
		if len(strData) > filter.MaxLength {
			log.CError(ctx, "Value didn't have required max length")
			return models.ErrMissingConstraint
		}
	}
	if filter.Not != nil {
		err := vs.validateFilter(ctx, data, filter.Not)
		if err == nil {
			log.CError(ctx, "Negative constraint was successfully evaluated")
			return models.ErrMissingConstraint
		}
	}
	if filter.Maximum != nil {
		err := interfaceComparison(ctx, filter.Format, data, filter.Maximum, ">=")
		if err != nil {
			return err
		}
	}
	if filter.Minimum != nil {
		err := interfaceComparison(ctx, filter.Format, data, filter.Minimum, "<=")
		if err != nil {
			return err
		}
	}
	if filter.ExclusiveMaximum != nil {
		err := interfaceComparison(ctx, filter.Format, data, filter.ExclusiveMaximum, ">")
		if err != nil {
			return err
		}
	}
	if filter.ExclusiveMinimum != nil {
		err := interfaceComparison(ctx, filter.Format, data, filter.ExclusiveMinimum, "<")
		if err != nil {
			return err
		}
	}
	return nil
}

func interfaceComparison(ctx echo.Context, format string, data, comparator interface{}, operation string) error {
	switch format {
	case "date":
		t, err := time.Parse(time.RFC3339, data.(string))
		if err != nil {
			log.CError(ctx, "Cannot convert expected data to date")
			return err
		}
		ct, err := time.Parse(time.RFC3339, comparator.(string))
		if err != nil {
			log.CError(ctx, "Cannot convert expected data to date")
			return err
		}
		switch operation {
		case ">=":
			if ct.Before(t) {
				log.CError(ctx, "Time max or equal constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		case ">":
			if !t.Before(ct) {
				log.CError(ctx, "Time max constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		case "<=":
			if ct.After(t) {
				log.CError(ctx, "Time min or equal constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		case "<":
			if !t.After(ct) {
				log.CError(ctx, "Time min constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		}
	default:
		//Ensure coversions
		nbData, ok := data.(float64)
		if !ok && reflect.TypeOf(data).String() == "string" {
			nbData = converter.MustFloat(data.(string))
		}
		nbComparator, ok := comparator.(float64)
		if !ok && reflect.TypeOf(comparator).String() == "string" {
			nbComparator = converter.MustFloat(comparator.(string))
		}
		switch operation {
		case ">=":
			if nbData < nbComparator {
				log.CError(ctx, "Time max or equal constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		case ">":
			if !(nbData > nbComparator) {
				log.CError(ctx, "Time max constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		case "<=":
			if nbData > nbComparator {
				log.CError(ctx, "Time min or equal constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		case "<":
			if !(nbData < nbComparator) {
				log.CError(ctx, "Time min constraint was successfully evaluated")
				return models.ErrMissingConstraint
			}
		}
	}
	return nil
}

func (vs *ValidatorServiceDIF) validateSubmissionRequirements(ctx echo.Context, result *models.VerificationResult, pd *models.PresentationDefinition, sub *models.PresentationSubmission, resp *models.VerifiablePresentation) error {
	for _, req := range pd.SubmissionRequirements {
		if req.Name == "Identity verification" {
			err := vs.validateSecondFactors(ctx, result, pd, resp)
			if err != nil {
				log.CError(ctx, "Could not validate submission identity verification: ", req.Name)
				return err
			}
			result.Checks = append(result.Checks, CheckIdentity)
		} else {
			err := vs.validateRequirement(ctx, result, &req, pd, sub)
			if err != nil {
				log.CError(ctx, "Could not validate submission requirement: ", req.Name)
				return err
			}
		}
	}
	result.Checks = append(result.Checks, CheckRequirements)
	return nil
}

func (vs *ValidatorServiceDIF) validateRequirement(ctx echo.Context, result *models.VerificationResult, req *models.SubmissionRequirement, pd *models.PresentationDefinition, subs *models.PresentationSubmission) error {
	if len(req.FromNested) > 0 {
		for _, nested := range req.FromNested {
			err := vs.validateRequirement(ctx, result, &nested, pd, subs)
			if err != nil {
				return err
			}
		}
		return nil
	}
	group := req.From
	ids := filterInputDescriptorsIdsInGroup(pd.InputDescriptors, group)
	count := 0
	for _, id := range ids {
		for _, sub := range subs.DescriptorMap {
			if sub.ID == id {
				count++
				break
			}
		}
	}
	switch req.Rule {
	case models.All:
		if count < len(ids) {
			log.CErrorf(ctx, "Required requirement %s wasn't satisfied: only %d out of %d satisfy the condition, not all ", req.Name, count, len(ids))
			result.Errors = append(result.Errors, fmt.Sprintf("Required requirement %s wasn't satisfied: only %d out of %d satisfy the condition, not all ", req.Name, count, len(ids)))
			return models.ErrMissingRequirement
		}
	case models.Pick:
		if req.Minimum != nil {
			if count < *req.Minimum {
				log.CErrorf(ctx, "Required requirement %s wasn't satisfied: only %d, less than %d desired, satisfy the condition", req.Name, count, req.Minimum)
				result.Errors = append(result.Errors, fmt.Sprintf("Required requirement %s wasn't satisfied: only %d, less than %d desired, satisfy the condition", req.Name, count, req.Minimum))
				return models.ErrMissingRequirement
			}
		}
		if req.Maximum != nil {
			if count > *req.Maximum {
				log.CErrorf(ctx, "Required requirement %s wasn't satisfied:  %d, more than of %d maximum, satisfy the condition", req.Name, count, req.Maximum)
				result.Errors = append(result.Errors, fmt.Sprintf("Required requirement %s wasn't satisfied:  %d, more than of %d maximum, satisfy the condition", req.Name, count, req.Maximum))
				return models.ErrMissingRequirement
			}
		}
		if req.Count != nil {
			if count != *req.Count {
				log.CErrorf(ctx, "Required requirement %s wasn't satisfied: only %d instead of %d desired satisfy the condition", req.Name, count, req.Count)
				result.Errors = append(result.Errors, fmt.Sprintf("Required requirement %s wasn't satisfied: only %d instead of %d desired satisfy the condition", req.Name, count, req.Count))
				return models.ErrMissingRequirement
			}
		}
	}
	return nil
}

func (vs *ValidatorServiceDIF) verifyStatus(ctx echo.Context, result *models.VerificationResult, cred *models.VerifiableCredential) error {
	status := cred.CredentialStatus
	if status == nil {
		log.CDebug(ctx, "Credential has no Status")
		result.Warnings = append(result.Warnings, "Credential status not available.")
		return nil
	}

	req, err := http.NewRequest("GET", status.Id, nil)
	if err != nil {
		log.CErrorf(ctx, "Error creating service request to query status %s. Error: %v", status.Id, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(log.HeaderXSpanId, log.GetTraceId(ctx))

	//Setting timeout as context on the request to prevent slow response on offline/block hosts
	client := &http.Client{Timeout: thresholdVCStatusCheck}

	res, err := client.Do(req)
	if err != nil {
		log.CErrorf(ctx, "Error requesting credential status %s. Error: %v", status.Id, err)
		return err
	} else if res.StatusCode != 200 {
		log.CErrorf(ctx, "Error requesting credential status. Status code not success", status.Id)
		return models.ErrStatusNotValid
	}

	vst := models.VerifiableStatus{}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&vst)
	if err != nil {
		log.CErrorf(ctx, "Error decoding requesting credential status %s. Error: %v", status.Id, err)
		return err
	}

	for _, c := range vst.VerifiableCredential {
		if c.Claim.Id == cred.Id {
			if c.Claim.CurrentStatus == constant.CredentialStatusIssued {
				log.CDebugf(ctx, "[Validation Status] Status validated")
				return nil
			}
			log.CDebugf(ctx, "[Validation Status] Status not validated. %s", c.Claim.CurrentStatus)
			return models.ErrStatusNotValid
		}
	}

	log.CWarnf(ctx, "[Validation Status] Credential status scanned but not found")

	return models.ErrStatusNotValid
}

func (vs *ValidatorServiceDIF) validateSecondFactors(ctx echo.Context, result *models.VerificationResult, pd *models.PresentationDefinition, resp *models.VerifiablePresentation) error {
	for _, id := range pd.InputDescriptors {
		if tools.Contains(id.Group, "identity") {
			expectedKey := regexp.MustCompile(id.Constraints.Fields[0].Filter.Pattern)
			found := false
			for _, p := range *resp.GetProofs().GetProof() {
				if len(expectedKey.Find([]byte(p.GetCreator()))) > 0 {
					found = true
					break
				}
			}
			if !found {
				log.CError(ctx, "Required signing key not found in validation", id.Constraints.Fields[0].Filter.Pattern)
				err := models.ErrMissingSecondFactor
				result.Errors = append(result.Errors, err.Error())
				return err
			}
		}
	}
	return nil
}

//Filtering functions
func findIssuerInProofs(ctx echo.Context, vc *models.VerifiableCredential, expectedIssuer string) error {
	proofs := vc.GetProofs()
	if proofs != nil {
		for _, p := range *proofs.GetProof() {
			if strings.Contains(p.GetCreator(), expectedIssuer) {
				return nil
			}
		}
	}
	return models.ErrMissingConstraint
}

func filterInputDescriptorsIdsInGroup(descs []models.InputDescriptor, group string) []string {
	filtered := []string{}
	for _, d := range descs {
		for _, g := range d.Group {
			if g == group {
				filtered = append(filtered, d.ID)
				break
			}
		}
	}
	return filtered
}

func findInputDescriptorWithId(descs []models.InputDescriptor, id string) *models.InputDescriptor {
	for _, d := range descs {
		if d.ID == id {
			return &d
		}
	}
	return nil
}

func findInPaths(ctx echo.Context, object map[string]interface{}, paths []string) []interface{} {
	foundData := []interface{}{}
	for _, pathDef := range paths {
		found, err := jsonPath(ctx, pathDef, object)
		if err != nil {
			log.CWarnf(ctx, "Nothing found in %s: %s", pathDef, err.Error())
		}
		if found != nil {
			foundData = append(foundData, found...)
		}
	}
	return foundData
}

func jsonPath(ctx echo.Context, path string, input interface{}) ([]interface{}, error) {
	mappedData, err := tools.ToMap(input)
	if err != nil {
		log.CError(ctx, "Cannot convert vc to map to process it")
		return nil, err
	}
	x, err := jp.ParseString(path)
	ys := x.Get(mappedData)
	return ys, nil
}

func normalizeResult(result *models.VerificationResult) *models.VerificationResult {
	result.Checks = tools.UniqueSlice(result.Checks)
	result.Warnings = tools.UniqueSlice(result.Warnings)
	return result
}

//Commented do to unuse, but might be useful
// func filterCredentialsWithSchemas(vcs []models.VerifiableCredential, schemas []string) []models.VerifiableCredential {
// 	filtered := []models.VerifiableCredential{}
// 	for _, vc := range vcs {
// 		for _, schema := range schemas {
// 			if vc.CredentialSchema.Id == schema {
// 				filtered = append(filtered, vc)
// 				break
// 			}
// 		}
// 	}
// 	return filtered
// }
