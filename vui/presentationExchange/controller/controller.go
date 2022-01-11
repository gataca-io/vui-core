package controller

import (
	"net/http"

	coreModels "github.com/gataca-io/vui-core/models"
	presentationexchange "github.com/gataca-io/vui-core/vui/presentationExchange"
	"github.com/labstack/echo/v4"
)

type peHandler struct {
	exchangeService presentationexchange.PresExchangeService
	baseURI         string
}

type PECreationResponse struct {
	ID  string `json:"id" example:"32f54163-7166-48f1-93d8-ff217bdb0653" description:"Presentation Exchange unique id"`
	URI string `json:"uri" example:"https://vui.gataca.io/api/presentations/v2/32f54163-7166-48f1-93d8-ff217bdb0653/definition" description:"URI to retrieve the presentation Definition for the exchange"`
}

type SIOPSubmission struct {
	VPToken struct {
		Format       string                            `json:"format" example:"ldp_vp" description:"Format of the verifiable presentation"`
		Presentation coreModels.VerifiablePresentation `json:"presentation" example:"" description:"Presentation in the stablished format"`
	} `json:"vp_token" description:"Verifiable Presentation as token"`
}

// NewPresentationExchangeHandler godoc
// Create a Controller for the Presentation Exchange API
func NewPresentationExchangeHandler(e *echo.Echo, exchangeService presentationexchange.PresExchangeService, baseURI string) {

	handler := &peHandler{
		exchangeService: exchangeService,
		baseURI:         baseURI,
	}
	e.POST("/api/v2/presentations", handler.createPresentationExchange)
	e.GET("/api/v2/presentations/:id/definition", handler.getPresentationDefinition)
	e.GET("/api/v2/presentations/:id/data_agreement", handler.getPresentationDataAgreement)
	e.POST("/api/v2/presentations/:id/submission", handler.submitPresentation)
	e.POST("/api/v2/authentication_responses", handler.submitSIOPToken) //does the same that the previous endpoint but updated
	e.GET("/api/v2/presentations/:id", handler.checkStatus)

}

// CreatePresentationExchange godoc
// @Summary Create Presentation Exchange
// @Description Create a new presentation exchange process by providing it's presentation definition. Relying parties with due authentication can perform this operation.
// @Accept  json
// @Produce  json
// @Param presentationDefinition body orecoreModels.PresentationDefinition true "Presentation definition of this exchange"
// @Success 201 {object} PECreationResponse "Reference to the exchange process"
// @Failure 400 {object} coreModels.ResponseMessage "Invalid input data."
// @Failure 403 {object} coreModels.ResponseMessage "Not Authorized to create exchanges."
// @Failure 500 {object} coreModels.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/presentations [post]
// @tag Presentations
// @tags Presentations,Connect
// @security Token
func (h *peHandler) createPresentationExchange(c echo.Context) error {
	var pr coreModels.PresentationDefinition

	err := c.Bind(&pr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	pe, err := h.exchangeService.Create(c, &pr)
	if err != nil {
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}

	response := PECreationResponse{
		ID:  pe.Id,
		URI: h.baseURI + "/api/v2/presentations/" + pe.Id + "/definition",
	}

	return c.JSON(http.StatusCreated, response)
}

// GetPresentationDefinition godoc
// @Summary Get a Presentation Definition
// @Description Upon scanning a QR, a Holder may retrieve the presentation definition associated to the process identifier in order to perform an exchange.
// @Accept  json
// @Produce  json
// @Param id path string false "Presentation exchange Id"
// @Success 200 {array} coreModels.PresentationDefinition "Tenant configurations requested."
// @Failure 404 {object} coreModels.ResponseMessage "Inexistent process Id"
// @Failure 409 {object} coreModels.ResponseMessage "Process Id cannot be retrieved"
// @Failure 500 {object} coreModels.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/presentations/{id}/definition [get]
// @tag Presentations
// @tags Presentations,Connect
func (h *peHandler) getPresentationDefinition(c echo.Context) error {
	id := c.Param("id")

	def, err := h.exchangeService.GetDefinition(c, id, false)
	if err != nil {
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}

	def.DataAgreement = &coreModels.DataAgreementRef{
		Ref: h.baseURI + "/api/v2/presentations/" + def.ID + "/data_agreement",
	}

	return c.JSON(http.StatusOK, def)
}

// GetPresentationDataAgreement godoc
// @Summary Get a the data agreement template of a Presentation
// @Description When expanding a presentation, the verifier may just offer the URI to the data agreement template linked to that service
// @Accept  json
// @Produce  json
// @Param id path string false "Presentation exchange Id"
// @Success 200 {array} coreModels.DataAgreement "Tenant configurations requested."
// @Failure 404 {object} coreModels.ResponseMessage "Inexistent process Id"
// @Failure 409 {object} coreModels.ResponseMessage "Process Id cannot be retrieved"
// @Failure 500 {object} coreModels.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/presentations/{id}/data_agreement [get]
// @tag Presentations
// @tags Presentations,Connect
func (h *peHandler) getPresentationDataAgreement(c echo.Context) error {
	id := c.Param("id")

	def, err := h.exchangeService.GetDefinition(c, id, true)
	if err != nil {
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, def.DataAgreement)
}

// SubmitPresentation godoc
// @Summary Submit a Verifiable Presentation
// @Description A Holder may submit a verifiable presentation in response to a given definition in order to fulfill the exchange.
// @Accept  json
// @Produce  json
// @Param id path string false "Presentation exchange Id"
// @Param submission body coreModels.VerifiablePresentation true "Verifiable Presentation"
// @Success 200 {array} coreModels.VerificationResult "Verification result."
// @Failure 400 {object} coreModels.ResponseMessage "Request body malformed"
// @Failure 403 {object} coreModels.ResponseMessage "Not Authorized to submit a presentation exchanges."
// @Failure 404 {object} coreModels.ResponseMessage "Inexistent process Id"
// @Failure 406 {object} coreModels.ResponseMessage "Presentation submission not acceptable"
// @Failure 409 {object} coreModels.ResponseMessage "Process Id cannot be modified"
// @Failure 500 {object} coreModels.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/presentations/{id}/submission [post]
// @tag Presentations
// @tags Presentations,Connect
func (h *peHandler) submitPresentation(c echo.Context) error {
	id := c.Param("id")
	var vp coreModels.VerifiablePresentation

	err := c.Bind(&vp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	verification, err := h.exchangeService.Submit(c, id, &vp)
	if err != nil {
		if verification != nil {
			return c.JSON(http.StatusNotAcceptable, verification)
		}
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, verification)

}

// SubmitSIOPToken godoc
// @Summary Submit a Verifiable Presentation under the siop standard
// @Description A Holder may submit a verifiable presentation in response to a given authentication_request in order to finish the exchange.
// @Accept  json
// @Produce  json
// @Param submission body SIOPSubmission true "Verifiable Presentation token for DID SIOP"
// @Success 200 {array} coreModels.VerificationResult "Verification result."
// @Failure 400 {object} coreModels.ResponseMessage "Request body malformed"
// @Failure 403 {object} coreModels.ResponseMessage "Not Authorized to submit a presentation exchanges."
// @Failure 404 {object} coreModels.ResponseMessage "Inexistent process Id"
// @Failure 406 {object} coreModels.ResponseMessage "Presentation submission not acceptable"
// @Failure 409 {object} coreModels.ResponseMessage "Process Id cannot be modified"
// @Failure 500 {object} coreModels.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/authentication_responses [post]
// @tag Presentations
// @tags Presentations,Connect
func (h *peHandler) submitSIOPToken(c echo.Context) error {
	var token SIOPSubmission

	err := c.Bind(&token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	verification, err := h.exchangeService.Submit(c, token.VPToken.Presentation.PresentationSubmission.DefinitionID, &token.VPToken.Presentation)
	if err != nil {
		if verification != nil {
			return c.JSON(http.StatusNotAcceptable, verification)
		}
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, verification)

}

// CheckStatus godoc
// @Summary Check the status of a presentation exchange
// @Description The relying party may at any time query the status of a given exchange at any time to see if the data has been validated.
// @Accept  json
// @Produce  json
// @Param id path string false "Presentation exchange Id"
// @Success 200 {array} coreModels.VerificationResult "Valid verification result."
// @Success 202 {array} coreModels.VerificationResult "Pending verification result. No submission in the exchange yet."
// @Failure 400 {object} coreModels.ResponseMessage "Process Id cannot be retrieved"
// @Failure 403 {object} coreModels.ResponseMessage "Not Authorized to retrieve the presentation exchange"
// @Failure 404 {object} coreModels.ResponseMessage "Inexistent process Id"
// @Failure 406 {object} coreModels.VerificationResult "Presentation submission in valid"
// @Failure 500 {object} coreModels.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/presentations/{id} [get]
// @tag Presentations
// @tags Presentations,Connect
func (h *peHandler) checkStatus(c echo.Context) error {
	id := c.Param("id")

	verification, err := h.exchangeService.GetVerification(c, id)
	if err != nil {

		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}
	if verification != nil {
		if !verification.Valid() {
			return c.JSON(http.StatusNotAcceptable, verification)
		}
		return c.JSON(http.StatusOK, verification)
	}
	return c.JSON(http.StatusAccepted, verification)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case coreModels.ErrInternalServerError:
		return http.StatusInternalServerError
	case coreModels.ErrNotFound:
		return http.StatusNotFound
	case coreModels.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
