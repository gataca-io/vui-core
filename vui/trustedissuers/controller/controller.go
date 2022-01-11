package controller

import (
	"net/http"

	"github.com/gataca-io/vui-core/models"
	"github.com/gataca-io/vui-core/vui/trustedissuers"
	"github.com/labstack/echo/v4"
)

type trustIssuerHandler struct {
	s trustedissuers.TrustedIssuers
}

func NewTrustedIssuerHandler(e *echo.Echo, s trustedissuers.TrustedIssuers) {
	handler := &trustIssuerHandler{
		s: s,
	}
	//Trusted Issuers
	e.GET("/api/v1/trusted-issuers", handler.getAllTrustedIssuers)
	e.GET("/api/v1/trusted-issuers/:id", handler.getTrustedIssuer)
}

// ********************
// *    Operations    *
// ********************

// GetAllTrustedIssuers godoc
// @Summary Get the information related to the Trusted Issuers
// @Description Retrieve the information related to the Trusted Issuers. That information must contains all the legal information, API services, eIDAS certificates, ... that provide a trusted relation with the legal entity.
// @Accept  json
// @Produce  json
// @Success 200 {object} StatusResponse Trusted Issuer List
// @Failure 500 {object} StatusResponse "Serverside error processing the request."
// @Router /api/v1/trusted-issuers [get]
// @tag trusted-issuers
// @tags trusted-issuers,vui
func (h *trustIssuerHandler) getAllTrustedIssuers(ctx echo.Context) error {

	til, err := h.s.GetAllTrustedIssuers(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrInternalServerError)
	}

	return ctx.JSON(http.StatusOK, til)
}

// GetAllTrustedIssuers godoc
// @Summary Get the information related to a concrete Trusted Issuer
// @Description Retrieve the information related to the Trusted Issuers. That information must contains all the legal information, API services, eIDAS certificates, ... that provide a trusted relation with the legal entity.
// @Accept  json
// @Produce  json
// @Success 200 {object} StatusResponse Trusted Issuer List
// @Failure 404 {object} StatusResponse "Not found trusted issuer for this {id}"
// @Failure 500 {object} StatusResponse "Serverside error processing the request."
// @Router /api/v1/trusted-issuers/{id} [get]
// @tag trusted-issuer
// @tags trusted-issuer,vui
func (h *trustIssuerHandler) getTrustedIssuer(ctx echo.Context) error {
	id := ctx.Param("id")

	ti, err := h.s.GetTrustedIssuer(ctx, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.ErrInternalServerError)
	}

	return ctx.JSON(http.StatusOK, ti)
}
