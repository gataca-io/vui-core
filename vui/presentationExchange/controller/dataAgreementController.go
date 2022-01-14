package controller

import (
	"net/http"
	"strconv"

	"github.com/gataca-io/vui-core/log"
	coreModels "github.com/gataca-io/vui-core/models"
	presentationexchange "github.com/gataca-io/vui-core/vui/presentationExchange"
	"github.com/labstack/echo/v4"
)

type daHandler struct {
	daService presentationexchange.DataAgreementService
}

// NewDataAgreementHandler godoc
// Create a Controller for the Data Agreement API
func NewDataAgreementHandler(e *echo.Echo, daService presentationexchange.DataAgreementService) {

	handler := &daHandler{
		daService: daService,
	}
	e.GET("/api/v2/data_agreements/:id/:version", handler.getDataAgreement)
	e.GET("/api/v2/data_agreements/:id", handler.getDataAgreement)
	e.POST("/api/v2/data_agreements", handler.createDataAgreement)
	e.DELETE("/api/v2/data_agreements/:id", handler.deleteDataAgreement)
	e.PATCH("/api/v2/data_agreements/:id", handler.updateDataAgreement)
}

// CreateDataAgreement godoc
// @Summary Create Data Agreement
// @Description Create a new data agreement to record by the Verifier, in order to hold the service
// @Accept  json
// @Produce  json
// @Param dataAgreement body coreModels.DataAgreement true "Data Agreement of this service"
// @Success 201 {object} coreModels.DataAgreement "Updated Data Agreement"
// @Failure 400 {object} models.ResponseMessage "Invalid input data."
// @Failure 403 {object} models.ResponseMessage "Not Authorized to create data agreements."
// @Failure 500 {object} models.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/data_agreements [post]
// @tag DataAgreements
// @tags DataAgreements,Connect
// @security Token
func (h *daHandler) createDataAgreement(c echo.Context) error {
	da := &coreModels.DataAgreement{}

	err := c.Bind(da)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	da, err = h.daService.Create(c, da)
	if err != nil {
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, da)
}

// GetDataAgreement godoc
// @Summary Retrieve an existing Data Agreement
// @Description Retrieve the last version of a Data Agreement with all the corresponding events.
// @Accept  json
// @Produce  json
// @Param id path string false "Data agreement Id"
// @Param version path string false "Optional: Version of the data agreement to recover. Default: last"
// @Success 200 {array} coreModels.DataAgreement "Data Agreement with that Id"
// @Failure 404 {object} models.ResponseMessage "Inexistent process Id"
// @Failure 409 {object} models.ResponseMessage "Process Id cannot be retrieved"
// @Failure 500 {object} models.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/data_agreements/{id}/{version} [get]
// @tag DataAgreements
// @tags DataAgreements,Connect
func (h *daHandler) getDataAgreement(c echo.Context) error {
	id := c.Param("id")
	versionParam := c.Param("version")
	version := -1
	if versionParam != "" {
		var err error
		version, err = strconv.Atoi(versionParam)
		if err != nil {
			log.CWarn(c, "Wrong version param, providing last version")
		}
	}

	da, err := h.daService.GetDataAgreement(c, id, version)
	if err != nil {
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, da)
}

// updateDataAgreement godoc
// @Summary Modify an existing data agreement
// @Description Use cases not implemented yet
// @Accept  json
// @Produce  json
// @Param id path string false "Data agreement Id"
// @Param dataAgreement body coreModels.DataAgreement true "Data Agreement of this service"
// @Success 200 {array} coreModels.DataAgreement "Updated data agreement"
// @Failure 404 {object} models.ResponseMessage "Inexistent process Id"
// @Failure 409 {object} models.ResponseMessage "Process Id cannot be retrieved"
// @Failure 500 {object} models.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/data_agreements/{id} [patch]
// @tag DataAgreements
// @tags DataAgreements,Connect
func (h *daHandler) updateDataAgreement(c echo.Context) error {
	id := c.Param("id")
	da := &coreModels.DataAgreement{}

	err := c.Bind(da)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if id != da.ID {
		log.CError(c, "Id param and body not matching")
		return c.JSON(http.StatusBadRequest, coreModels.ResponseMessage{Message: "Id param and body not matching"})
	}

	da, err = h.daService.Update(c, da)
	if err != nil {
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, da)
}

// deleteDataAgreement godoc
// @Summary Delete a current data agreement
// @Description Use case not implemented yet
// @Accept  json
// @Produce  json
// @Param id path string false "Data agreement Id"
// @Success 200 {array} coreModels.DataAgreement "Revocated Data Agreement"
// @Failure 404 {object} models.ResponseMessage "Inexistent process Id"
// @Failure 409 {object} models.ResponseMessage "Process Id cannot be retrieved"
// @Failure 500 {object} models.ResponseMessage "Serverside error processing the request."
// @Router /api/v2/data_agreements/{id} [delete]
// @tag DataAgreements
// @tags DataAgreements,Connect
func (h *daHandler) deleteDataAgreement(c echo.Context) error {
	id := c.Param("id")
	da := &coreModels.DataAgreement{}

	err := c.Bind(da)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if id != da.ID {
		log.CError(c, "Id param and body not matching")
		return c.JSON(http.StatusBadRequest, coreModels.ResponseMessage{Message: "Id param and body not matching"})
	}

	da, err = h.daService.Delete(c, da)
	if err != nil {
		return c.JSON(getStatusCode(err), coreModels.ResponseMessage{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, da)
}
