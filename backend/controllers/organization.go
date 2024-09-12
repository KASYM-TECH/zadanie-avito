//nolint:lll
package controllers

import (
	"avito/db/model"
	"avito/domain"
	"avito/log"
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type OrganizationService interface {
	Create(ctx context.Context, org *model.Organization) (string, error)
	MakeResponsible(ctx context.Context, empId, orgId string) (string, error)
}

type OrganizationController struct {
	log        log.Logger
	orgService OrganizationService
}

func NewOrganizationController(log log.Logger, orgService OrganizationService) *OrganizationController {
	return &OrganizationController{log: log, orgService: orgService}
}

func (t *OrganizationController) Create(ctx context.Context, createReq domain.CreateOrganizationReq, rd domain.RequestData) (string, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "name", createReq.Name)
	t.log.Info(ctx, "org create handler")

	org := &model.Organization{
		Name:        createReq.Name,
		Description: createReq.Description,
		Type:        createReq.Type,
	}
	id, err := t.orgService.Create(ctx, org)

	if err == nil {
		return id, nil
	}

	return "", &domain.HTTPError{Cause: err, Reason: "could not create organization", Status: domain.ServerFailureCode}
}

func (o *OrganizationController) MakeResponsible(ctx context.Context, bondReq domain.BondReq, rd domain.RequestData) (string, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "orgId", bondReq.OrganizationId)
	o.log.Info(ctx, "MakeResponsible handler")

	bondId, err := o.orgService.MakeResponsible(ctx, bondReq.UserId, bondReq.OrganizationId)

	if err == nil {
		return bondId, nil
	}

	pgErr := &pgconn.PgError{}
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return "", &domain.HTTPError{Cause: err, Reason: "such employee or company does not exist", Status: domain.BadRequestCode}
		}
	}

	return "", &domain.HTTPError{Cause: err, Reason: "could not bond", Status: domain.ServerFailureCode}
}
