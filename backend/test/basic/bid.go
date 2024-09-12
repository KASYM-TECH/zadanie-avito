package basic

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/txix-open/isp-kit/http/httpcli"
)

func SetBidStatus(test *Test, token, bidId, username string, status model.BidStatus) (domain.SetStatusBidResp, *httpcli.Response) {
	assert := test.Assertions

	var bidResp domain.SetStatusBidResp
	resp, err := test.Cli.Put(test.URL+"/api/bids/"+bidId+"/status").
		QueryParams(map[string]any{"username": username, "status": status}).
		JsonResponseBody(&bidResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return bidResp, resp
}

func CreateBid(test *Test, token string, bidReq domain.CreateBidReq) (domain.CreateBidResp, *httpcli.Response) {
	assert := test.Assertions

	var bid domain.CreateBidResp
	resp, err := test.Cli.Post(test.URL+"/api/bids/new").
		JsonRequestBody(&bidReq).
		JsonResponseBody(&bid).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return bid, resp
}

func GetBidByUsername(test *Test, token, username string, offset, limit int) ([]domain.GetBidResp, *httpcli.Response) {
	assert := test.Assertions

	var bidsResp []domain.GetBidResp
	resp, err := test.Cli.Get(test.URL+"/api/bids/my").
		QueryParams(map[string]any{"username": username, "offset": offset, "limit": limit}).
		JsonResponseBody(&bidsResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return bidsResp, resp
}

func EditBid(test *Test, token, bidId, username string, req domain.EditBidReq) (domain.EditBidResp, *httpcli.Response) {
	assert := test.Assertions

	var tenderBidResp domain.EditBidResp
	resp, err := test.Cli.Patch(test.URL+"/api/bids/"+username+"/edit").
		JsonRequestBody(req).
		QueryParams(map[string]any{"username": username}).
		JsonResponseBody(&tenderBidResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return tenderBidResp, resp
}

func RollbackBid(test *Test, token, bidId, username, version string) (domain.RollbackTenderResp, *httpcli.Response) {
	assert := test.Assertions

	var rollbackResp domain.RollbackTenderResp
	resp, err := test.Cli.Put(test.URL+"/api/bids/"+bidId+"/rollback/"+version).
		QueryParams(map[string]any{"username": username}).
		JsonResponseBody(&rollbackResp).
		Header("Authorization", token).
		Do(context.Background())

	assert.NoError(err)

	return rollbackResp, resp
}
