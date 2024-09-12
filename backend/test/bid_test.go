package test

import (
	"avito/db/model"
	"avito/domain"
	"avito/test/basic"
	"net/http"
	"testing"
)

func TestBid(t *testing.T) {
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

	bidReq := domain.CreateBidReq{
		Name:        "n1",
		Description: "d1",
		TenderID:    tender.Id,
		AuthorType:  model.BidAuthorTypeOrganization,
		AuthorID:    martinOrg.EmployeeId,
	}
	_, bidResp := basic.CreateBid(test, martinOrg.Token, bidReq)
	test.Assertions.Equal(http.StatusOK, bidResp.StatusCode())

	tenders, myResp := basic.GetBidByUsername(test, martinOrg.Token, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)
}
