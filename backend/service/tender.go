//nolint:interfacebloat
package service

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/pkg/errors"
)

type TenderRep interface {
	Insert(ctx context.Context, newTender *model.Tender, authorUsername string) (string, error)
	GetPublished(ctx context.Context, offset, limit int, types []string) ([]model.Tender, error)
	GetById(ctx context.Context, tenderId string) (*model.Tender, error)
	GetByUsername(ctx context.Context, offset, limit int, username string) ([]model.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId string) (string, error)
	SetTenderStatus(ctx context.Context, tenderId, status string) error
	UpdateById(ctx context.Context, tender *model.Tender) error
	Rollback(ctx context.Context, tenderId string, version int) error
	AuthorByTenderId(ctx context.Context, tenderId string) (string, error)
	UsernameBelongsToTenderOrg(ctx context.Context, username, tenderId string) (bool, error)
}

type TenderService struct {
	tenderRep TenderRep
	orgRep    OrganizationRep
}

func NewTenderService(tenderRep TenderRep, orgRep OrganizationRep) TenderService {
	return TenderService{tenderRep: tenderRep, orgRep: orgRep}
}

func (t TenderService) Create(ctx context.Context, tender *domain.CreateTenderReq) (*domain.CreateTenderResp, error) {
	tenderDom := &model.Tender{
		Name:           tender.Name,
		Description:    tender.Description,
		Status:         tender.Status,
		ServiceType:    tender.ServiceType,
		OrganizationId: tender.OrganizationId,
	}

	isResponsible, err := t.orgRep.EmpBelongs(ctx, tender.CreatorUsername, tender.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	tenderId, err := t.tenderRep.Insert(ctx, tenderDom, tender.CreatorUsername)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender insert tender")
	}

	tenderNew, err := t.tenderRep.GetById(ctx, tenderId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender insert tender")
	}

	return &domain.CreateTenderResp{
		Id:          tenderNew.Id,
		Name:        tenderNew.Name,
		Description: tenderNew.Description,
		Status:      tenderNew.Status,
		ServiceType: tenderNew.ServiceType,
		CreatedAt:   tenderNew.CreatedAt,
		Version:     tenderNew.Version,
	}, nil
}

func (t TenderService) GetPublished(ctx context.Context, offset, limit int, types []string) ([]domain.GetTendersResp, error) {
	tenders, err := t.tenderRep.GetPublished(ctx, offset, limit, types)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.GetTendersResp, len(tenders))
	for i := range tenders {
		resp[i] = domain.GetTendersResp{
			Id:          tenders[i].Id,
			Name:        tenders[i].Name,
			Description: tenders[i].Description,
			Status:      tenders[i].Status,
			ServiceType: tenders[i].ServiceType,
			CreatedAt:   tenders[i].CreatedAt,
			Version:     tenders[i].Version,
		}
	}

	return resp, nil
}

func (t TenderService) GetByUsername(ctx context.Context, offset, limit int, username string) ([]domain.GetTendersResp, error) {
	tenders, err := t.tenderRep.GetByUsername(ctx, offset, limit, username)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.GetTendersResp, len(tenders))
	for i := range tenders {
		resp[i] = domain.GetTendersResp{
			Id:          tenders[i].Id,
			Name:        tenders[i].Name,
			Description: tenders[i].Description,
			Status:      tenders[i].Status,
			ServiceType: tenders[i].ServiceType,
			CreatedAt:   tenders[i].CreatedAt,
			Version:     tenders[i].Version,
		}
	}

	return resp, nil
}

func (t TenderService) GetStatus(ctx context.Context, tenderId, username string) (string, error) {
	authorName, err := t.tenderRep.AuthorByTenderId(ctx, tenderId)
	if err != nil {
		return "", err
	}

	status, err := t.tenderRep.GetTenderStatus(ctx, tenderId)

	if username != authorName && status == string(model.BidStatusCreated) {
		return "", domain.ErrTenderDoesNotExist
	}

	if err != nil {
		return "", err
	}

	return status, nil
}

func (t TenderService) SetStatus(ctx context.Context, tenderId, status, username string) (*domain.SetStatusTenderResp, error) {
	isResponsible, err := t.tenderRep.UsernameBelongsToTenderOrg(ctx, username, tenderId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	err = t.tenderRep.SetTenderStatus(ctx, tenderId, status)

	if err != nil {
		return nil, err
	}

	tenderUpdated, err := t.tenderRep.GetById(ctx, tenderId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender set status")
	}

	tenderDom := &domain.SetStatusTenderResp{
		Id:          tenderUpdated.Id,
		Name:        tenderUpdated.Name,
		Description: tenderUpdated.Description,
		Status:      tenderUpdated.Status,
		CreatedAt:   tenderUpdated.CreatedAt,
		Version:     tenderUpdated.Version,
		ServiceType: tenderUpdated.ServiceType,
	}

	return tenderDom, nil
}

func (t TenderService) Edit(ctx context.Context, username, tenderId string, tender *domain.EditTenderReq) (*domain.EditTenderResp, error) {
	tenderEdit := model.Tender{
		Id:          tenderId,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
	}

	isResponsible, err := t.tenderRep.UsernameBelongsToTenderOrg(ctx, username, tenderId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	err = t.tenderRep.UpdateById(ctx, &tenderEdit)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender edit tender")
	}

	tenderUpdated, err := t.tenderRep.GetById(ctx, tenderId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender get tender")
	}

	tenderDom := &domain.EditTenderResp{
		Id:          tenderUpdated.Id,
		Name:        tenderUpdated.Name,
		Description: tenderUpdated.Description,
		Status:      tenderUpdated.Status,
		CreatedAt:   tenderUpdated.CreatedAt,
		Version:     tenderUpdated.Version,
		ServiceType: tenderUpdated.ServiceType,
	}
	return tenderDom, nil
}

func (t TenderService) Rollback(ctx context.Context, username string, tenderId string, version int) (*domain.RollbackTenderResp, error) {
	isResponsible, err := t.tenderRep.UsernameBelongsToTenderOrg(ctx, username, tenderId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	err = t.tenderRep.Rollback(ctx, tenderId, version)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender rollback tender")
	}

	tenderUpdated, err := t.tenderRep.GetById(ctx, tenderId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender rollback")
	}

	tenderDom := &domain.RollbackTenderResp{
		Id:          tenderUpdated.Id,
		Name:        tenderUpdated.Name,
		Description: tenderUpdated.Description,
		Status:      tenderUpdated.Status,
		CreatedAt:   tenderUpdated.CreatedAt,
		Version:     tenderUpdated.Version,
		ServiceType: tenderUpdated.ServiceType,
	}

	return tenderDom, nil
}
