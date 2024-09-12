package test

import (
	"avito/db/model"
	"avito/domain"
	"avito/test/basic"
	"net/http"
	"testing"
)

func TestTender(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationID:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	_, tenderResp := basic.CreateTender(test, martinOrg.Token, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	tenders, myResp := basic.GetTenderByUsername(test, martinOrg.Token, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	aliceOrg := basic.CreateOrgEmployee(test, "Alice")

	tenderForbiddenReq := domain.CreateTenderReq{
		Name:        "n1",
		Description: "d1",
		ServiceType: model.TenderServiceTypeConstruction,
		Status:      model.TenderStatusCreated,
		// ATTENTION HERE
		OrganizationID:  martinOrg.OrgId,
		CreatorUsername: aliceOrg.Username,
	}
	_, tenderForbiddenResp := basic.CreateTender(test, aliceOrg.Token, tenderForbiddenReq)
	test.Assertions.Equal(http.StatusForbidden, tenderForbiddenResp.StatusCode())
}

func TestTenderStatus(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationID:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, martinOrg.Token, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	tenders, myResp := basic.GetTenderByUsername(test, martinOrg.Token, martinOrg.Username, 0, 1)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	tenderSetStatus, tenderSetStatusResp := basic.SetTenderStatus(test, martinOrg.Token, tender.Id, martinOrg.Username, model.TenderStatusClosed)
	test.Assertions.Equal(http.StatusOK, tenderSetStatusResp.StatusCode())
	test.Assertions.Equal(model.TenderStatusClosed, tenderSetStatus.Status)

	tenderGetStatus, tenderSetStatusResp := basic.GetTenderStatus(test, martinOrg.Token, tender.Id, martinOrg.Username)
	test.Assertions.Equal(http.StatusOK, tenderSetStatusResp.StatusCode())
	test.Assertions.Equal(model.TenderStatusClosed, tenderGetStatus)
}

func TestTenderEdit(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationID:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, martinOrg.Token, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	tenders, myResp := basic.GetTenderByUsername(test, martinOrg.Token, martinOrg.Username, 0, 1)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	editReq := domain.EditTenderReq{
		Name:        "changed",
		Description: "changed",
		ServiceType: model.TenderServiceTypeConstruction,
	}
	editTenderResp, myResp := basic.EditTender(test, martinOrg.Token, tender.Id, martinOrg.Username, editReq)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Equal(editTenderResp.Name, editReq.Name)
	test.Assertions.Equal(editTenderResp.Description, editReq.Description)
	test.Assertions.Equal(editTenderResp.ServiceType, editReq.ServiceType)

	aliceOrg := basic.CreateOrgEmployee(test, "Alice")
	editTenderResp, myResp = basic.EditTender(test, aliceOrg.Token, tender.Id, aliceOrg.Username, editReq)
	test.Assertions.Equal(http.StatusForbidden, myResp.StatusCode())
}

func TestTenderRollback(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationID:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, martinOrg.Token, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	tenders, myResp := basic.GetTenderByUsername(test, martinOrg.Token, martinOrg.Username, 0, 1)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	editReq := domain.EditTenderReq{
		Name:        "changed",
		Description: "changed",
		ServiceType: model.TenderServiceTypeConstruction,
	}
	editTenderResp, myResp := basic.EditTender(test, martinOrg.Token, tender.Id, martinOrg.Username, editReq)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Equal(editTenderResp.Name, editReq.Name)
	test.Assertions.Equal(editTenderResp.Description, editReq.Description)
	test.Assertions.Equal(editTenderResp.ServiceType, editReq.ServiceType)

	rollbackResp, resp := basic.RollbackTender(test, martinOrg.Token, tender.Id, martinOrg.Username, "1")
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(tenderReq.Name, rollbackResp.Name)
	test.Assertions.Equal(tenderReq.Description, rollbackResp.Description)
	test.Assertions.Equal(tenderReq.ServiceType, rollbackResp.ServiceType)

	aliceOrg := basic.CreateOrgEmployee(test, "Alice")
	rollbackResp, resp = basic.RollbackTender(test, aliceOrg.Token, tender.Id, aliceOrg.Username, "1")
	test.Assertions.Equal(http.StatusForbidden, resp.StatusCode())
}
