package trustedissuers

import (
	"github.com/gataca-io/vui-core/models"
	"github.com/labstack/echo/v4"
)

type TrustedIssuers interface {
	GetTrustedIssuer(ctx echo.Context, id string) (*models.TrustedIssuer, error)
	GetAllTrustedIssuers(ctx echo.Context) (*models.TrustedIssuerList, error)
}
