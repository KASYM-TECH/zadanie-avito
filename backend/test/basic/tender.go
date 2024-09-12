package basic

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/txix-open/isp-kit/http/httpcli"
)

func CreateTender(test *Test, token string, tender domain.CreateTenderReq) (domain.CreateTenderResp, *httpcli.Response) {
	assert := test.Assertions

	var tenderResp domain.CreateTenderResp
	resp, err := test.Cli.Post(test.URL+"/api/tenders/new").
		JsonRequestBody(&tender).
		JsonResponseBody(&tenderResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return tenderResp, resp
}

func SetTenderStatus(test *Test, token, tenderId, username string, status model.TenderStatus) (domain.SetStatusTenderResp, *httpcli.Response) {
	assert := test.Assertions

	var tender domain.SetStatusTenderResp
	resp, err := test.Cli.Put(test.URL+"/api/tenders/"+tenderId+"/status").
		QueryParams(map[string]any{"username": username, "status": status}).
		JsonResponseBody(&tender).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return tender, resp
}

func GetTenderStatus(test *Test, token, tenderId, username string) (model.TenderStatus, *httpcli.Response) {
	assert := test.Assertions

	var status model.TenderStatus
	resp, err := test.Cli.Get(test.URL+"/api/tenders/"+tenderId+"/status").
		QueryParams(map[string]any{"username": username}).
		JsonResponseBody(&status).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return status, resp
}

func GetTenderByUsername(test *Test, token, username string, offset, limit int) ([]domain.GetTendersResp, *httpcli.Response) {
	assert := test.Assertions

	var tendersResp []domain.GetTendersResp
	resp, err := test.Cli.Get(test.URL+"/api/tenders/my").
		QueryParams(map[string]any{"username": username, "offset": offset, "limit": limit}).
		JsonResponseBody(&tendersResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return tendersResp, resp
}

func EditTender(test *Test, token, tenderId, username string, req domain.EditTenderReq) (domain.EditTenderResp, *httpcli.Response) {
	assert := test.Assertions

	var tenderEditResp domain.EditTenderResp
	resp, err := test.Cli.Patch(test.URL+"/api/tenders/"+tenderId+"/edit").
		JsonRequestBody(req).
		QueryParams(map[string]any{"username": username}).
		JsonResponseBody(&tenderEditResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return tenderEditResp, resp
}

func RollbackTender(test *Test, token, tenderId, username, version string) (domain.RollbackTenderResp, *httpcli.Response) {
	assert := test.Assertions

	var tenderRollbackResp domain.RollbackTenderResp
	resp, err := test.Cli.Put(test.URL+"/api/tenders/"+tenderId+"/rollback/"+version).
		QueryParams(map[string]any{"username": username}).
		JsonResponseBody(&tenderRollbackResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return tenderRollbackResp, resp
}
