//nolint:lll,gosimple
package controllers

import (
	"avito/db/model"
	"avito/domain"
	"avito/log"
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"strconv"
)

type TenderService interface {
	Create(ctx context.Context, tender *model.Tender) (*model.Tender, error)
	GetPublished(ctx context.Context, offset, limit int, types []string) ([]model.Tender, error)
	GetByUsername(ctx context.Context, offset, limit int, username string) ([]model.Tender, error)
	GetStatus(ctx context.Context, tenderId string) (string, error)
	SetStatus(ctx context.Context, tenderId, status, username string) (*model.Tender, error)
	Edit(ctx context.Context, username string, tender *model.Tender) (*model.Tender, error)
	Rollback(ctx context.Context, username string, tenderId string, version int) (*model.Tender, error)
}

type TenderController struct {
	log           log.Logger
	tenderService TenderService
}

func NewTenderController(log log.Logger, tenderService TenderService) *TenderController {
	return &TenderController{log: log, tenderService: tenderService}
}

func (t *TenderController) Create(ctx context.Context, createReq domain.CreateTenderReq, rd domain.RequestData) (*domain.CreateTenderResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "name", createReq.Name)
	t.log.Info(ctx, "tender create handler")

	if rd.Claims.Username != createReq.CreatorUsername {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	tenderDom := &model.Tender{
		Name:           createReq.Name,
		Description:    createReq.Description,
		Status:         createReq.Status,
		ServiceType:    createReq.ServiceType,
		OrganizationID: createReq.OrganizationID,
		UserId:         rd.UserId,
	}

	tender, err := t.tenderService.Create(ctx, tenderDom)

	if err == nil {
		return &domain.CreateTenderResp{
			Id:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			ServiceType: tender.ServiceType,
			CreatedAt:   tender.CreatedAt,
			Version:     tender.Version,
		}, nil
	}

	pgErr := &pgconn.PgError{}
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return nil, &domain.HTTPError{Cause: err, Reason: "organization with this id does not exist", Status: domain.BadRequestCode}
		}
	}

	switch {
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) GetPublished(ctx context.Context, rd domain.RequestData) ([]domain.GetTendersResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	t.log.Info(ctx, "tender Get handler")

	var (
		offset, limit int
		err           error
		types         []string
	)

	if offsetStr, ok := rd.Request.URL.Query()["offset"]; ok {
		if offset, err = strconv.Atoi(offsetStr[0]); err != nil && offset < 0 {
			offset = 0
		}
	}

	if limitStr, ok := rd.Request.URL.Query()["limit"]; ok {
		if limit, err = strconv.Atoi(limitStr[0]); err != nil && limit < 0 {
			limit = 0
		}
	}

	types, _ = rd.Request.URL.Query()["service_type"]

	tenders, err := t.tenderService.GetPublished(ctx, offset, limit, types)
	resp := make([]domain.GetTendersResp, len(tenders))
	for i := range tenders {
		resp[i] = domain.GetTendersResp{
			Id:          tenders[i].ID,
			Name:        tenders[i].Name,
			Description: tenders[i].Description,
			Status:      tenders[i].Status,
			ServiceType: tenders[i].ServiceType,
			CreatedAt:   tenders[i].CreatedAt,
			Version:     tenders[i].Version,
		}
	}

	if err == nil {
		return resp, nil
	}

	return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
}

func (t *TenderController) GetByUsername(ctx context.Context, rd domain.RequestData) ([]domain.GetTendersResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	t.log.Info(ctx, "tender GetByUserId handler")

	username, ok := rd.Request.URL.Query()["username"]
	if !ok || len(username) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	var (
		offset, limit int
		err           error
	)

	if offsetStr, ok := rd.Request.URL.Query()["offset"]; ok {
		if offset, err = strconv.Atoi(offsetStr[0]); err != nil && offset < 0 {
			offset = 0
		}
	}

	if limitStr, ok := rd.Request.URL.Query()["limit"]; ok {
		if limit, err = strconv.Atoi(limitStr[0]); err != nil && limit < 0 {
			limit = 0
		}
	}

	tenders, err := t.tenderService.GetByUsername(ctx, offset, limit, username[0])
	resp := make([]domain.GetTendersResp, len(tenders))
	for i := range tenders {
		resp[i] = domain.GetTendersResp{
			Id:          tenders[i].ID,
			Name:        tenders[i].Name,
			Description: tenders[i].Description,
			Status:      tenders[i].Status,
			ServiceType: tenders[i].ServiceType,
			CreatedAt:   tenders[i].CreatedAt,
			Version:     tenders[i].Version,
		}
	}

	if err == nil {
		return resp, nil
	}

	return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
}

func (t *TenderController) GetStatus(ctx context.Context, rd domain.RequestData) (string, *domain.HTTPError) {
	tenderId, ok := mux.Vars(rd.Request)["tenderId"]
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	t.log.Info(ctx, "tender GetStatus handler")

	if !ok || len(tenderId) == 0 {
		return "", &domain.HTTPError{Cause: nil, Reason: "tenderId is required query", Status: domain.BadRequestCode}
	}

	status, err := t.tenderService.GetStatus(ctx, tenderId)
	if err == nil {
		return status, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return "", &domain.HTTPError{Cause: err, Reason: "Tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return "", &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return "", &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) SetStatus(ctx context.Context, rd domain.RequestData) (*domain.SetStatusTenderResp, *domain.HTTPError) {
	tenderId, ok := mux.Vars(rd.Request)["tenderId"]
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	t.log.Info(ctx, "tender SetStatus handler")

	if !ok || len(tenderId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required query", Status: domain.BadRequestCode}
	}

	var (
		username, status []string
	)

	if status, ok = rd.Request.URL.Query()["status"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "status is required query", Status: domain.BadRequestCode}
	}

	if username, ok = rd.Request.URL.Query()["username"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	tender, err := t.tenderService.SetStatus(ctx, tenderId, status[0], username[0])
	if err == nil {
		tenderDom := &domain.SetStatusTenderResp{
			Id:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			CreatedAt:   tender.CreatedAt,
			Version:     tender.Version,
			ServiceType: tender.ServiceType,
		}
		return tenderDom, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) Edit(ctx context.Context, req domain.EditTenderReq, rd domain.RequestData) (*domain.EditTenderResp, *domain.HTTPError) {
	tenderId, ok := mux.Vars(rd.Request)["tenderId"]
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	t.log.Info(ctx, "tender Edit handler")

	if !ok || len(tenderId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required query", Status: domain.BadRequestCode}
	}

	var username []string
	if username, ok = rd.Request.URL.Query()["username"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	tenderReq := model.Tender{
		ID:          tenderId,
		Name:        req.Name,
		Description: req.Description,
		ServiceType: req.ServiceType,
		UserId:      rd.UserId,
	}

	tender, err := t.tenderService.Edit(ctx, username[0], &tenderReq)
	if err == nil {
		tenderDom := &domain.EditTenderResp{
			Id:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			CreatedAt:   tender.CreatedAt,
			Version:     tender.Version,
			ServiceType: tender.ServiceType,
		}
		return tenderDom, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) Rollback(ctx context.Context, rd domain.RequestData) (*domain.RollbackTenderResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	t.log.Info(ctx, "tender Rollback handler")

	tenderId, ok := mux.Vars(rd.Request)["tenderId"]
	if !ok || len(tenderId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required query", Status: domain.BadRequestCode}
	}
	versionStr, ok := mux.Vars(rd.Request)["version"]
	if !ok || len(versionStr) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version is required query", Status: domain.BadRequestCode}
	}

	var username []string
	if username, ok = rd.Request.URL.Query()["username"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	var (
		version int
		err     error
	)
	if version, err = strconv.Atoi(versionStr); err != nil || version < 1 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version must be positive integer", Status: domain.BadRequestCode}
	}

	tender, err := t.tenderService.Rollback(ctx, username[0], tenderId, version)
	if err == nil {
		tenderDom := &domain.RollbackTenderResp{
			Id:          tender.ID,
			Name:        tender.Name,
			Description: tender.Description,
			Status:      tender.Status,
			CreatedAt:   tender.CreatedAt,
			Version:     tender.Version,
			ServiceType: tender.ServiceType,
		}
		return tenderDom, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}
