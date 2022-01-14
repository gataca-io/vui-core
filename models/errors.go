package models

import "errors"

// ResponseMessage represent the response error struct
type ResponseMessage struct {
	Message string `json:"message" example:"Some description message" description:"Operation result message description"`
}

var (
	//Common or admin
	ErrBadParamInput       = errors.New("request is not valid")
	ErrNotFound            = errors.New("your requested item is not found")
	ErrNotInitialized      = errors.New("application not initialized")
	ErrAlreadyInitialized  = errors.New("application already initialized")
	ErrNotConfigured       = errors.New("configuration available but not applied")
	ErrInternalServerError = errors.New("internal Server Error")
	ErrConflict            = errors.New("your item already exists")

	//Connect
	ErrEmptySession   = errors.New("your session is empty")
	ErrInvalidSession = errors.New("your session is invalid or cannot be queried")

	// DID
	ErrDIDNotAvailable     = errors.New("DID is not available")
	ErrCatalogNotAvailable = errors.New("catalog is not available")
	ErrMissingKey          = errors.New("verification method not present in did")
	ErrInvalidDIDMethod    = errors.New("cannot register DIDs with non GATC method")

	//Validations
	ErrNotMatch            = errors.New("presentation response does not match")
	ErrRepeatedClaim       = errors.New("required claim is found repeated")
	ErrInvalidFormat       = errors.New("document is not valid json")
	ErrUnwantedClaim       = errors.New("found unrequested claim")
	ErrMissingClaim        = errors.New("required claim is missing")
	ErrMissingConstraint   = errors.New("required constraint couldn't be satisfied")
	ErrMissingRequirement  = errors.New("submission requirement couldn't be satisfied")
	ErrMissingSecondFactor = errors.New("required second factor is missing")
	ErrMissingSFProof      = errors.New("second factor proof is missing")
	ErrSFValidation        = errors.New("second factor could not be validated")
	ErrInvalidContext      = errors.New("document linked data context couldn't be validated")
	ErrMissingVerifiable   = errors.New("missing object to verify")
	ErrSessionRequested    = errors.New("your session is already in use")
	ErrConsentValidation   = errors.New("claim consent validation fail")
	ErrCredentialsNotMatch = errors.New("credentials requested are not avaliable in tenant")
	ErrRenewDisallowed     = errors.New("renew service is not activated")

	//Status
	ErrStatusNotValid = errors.New("credential status not valid")
)
