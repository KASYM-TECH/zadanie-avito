package model

import "time"

type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"
	OrganizationTypeLLC OrganizationType = "LLC"
	OrganizationTypeJSC OrganizationType = "JSC"
)

type Organization struct {
	Id          string           `db:"id"`
	Name        string           `db:"name"`
	Description string           `db:"description"`
	Type        OrganizationType `db:"type"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
}

type OrganizationResponsible struct {
	Id             string
	OrganizationId string
	UserId         string
}
