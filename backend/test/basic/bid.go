package basic

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/txix-open/isp-kit/http/httpcli"
)

func SetBidStatus(test *Test, bidId, username string, status model.BidStatus) (domain.SetStatusBidResp, *httpcli.Response) {
	assert := test.Assertions

	var bidResp domain.SetStatusBidResp
	resp, err := test.Cli.Put(test.URL + "/api/bids/" + bidId + "/status").
		QueryParams(map[string]any{"username": username, "status": status}).
		JsonResponseBody(&bidResp).
		Do(context.Background())

	assert.NoError(err)

	return bidResp, resp
}

func CreateBid(test *Test, bidReq domain.CreateBidReq) (domain.CreateBidResp, *httpcli.Response) {
	assert := test.Assertions

	var bid domain.CreateBidResp
	resp, err := test.Cli.Post(test.URL + "/api/bids/new").
		JsonRequestBody(&bidReq).
		JsonResponseBody(&bid).
		Do(context.Background())

	assert.NoError(err)

	return bid, resp
}

func GetBidByUsername(test *Test, username string, offset, limit int) ([]domain.GetBidResp, *httpcli.Response) {
	assert := test.Assertions

	var bidsResp []domain.GetBidResp
	resp, err := test.Cli.Get(test.URL + "/api/bids/my").
		QueryParams(map[string]any{"username": username, "offset": offset, "limit": limit}).
		JsonResponseBody(&bidsResp).
		Do(context.Background())

	assert.NoError(err)

	return bidsResp, resp
}

func EditBid(test *Test, bidId, username string, req domain.EditBidReq) (domain.EditBidResp, *httpcli.Response) {
	assert := test.Assertions

	var tenderBidResp domain.EditBidResp
	resp, err := test.Cli.Patch(test.URL + "/api/bids/" + bidId + "/edit").
		JsonRequestBody(req).
		QueryParams(map[string]any{"username": username}).
		JsonResponseBody(&tenderBidResp).
		Do(context.Background())

	assert.NoError(err)

	return tenderBidResp, resp
}

func RollbackBid(test *Test, bidId, username, version string) (domain.RollbackTenderResp, *httpcli.Response) {
	assert := test.Assertions

	var rollbackResp domain.RollbackTenderResp
	resp, err := test.Cli.Put(test.URL + "/api/bids/" + bidId + "/rollback/" + version).
		QueryParams(map[string]any{"username": username}).
		JsonResponseBody(&rollbackResp).
		Do(context.Background())

	assert.NoError(err)

	return rollbackResp, resp
}

func SubmitDecisionBid(test *Test, bidId, username, decision string) (domain.SubmitDecisionBidResp, *httpcli.Response) {
	assert := test.Assertions

	var submitDecisionResp domain.SubmitDecisionBidResp
	resp, err := test.Cli.Put(test.URL + "/api/bids/" + bidId + "/submit_decision").
		QueryParams(map[string]any{"username": username, "decision": decision}).
		JsonResponseBody(&submitDecisionResp).
		Do(context.Background())

	assert.NoError(err)

	return submitDecisionResp, resp
}

func SubmitFeedbackBid(test *Test, bidId, username, feedback string) (domain.FeedbackBidResp, *httpcli.Response) {
	assert := test.Assertions

	var feedbackBidResp domain.FeedbackBidResp
	resp, err := test.Cli.Put(test.URL + "/api/bids/" + bidId + "/feedback").
		QueryParams(map[string]any{"username": username, "feedback": feedback}).
		JsonResponseBody(&feedbackBidResp).
		Do(context.Background())

	assert.NoError(err)

	return feedbackBidResp, resp
}

func ReviewBid(test *Test, tenderId, authorUsername, requesterUsername string, offset, limit int) ([]domain.ReviewResp, *httpcli.Response) {
	assert := test.Assertions

	var reviewResp []domain.ReviewResp
	resp, err := test.Cli.Get(test.URL + "/api/bids/" + tenderId + "/reviews").
		QueryParams(map[string]any{"limit": limit,
			"offset":            offset,
			"authorUsername":    authorUsername,
			"requesterUsername": requesterUsername}).
		JsonResponseBody(&reviewResp).
		Do(context.Background())

	assert.NoError(err)

	return reviewResp, resp
}
