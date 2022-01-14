package service

import (
	"time"

	"github.com/gataca-io/vui-core/log"
	coreModels "github.com/gataca-io/vui-core/models"
	coreServices "github.com/gataca-io/vui-core/service"
	presentationexchange "github.com/gataca-io/vui-core/vui/presentationExchange"
	"github.com/labstack/echo/v4"
)

type daService struct {
	daRepo     presentationexchange.DataAgreementDao
	ssiService coreServices.SSIService
}

func NewDataAgreementService(daRepo presentationexchange.DataAgreementDao, ssiService coreServices.SSIService) presentationexchange.DataAgreementService {
	return &daService{
		daRepo:     daRepo,
		ssiService: ssiService,
	}
}

func (das *daService) Create(c echo.Context, da *coreModels.DataAgreement) (*coreModels.DataAgreement, error) {
	receiverID := da.DataReceiver.ID
	err := das.ssiService.VerifyDataAgreement(c, da, receiverID)
	if err != nil {
		log.CError(c, "Invalid data agreement", err)
		return nil, err
	}

	newEvent := coreModels.Event{
		PrincipleDid: receiverID,
		Version:      da.Version,
		TimeStamp:    time.Now().UnixMilli(),
		State:        "Capture",
	}
	da.Event = append(da.Event, newEvent)

	err = das.ssiService.SignDataAgreement(c, da, receiverID)
	if err != nil {
		log.CError(c, "Cannot sign updated data agreement", err)
		return nil, err
	}

	err = das.daRepo.Create(c, da)
	if err != nil {
		log.CError(c, "Cannot save data agreement in database", err)
		return nil, err
	}
	return da, nil
}

func (das *daService) GetDataAgreement(c echo.Context, id string, version int) (*coreModels.DataAgreement, error) {
	da, err := das.daRepo.GetDataAgreement(c, id)
	if err != nil {
		log.CErrorf(c, "Cannot retrieve data agreement %s from database", id, err)
		return nil, err
	}
	return da, nil
}

func (das *daService) Update(c echo.Context, da *coreModels.DataAgreement) (*coreModels.DataAgreement, error) {
	receiverID := da.DataReceiver.ID
	err := das.ssiService.VerifyDataAgreement(c, da, receiverID)
	if err != nil {
		log.CError(c, "Invalid data agreement", err)
		return nil, err
	}

	newEvent := coreModels.Event{
		PrincipleDid: receiverID,
		Version:      da.Version,
		TimeStamp:    time.Now().UnixMilli(),
		State:        "Modification",
	}
	da.Event = append(da.Event, newEvent)

	err = das.ssiService.SignDataAgreement(c, da, receiverID)
	if err != nil {
		log.CError(c, "Cannot sign updated data agreement", err)
		return nil, err
	}

	err = das.daRepo.UpdateDataAgreement(c, da)
	if err != nil {
		log.CError(c, "Cannot save data agreement in database", err)
		return nil, err
	}
	return da, nil
}

func (das *daService) Delete(c echo.Context, da *coreModels.DataAgreement) (*coreModels.DataAgreement, error) {
	receiverID := da.DataReceiver.ID
	err := das.ssiService.VerifyDataAgreement(c, da, receiverID)
	if err != nil {
		log.CError(c, "Invalid data agreement", err)
		return nil, err
	}

	newEvent := coreModels.Event{
		PrincipleDid: receiverID,
		Version:      da.Version,
		TimeStamp:    time.Now().UnixMilli(),
		State:        "Revocation",
	}
	da.Event = append(da.Event, newEvent)

	err = das.ssiService.SignDataAgreement(c, da, receiverID)
	if err != nil {
		log.CError(c, "Cannot sign updated data agreement", err)
		return nil, err
	}

	err = das.daRepo.UpdateDataAgreement(c, da)
	if err != nil {
		log.CError(c, "Cannot save data agreement in database", err)
		return nil, err
	}

	err = das.daRepo.DeleteLogical(c, da.ID)
	if err != nil {
		log.CError(c, "Cannot save data agreement in database", err)
		return nil, err
	}
	return da, nil
}
