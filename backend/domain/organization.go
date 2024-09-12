package domain

import (
	"avito/db/model"
)

type CreateOrganizationReq struct {
	Name        string                 `validate:"required,lte=100" json:"name"`
	Description string                 `validate:"required,lte=500" json:"description"`
	Type        model.OrganizationType `validate:"required" json:"type"`
}

type BondReq struct {
	OrganizationId string `validate:"required" json:"organizationId"`
	UserId         string `validate:"required" json:"userId"`
}
