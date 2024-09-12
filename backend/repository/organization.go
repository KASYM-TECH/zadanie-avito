package repository

import (
	"avito/db"
	"avito/db/model"
	"avito/log"
	"context"
	"github.com/pkg/errors"
)

type OrganizationRep struct {
	cli    db.DB
	logger log.Logger
}

func NewOrganizationRep(logger log.Logger, cli db.DB) *OrganizationRep {
	return &OrganizationRep{
		logger: logger, cli: cli,
	}
}

func (rep *OrganizationRep) Insert(ctx context.Context, org *model.Organization) (string, error) {
	var id string
	err := rep.cli.SelectRow(ctx, &id,
		"insert into organization(name, description, type) values ($1, $2, $3) returning id",
		org.Name, org.Description, org.Type)

	if err != nil {
		return id, errors.WithMessage(err, "Repository.Organization.Insert with name: "+org.Name)
	}

	return id, nil
}

func (rep *OrganizationRep) MakeResponsible(ctx context.Context, empId, orgId string) (string, error) {
	var id string
	err := rep.cli.SelectRow(ctx, &id,
		"insert into organization_responsible(organization_id, user_id) values ($1, $2) returning id",
		orgId, empId)

	if err != nil {
		return id, errors.WithMessage(err, "Repository.Organization.MakeResponsible with org id: "+orgId)
	}

	return id, nil
}

func (rep *OrganizationRep) EmpBelongs(ctx context.Context, empId, orgId string) (bool, error) {
	var belongs bool
	err := rep.cli.SelectRow(ctx, &belongs,
		`SELECT EXISTS(SELECT 1 from organization_responsible where organization_id = $1 and user_id = $2)`,
		orgId, empId)

	if err != nil {
		return false, errors.WithMessage(err, "Repository.Org.EmpBelongs with orgId: "+orgId)
	}

	return belongs, nil
}
