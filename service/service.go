package service

import (
	"github.com/gataca-io/vui-core/models"
	"github.com/labstack/echo/v4"
)

type SSIService interface {
	ValidateLdContext(ctx echo.Context, ldv models.LdContext) (string, error)

	SignCredential(ctx echo.Context, vc *models.VerifiableCredential, vm string, proofType string) error
	SignQualifiedCredential(ctx echo.Context, vc *models.VerifiableCredential, vm string) error
	VerifyCredential(ctx echo.Context, vc *models.VerifiableCredential, requester string, sbx bool) (int, error)

	SignPresentation(ctx echo.Context, vc *models.VerifiablePresentation, vmethod string) error
	VerifyPresentation(ctx echo.Context, fc *models.VerifiablePresentation, requester string) error

	SignDataAgreement(ctx echo.Context, da *models.DataAgreement, vmethod string) error
	VerifyDataAgreement(ctx echo.Context, da *models.DataAgreement, requester string) error

	SignPresentationDefinition(ctx echo.Context, pd *models.PresentationDefinition, vmethod string) error
	VerifyPresentationDefinition(ctx echo.Context, pd *models.PresentationDefinition, requester string) error

	VerifyDIDDocument(ctx echo.Context, fc *models.DIDDocument, vmethods []*models.PublicKey) error
	SignDIDDocument(ctx echo.Context, fc *models.DIDDocument, vmethod string) error
}

type JSONValidator interface {
	Validate(document models.JSONSchema) error
	ValidateWithRef(document models.JSONSchema, ref string) error
	ValidateStrings(schema, document string) error
}

type Validator interface {
	ValidatePresentationResponse(ctx echo.Context, pr models.ExchangeRequest, resp models.ExchangeResponse, token string, signedToken string, requesterVMethod string) (*models.VerificationResult, error)
}

type DidService interface {
	GetDID(ctx echo.Context, did string) (*models.DIDDocument, error)
	CreateDID(ctx echo.Context, did *models.DIDDocument) error
	UpdateDID(ctx echo.Context, did *models.DIDDocument) error
	RevokeDID(ctx echo.Context, did *models.DIDDocument) error
}

type GovernanceService interface {
	GetTrustedIssuersForSchemas(ctx echo.Context, credentialType string, schemasOrContexts []string) ([]models.TrustedIssuer, error)
}
