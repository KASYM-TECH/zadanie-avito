//nolint:lll,gosimple
package controllers

import (
	"avito/domain"
	"avito/log"
	"context"
	"github.com/pkg/errors"
	"strconv"
)

type TenderService interface {
	Create(ctx context.Context, tender *domain.CreateTenderReq) (*domain.CreateTenderResp, error)
	GetPublished(ctx context.Context, offset, limit int, types []string) ([]domain.GetTendersResp, error)
	GetByUsername(ctx context.Context, offset, limit int, username string) ([]domain.GetTendersResp, error)
	GetStatus(ctx context.Context, tenderId, username string) (string, error)
	SetStatus(ctx context.Context, tenderId, status, username string) (*domain.SetStatusTenderResp, error)
	Edit(ctx context.Context, username, tenderId string, tender *domain.EditTenderReq) (*domain.EditTenderResp, error)
	Rollback(ctx context.Context, username string, tenderId string, version int) (*domain.RollbackTenderResp, error)
}

type TenderController struct {
	log           log.Logger
	tenderService TenderService
}

func NewTenderController(log log.Logger, tenderService TenderService) *TenderController {
	return &TenderController{log: log, tenderService: tenderService}
}

func (t *TenderController) Create(ctx context.Context, createReq domain.CreateTenderReq, rd domain.RequestData) (*domain.CreateTenderResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "tenderName", createReq.Name)
	t.log.Info(ctx, "tender create handler")

	tender, err := t.tenderService.Create(ctx, &createReq)
	if err == nil {
		return tender, nil
	}

	switch {
	case errors.Is(err, domain.ErrUserWithNameNotFound):
		return nil, &domain.HTTPError{Cause: err, Reason: "user with this username does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrOrganizationDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "organization with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) GetPublished(ctx context.Context, rd domain.RequestData) ([]domain.GetTendersResp, *domain.HTTPError) {
	var (
		offsetStr, limitStr string
		offset, limit       int
		types               []string
		err                 error
	)

	t.log.Info(ctx, "tender get handler")

	offsetStr, _ = ExtractQuery(rd.Request, "offset", "0")
	if offset, err = strconv.Atoi(offsetStr); err != nil && offset < 0 {
		offset = 0
	}

	limitStr, _ = ExtractQuery(rd.Request, "limit", "0")
	if limit, err = strconv.Atoi(limitStr); err != nil && limit < 0 {
		limit = 0
	}

	types, _ = ExtractQueryMany(rd.Request, "service_type")

	resp, err := t.tenderService.GetPublished(ctx, offset, limit, types)

	if err == nil {
		return resp, nil
	}

	return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
}

func (t *TenderController) GetByUsername(ctx context.Context, rd domain.RequestData) ([]domain.GetTendersResp, *domain.HTTPError) {
	var (
		username            string
		ok                  bool
		offsetStr, limitStr string
		offset, limit       int
		err                 error
	)

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "username", username)
	t.log.Info(ctx, "tender get handler")

	offsetStr, _ = ExtractQuery(rd.Request, "offset", "0")
	if offset, err = strconv.Atoi(offsetStr); err != nil && offset < 0 {
		offset = 0
	}

	limitStr, _ = ExtractQuery(rd.Request, "limit", "0")
	if limit, err = strconv.Atoi(limitStr); err != nil && limit < 0 {
		limit = 0
	}

	tenders, err := t.tenderService.GetByUsername(ctx, offset, limit, username)

	if err == nil {
		return tenders, nil
	}

	return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
}

func (t *TenderController) GetStatus(ctx context.Context, rd domain.RequestData) (string, *domain.HTTPError) {
	var (
		username, tenderId string
		ok                 bool
	)

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return "", &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if tenderId, ok = ExtractParam(rd.Request, "tenderId", ""); !ok {
		return "", &domain.HTTPError{Cause: nil, Reason: "tenderId is required", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "tenderId", tenderId)
	t.log.Info(ctx, "tender GetStatus handler")

	status, err := t.tenderService.GetStatus(ctx, tenderId, username)
	if err == nil {
		return status, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return "", &domain.HTTPError{Cause: err, Reason: "tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return "", &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return "", &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) SetStatus(ctx context.Context, rd domain.RequestData) (*domain.SetStatusTenderResp, *domain.HTTPError) {
	var (
		username, tenderId, status string
		ok                         bool
	)

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if tenderId, ok = ExtractParam(rd.Request, "tenderId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required", Status: domain.BadRequestCode}
	}

	if status, ok = ExtractQuery(rd.Request, "status", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "status is required query", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "username", username)
	t.log.Info(ctx, "tender SetStatus handler")

	tender, err := t.tenderService.SetStatus(ctx, tenderId, status, username)
	if err == nil {
		return tender, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) Edit(ctx context.Context, req domain.EditTenderReq, rd domain.RequestData) (*domain.EditTenderResp, *domain.HTTPError) {
	var (
		username, tenderId string
		ok                 bool
	)

	if tenderId, ok = ExtractParam(rd.Request, "tenderId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required", Status: domain.BadRequestCode}
	}
	ctx = log.AddKeyVal(ctx, "tenderId", tenderId)
	t.log.Info(ctx, "tender Edit handler")

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	tender, err := t.tenderService.Edit(ctx, username, tenderId, &req)
	if err == nil {
		return tender, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (t *TenderController) Rollback(ctx context.Context, rd domain.RequestData) (*domain.RollbackTenderResp, *domain.HTTPError) {
	var (
		username, tenderId, versionStr string
		version                        int
		ok                             bool
		err                            error
	)

	ctx = log.AddKeyVal(ctx, "tenderId", tenderId)
	t.log.Info(ctx, "tender Edit handler")

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if versionStr, ok = ExtractParam(rd.Request, "version", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version is required query", Status: domain.BadRequestCode}
	}

	if tenderId, ok = ExtractParam(rd.Request, "tenderId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version is required query", Status: domain.BadRequestCode}
	}

	if version, err = strconv.Atoi(versionStr); err != nil || version < 1 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version must be positive integer", Status: domain.BadRequestCode}
	}

	tender, err := t.tenderService.Rollback(ctx, username, tenderId, version)
	if err == nil {
		return tender, nil
	}

	switch {
	case errors.Is(err, domain.ErrTenderDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "tender with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "user does not belong to org", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}
