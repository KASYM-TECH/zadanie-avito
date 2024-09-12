//nolint:interfacebloat
package service

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/pkg/errors"
)

type TenderRep interface {
	Insert(ctx context.Context, newTender *model.Tender) (string, error)
	GetPublished(ctx context.Context, offset, limit int, types []string) ([]model.Tender, error)
	GetById(ctx context.Context, tenderId string) (*model.Tender, error)
	GetByUsername(ctx context.Context, offset, limit int, username string) ([]model.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId string) (string, error)
	SetTenderStatus(ctx context.Context, tenderId, status string) error
	UpdateById(ctx context.Context, tender *model.Tender) error
	Rollback(ctx context.Context, tenderId string, version int) error
	AuthorByTenderId(ctx context.Context, tenderId string) (string, error)
	UsernameBelongsToTenderOrg(ctx context.Context, username, tenderId string) (bool, error)
	UserIdBelongsToTenderOrg(ctx context.Context, userId, tenderId string) (bool, error)
}

type OrgRep interface {
	EmpBelongs(ctx context.Context, empId, orgId string) (bool, error)
}

type TenderService struct {
	tenderRep TenderRep
	orgRep    OrgRep
}

func NewTenderService(tenderRep TenderRep, orgRep OrgRep) TenderService {
	return TenderService{tenderRep: tenderRep, orgRep: orgRep}
}

func (t TenderService) Create(ctx context.Context, tender *model.Tender) (*model.Tender, error) {
	isResponsible, err := t.orgRep.EmpBelongs(ctx, tender.UserId, tender.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	tenderId, err := t.tenderRep.Insert(ctx, tender)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender insert tender")
	}

	tenderNew, err := t.tenderRep.GetById(ctx, tenderId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender insert tender")
	}

	return tenderNew, nil
}

func (t TenderService) GetPublished(ctx context.Context, offset, limit int, types []string) ([]model.Tender, error) {
	tenders, err := t.tenderRep.GetPublished(ctx, offset, limit, types)
	if err != nil {
		return nil, err
	}

	return tenders, nil
}

func (t TenderService) GetByUsername(ctx context.Context, offset, limit int, username string) ([]model.Tender, error) {
	tenders, err := t.tenderRep.GetByUsername(ctx, offset, limit, username)
	if err != nil {
		return nil, err
	}

	return tenders, nil
}

func (t TenderService) GetStatus(ctx context.Context, tenderId string) (string, error) {
	status, err := t.tenderRep.GetTenderStatus(ctx, tenderId)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (t TenderService) SetStatus(ctx context.Context, tenderId, status, username string) (*model.Tender, error) {
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

	return tenderUpdated, nil
}

func (t TenderService) Edit(ctx context.Context, username string, tender *model.Tender) (*model.Tender, error) {
	isResponsible, err := t.tenderRep.UsernameBelongsToTenderOrg(ctx, username, tender.ID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	err = t.tenderRep.UpdateById(ctx, tender)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender edit tender")
	}

	tenderUpdated, err := t.tenderRep.GetById(ctx, tender.ID)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Tender get tender")
	}

	return tenderUpdated, nil
}

func (t TenderService) Rollback(ctx context.Context, username string, tenderId string, version int) (*model.Tender, error) {
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

	return tenderUpdated, nil
}
