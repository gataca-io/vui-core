package presentationexchange

import (
	coreModels "github.com/gataca-io/vui-core/models"
	"github.com/labstack/echo/v4"
)

type PresExchangeService interface {
	CreateFromTenant(c echo.Context, tenant string, credentialsRequested []string) (*coreModels.PresentationDefinition, error)

	Create(c echo.Context, pe *coreModels.PresentationDefinition) (*coreModels.PExchange, error)
	Submit(c echo.Context, id string, verifiablePresentation *coreModels.VerifiablePresentation) (*coreModels.VerificationResult, error)
	GetExchange(c echo.Context, id string) (*coreModels.PExchange, error)
	GetDefinition(c echo.Context, id string, dataAgreementOnly bool) (*coreModels.PresentationDefinition, error)
	GetVerification(c echo.Context, id string) (*coreModels.VerificationResult, error)
	GetSubmittedData(c echo.Context, id string) (map[string]interface{}, error)
	Delete(c echo.Context, id string) error
}

type DataAgreementService interface {
	Create(c echo.Context, da *coreModels.DataAgreement) (*coreModels.DataAgreement, error)
	GetDataAgreement(c echo.Context, id string, version int) (*coreModels.DataAgreement, error)
	Update(c echo.Context, da *coreModels.DataAgreement) (*coreModels.DataAgreement, error)
	Delete(c echo.Context, da *coreModels.DataAgreement) (*coreModels.DataAgreement, error)
}

type PresExchangeDao interface {
	Create(c echo.Context, pe *coreModels.PExchange) error
	Update(c echo.Context, pe *coreModels.PExchange) error
	DeleteLogical(id string) error
	Delete(id string) error
	GetByID(c echo.Context, id string) (*coreModels.PExchange, error)
}

type DataAgreementDao interface {
	Create(c echo.Context, da *coreModels.DataAgreement) error
	GetDataAgreement(c echo.Context, id string) (*coreModels.DataAgreement, error)
	UpdateDataAgreement(c echo.Context, da *coreModels.DataAgreement) error
	Delete(c echo.Context, id string) error
	DeleteLogical(c echo.Context, id string) error
}

type TenantDao interface {
	GetTenantConfig(c echo.Context, tenant string) (*coreModels.TenantConfig, error)
	GetTenantConfigs(c echo.Context, tenants []string) ([]coreModels.TenantConfig, error)
	GetAllConfigs(c echo.Context) ([]coreModels.TenantConfig, error)
	CreateConfig(c echo.Context, config *coreModels.TenantConfig) error
	UpdateConfig(c echo.Context, config *coreModels.TenantConfig) error
	DeleteConfig(c echo.Context, config string) error
}
