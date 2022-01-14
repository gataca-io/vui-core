package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)

type StatusResponse struct {
	Status string `json:"status" example:"Ok" description:"Status of the user initalization."`
}

type versionHandler struct {
}

type VersionResponse struct {
	Version     string `json:"version" example:"1.0.0" description:"Version of the service running"`
	ReleaseDate string `json:"releaseDate" example:"2021-08-11T11:05:20Z" description:"Version of the service running"`
}

func NewVersionHandler(e *echo.Echo) {
	handler := &versionHandler{}

	e.GET("/api/v1/health", handler.health)
	e.GET("/api/v1/version", handler.version)
}

// Health godoc
// @Summary Check the health of the server.
// @Description For monitoring purposes. Returns 200 if the service is up and running.
// @Produce  json
// @Success 200 {object} StatusResponse "RUNNING"
// @Router /api/v1/health [get]
// @tag Common
// @tags Common, Monitoring,Core
func (vh *versionHandler) health(c echo.Context) error {
	return c.JSON(http.StatusOK, &StatusResponse{Status: "RUNNING"})
}

// Version godoc
// @Summary Check the current version of the service.
// @Description For monitoring purposes. Returns the current version of the service
// @Produce  json
// @Success 200 {object} VersionResponse "1.0.0"
// @Router /api/v1/version [get]
// @tag Common
// @tags Common, Monitoring,Core
func (vh *versionHandler) version(c echo.Context) error {
	return c.JSON(http.StatusOK, &VersionResponse{
		Version:     BuildVersion,
		ReleaseDate: BuildTime,
	})
}
