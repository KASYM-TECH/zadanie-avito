package service

import (
	"avito/db/model"
	"context"
	"github.com/pkg/errors"
)

type OrganizationRep interface {
	Insert(ctx context.Context, org *model.Organization) (string, error)
	MakeResponsible(ctx context.Context, empId, orgId string) (string, error)
}

type OrganizationService struct {
	orgRep OrganizationRep
}

func NewOrganizationService(orgRep OrganizationRep) OrganizationService {
	return OrganizationService{orgRep: orgRep}
}

func (u OrganizationService) Create(ctx context.Context, org *model.Organization) (string, error) {
	id, err := u.orgRep.Insert(ctx, org)
	if err != nil {
		return "", errors.WithMessage(err, "Service.Organization.Create insert org")
	}

	return id, nil
}

func (u OrganizationService) MakeResponsible(ctx context.Context, empId, orgId string) (string, error) {
	id, err := u.orgRep.MakeResponsible(ctx, empId, orgId)
	if err != nil {
		return "", errors.WithMessage(err, "Service.Organization.MakeResponsible bond")
	}

	return id, nil
}
