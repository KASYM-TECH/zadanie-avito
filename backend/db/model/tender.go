//nolint:gochecknoglobals
package model

import "time"

type TenderServiceType string

const (
	TenderServiceTypeConstruction TenderServiceType = "Construction"
	TenderServiceTypeDelivery     TenderServiceType = "Delivery"
	TenderServiceTypeManufacture  TenderServiceType = "Manufacture"
)

type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"
	TenderStatusPublished TenderStatus = "Published"
	TenderStatusClosed    TenderStatus = "Closed"
)

type Tender struct {
	Id             string            `db:"id"`
	Name           string            `db:"name"`
	Description    string            `db:"description"`
	ServiceType    TenderServiceType `db:"service_type"`
	Status         TenderStatus      `db:"status"`
	OrganizationId string            `db:"organization_id"`
	Version        int               `db:"version"`
	CreatedAt      time.Time         `db:"created_at"`
	UserId         string            `db:"user_id"`
}
