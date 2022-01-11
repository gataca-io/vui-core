package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gataca-io/vui-core/log"
	coreModels "github.com/gataca-io/vui-core/models"
	"github.com/labstack/echo/v4"
)

type universalDidService struct {
	contextTimeout        time.Duration
	universalResolverHost string
}

const RouteResolver = "/identifiers/"

// NewDidService
func NewUniversalDidService(timeout time.Duration, universalUrl string) DidService {
	return &universalDidService{
		contextTimeout:        timeout,
		universalResolverHost: universalUrl,
	}
}

// Create
func (s *universalDidService) CreateDID(ctx echo.Context, did *coreModels.DIDDocument) error {
	return coreModels.ErrInvalidDIDMethod
}

// GetByAddress
func (s *universalDidService) GetDID(ctx echo.Context, did string) (*coreModels.DIDDocument, error) {

	url := s.universalResolverHost + RouteResolver + did
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.CError(ctx, "Error creating http request: ", err.Error())
		return nil, err
	}

	req.Header.Set("Accept", "application/did+ld+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.CError(ctx, "Error connecting to Universal resolver service ", url, err.Error())
		return nil, err
	}
	if resp.StatusCode >= 400 {
		log.CError(ctx, "Error on universal resolver request ", url, resp.StatusCode)
		return nil, coreModels.ErrNotFound
	}

	log.CDebug(ctx, "DID resolved ok", url)
	var didDoc coreModels.DIDDocument
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&didDoc) //Allow unkown fields by default
	if err != nil {
		log.CError(ctx, "Error parsing DID response", resp.Body, didDoc)
		return nil, err
	}
	log.CDebugf(ctx, "Response to DID Request", url, didDoc)
	return &didDoc, nil
}

func (s *universalDidService) UpdateDID(ctx echo.Context, didDocument *coreModels.DIDDocument) error {
	return coreModels.ErrInvalidDIDMethod
}

func (s *universalDidService) GetByDIDWithFilters(ctx echo.Context, did string, timestamp *time.Time) ([]*coreModels.DIDDocument, error) {
	return nil, coreModels.ErrInvalidDIDMethod
}

func (s *universalDidService) RevokeDID(ctx echo.Context, didDocument *coreModels.DIDDocument) error {
	return coreModels.ErrInvalidDIDMethod
}
