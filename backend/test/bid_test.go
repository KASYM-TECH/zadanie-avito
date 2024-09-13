//nolint:wastedassign,ineffassign
package test

import (
	"avito/db/model"
	"avito/domain"
	"avito/test/basic"
	"net/http"
	"testing"
)

func TestBidCreate(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationId:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	bidReq := domain.CreateBidReq{
		Name:        "n1",
		Description: "d1",
		TenderId:    tender.Id,
		AuthorType:  model.BidAuthorTypeOrganization,
		AuthorId:    martinOrg.EmployeeId,
	}
	_, bidResp := basic.CreateBid(test, bidReq)
	test.Assertions.Equal(http.StatusOK, bidResp.StatusCode())

	tenders, myResp := basic.GetBidByUsername(test, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)
}

func TestBidChangeStatus(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationId:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	bidReq := domain.CreateBidReq{
		Name:        "n1",
		Description: "d1",
		TenderId:    tender.Id,
		AuthorType:  model.BidAuthorTypeOrganization,
		AuthorId:    martinOrg.EmployeeId,
	}
	bidResp, resp := basic.CreateBid(test, bidReq)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())

	tenders, myResp := basic.GetBidByUsername(test, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	// MARTIN CAN CHANGE TO CANCELED
	setStatusResp, resp := basic.SetBidStatus(test, bidResp.Id, martinOrg.Username, model.BidStatusCanceled)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(setStatusResp.Status, model.BidStatusCanceled)

	// WE CAN NOT CHANGE TO REJECTED
	_, resp = basic.SetBidStatus(test, bidResp.Id, martinOrg.Username, model.BidStatusRejected)
	test.Assertions.Equal(http.StatusBadRequest, resp.StatusCode())

	// ONLY AUTHOR CAN CHANGE STATUS
	aliceOrg := basic.CreateOrgEmployee(test, "Alice")
	_, resp = basic.SetBidStatus(test, bidResp.Id, aliceOrg.Username, model.BidStatusCanceled)
	test.Assertions.Equal(http.StatusBadRequest, resp.StatusCode())
}

func TestBidEdit(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationId:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	bidReq := domain.CreateBidReq{
		Name:        "n1",
		Description: "d1",
		TenderId:    tender.Id,
		AuthorType:  model.BidAuthorTypeOrganization,
		AuthorId:    martinOrg.EmployeeId,
	}
	bidResp, resp := basic.CreateBid(test, bidReq)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())

	tenders, myResp := basic.GetBidByUsername(test, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	// MARTIN CAN EDIT BECAUSE HE IS AUTHOR
	editBidReq := domain.EditBidReq{
		Name:        "changed name",
		Description: "changed description",
	}
	editResp, resp := basic.EditBid(test, bidResp.Id, martinOrg.Username, editBidReq)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(editResp.Name, editBidReq.Name)

	// ALICE CAN NOT EDIT BECAUSE SHE IS ALICE
	aliceOrg := basic.CreateOrgEmployee(test, "Alice")
	_, resp = basic.EditBid(test, bidResp.Id, aliceOrg.Username, editBidReq)
	test.Assertions.Equal(http.StatusForbidden, resp.StatusCode())
}

func TestBidRollback(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationId:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	bidReq := domain.CreateBidReq{
		Name:        "n1",
		Description: "d1",
		TenderId:    tender.Id,
		AuthorType:  model.BidAuthorTypeOrganization,
		AuthorId:    martinOrg.EmployeeId,
	}
	bidResp, resp := basic.CreateBid(test, bidReq)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())

	tenders, myResp := basic.GetBidByUsername(test, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	editBidReq := domain.EditBidReq{
		Name:        "changed name",
		Description: "changed description",
	}
	editResp, resp := basic.EditBid(test, bidResp.Id, martinOrg.Username, editBidReq)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(editResp.Name, editBidReq.Name)

	// MARTIN CAN ROLL IT BACK BECAUSE HE IS AUTHOR
	rollBackResp, resp := basic.RollbackBid(test, bidResp.Id, martinOrg.Username, "1")
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(bidReq.Name, rollBackResp.Name)

	// ALICE CAN NOT ROLL IT BACK BECAUSE SHE IS ALICE
	aliceOrg := basic.CreateOrgEmployee(test, "Alice")
	_, resp = basic.RollbackBid(test, bidResp.Id, aliceOrg.Username, "1")
	test.Assertions.Equal(http.StatusForbidden, resp.StatusCode())
}

func TestBidSubmitDecision(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationId:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	bidReq := domain.CreateBidReq{
		Name:        "n1",
		Description: "d1",
		TenderId:    tender.Id,
		AuthorType:  model.BidAuthorTypeOrganization,
		AuthorId:    martinOrg.EmployeeId,
	}
	bidResp, resp := basic.CreateBid(test, bidReq)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())

	tenders, myResp := basic.GetBidByUsername(test, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	// MARTIN CAN NOT APPROVE IT BECAUSE BId IS NOT YET PUBLISHED
	submitDecResp, resp := basic.SubmitDecisionBid(test, bidResp.Id, martinOrg.Username, "Approved")
	test.Assertions.Equal(http.StatusBadRequest, resp.StatusCode())

	// MARTIN CAN CHANGE TO PUBLISHED
	setStatusResp, resp := basic.SetBidStatus(test, bidResp.Id, martinOrg.Username, model.BidStatusPublished)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(setStatusResp.Status, model.BidStatusPublished)

	// ALICE CAN NOT APPROVE IT BECAUSE SHE IS NOT FROM ORGANIZATION OF THE TENDER
	aliceOrg := basic.CreateOrgEmployee(test, "Alice")
	submitDecResp, resp = basic.SubmitDecisionBid(test, bidResp.Id, aliceOrg.Username, "Approved")
	test.Assertions.Equal(http.StatusForbidden, resp.StatusCode())

	// MARTIN CAN APPROVE IT BECAUSE HE IS FROM ORGANIZATION OF THE TENDER
	submitDecResp, resp = basic.SubmitDecisionBid(test, bidResp.Id, martinOrg.Username, "Approved")
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(model.BidStatusApproved, submitDecResp.Status)
}

func TestBidSubmitFeedback(t *testing.T) {
	t.Parallel()
	test := basic.InitTest(t)

	martinOrg := basic.CreateOrgEmployee(test, "Martin")

	tenderReq := domain.CreateTenderReq{
		Name:            "n1",
		Description:     "d1",
		ServiceType:     model.TenderServiceTypeConstruction,
		Status:          model.TenderStatusCreated,
		OrganizationId:  martinOrg.OrgId,
		CreatorUsername: martinOrg.Username,
	}
	tender, tenderResp := basic.CreateTender(test, tenderReq)
	test.Assertions.Equal(http.StatusOK, tenderResp.StatusCode())

	bidReq := domain.CreateBidReq{
		Name:        "n1",
		Description: "d1",
		TenderId:    tender.Id,
		AuthorType:  model.BidAuthorTypeOrganization,
		AuthorId:    martinOrg.EmployeeId,
	}
	bidResp, resp := basic.CreateBid(test, bidReq)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())

	tenders, myResp := basic.GetBidByUsername(test, martinOrg.Username, 0, 3)
	test.Assertions.Equal(http.StatusOK, myResp.StatusCode())
	test.Assertions.Len(tenders, 1)

	feedback := "my feedback"
	feedbackResp, resp := basic.SubmitFeedbackBid(test, bidResp.Id, martinOrg.Username, feedback)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(bidReq.Name, feedbackResp.Name)

	// MARTIN CAN SEE REVIEW BECAUSE HE IS THE CREATOR OF THE TENDER AND THE BId
	reviewResp, resp := basic.ReviewBid(test, tender.Id, martinOrg.Username, martinOrg.Username, 0, 1)
	test.Assertions.Equal(http.StatusOK, resp.StatusCode())
	test.Assertions.Equal(feedback, reviewResp[0].Description)

	// ELICE CAN NOT SEE REVIEW BECAUSE SHE IS NOT THE AUTHOR
	aliceOrg := basic.CreateOrgEmployee(test, "Alice")
	_, resp = basic.ReviewBid(test, tender.Id, aliceOrg.Username, martinOrg.Username, 0, 1)
	test.Assertions.Equal(http.StatusBadRequest, resp.StatusCode())

	// ELICE CAN NOT SEE REVIEW BECAUSE SHE IS NOT THE AUTHOR OF THE TENDER
	_, resp = basic.ReviewBid(test, tender.Id, martinOrg.Username, aliceOrg.Username, 0, 1)
	test.Assertions.Equal(http.StatusForbidden, resp.StatusCode())
}
