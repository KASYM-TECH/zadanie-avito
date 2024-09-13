package domain

import (
	"avito/db/model"
	"time"
)

type CreateTenderReq struct {
	Name            string                  `validate:"required,lte=100" json:"name"`
	Description     string                  `validate:"required,lte=500" json:"description"`
	ServiceType     model.TenderServiceType `json:"serviceType"`
	Status          model.TenderStatus      `json:"status"`
	OrganizationId  string                  `validate:"required,lte=100" json:"organizationId"`
	CreatorUsername string                  `json:"creatorUsername"`
}

type CreateTenderResp struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ServiceType model.TenderServiceType `json:"serviceType"`
	Status      model.TenderStatus      `json:"status"`
	Version     int                     `json:"version"`
	CreatedAt   time.Time               `json:"createdAt"`
}

type SetStatusTenderResp struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ServiceType model.TenderServiceType `json:"serviceType"`
	Status      model.TenderStatus      `json:"status"`
	Version     int                     `json:"version"`
	CreatedAt   time.Time               `json:"createdAt"`
}

type GetTendersResp struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ServiceType model.TenderServiceType `json:"serviceType"`
	Status      model.TenderStatus      `json:"status"`
	Version     int                     `json:"version"`
	CreatedAt   time.Time               `json:"createdAt"`
}

type EditTenderReq struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ServiceType model.TenderServiceType `json:"serviceType"`
}

type EditTenderResp struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ServiceType model.TenderServiceType `json:"serviceType"`
	Status      model.TenderStatus      `json:"status"`
	Version     int                     `json:"version"`
	CreatedAt   time.Time               `json:"createdAt"`
}

type RollbackTenderResp struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ServiceType model.TenderServiceType `json:"serviceType"`
	Status      model.TenderStatus      `json:"status"`
	Version     int                     `json:"version"`
	CreatedAt   time.Time               `json:"createdAt"`
}
