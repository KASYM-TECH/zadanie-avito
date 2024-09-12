//nolint:cyclop,gosimple
package controllers

import (
	"avito/db/model"
	"avito/domain"
	"avito/log"
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"strconv"
)

type BidService interface {
	Create(ctx context.Context, tender *model.Bid) (*model.Bid, error)
	GetByUserId(ctx context.Context, offset, limit int, userId string) ([]model.Bid, error)
	GetByTenderId(ctx context.Context, offset, limit int, tenderId string) ([]model.Bid, error)
	GetStatus(ctx context.Context, bidId, username string) (string, error)
	SetStatus(ctx context.Context, bidId, username, status string) (*model.Bid, error)
	Edit(ctx context.Context, userId string, bid *model.Bid) (*model.Bid, error)
	SubmitDecision(ctx context.Context, userId, bidId string, decision model.Decision) (*model.Bid, error)
	SubmitFeedback(ctx context.Context, feedback *model.Feedback) (*model.Bid, error)
	Rollback(ctx context.Context, userId, tenderId string, version int) (*model.Bid, error)
	Reviews(ctx context.Context, requesterName, authorName, tenderId string, offset, limit int) ([]model.Feedback, error)
}

type BidController struct {
	log        log.Logger
	bidService BidService
}

func NewBidController(log log.Logger, bidService BidService) *BidController {
	return &BidController{log: log, bidService: bidService}
}

func (b *BidController) Create(ctx context.Context, createReq domain.CreateBidReq, rd domain.RequestData) (*domain.CreateBidResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "tender_id", createReq.TenderID)
	b.log.Info(ctx, "bid create handler")

	bidMod := &model.Bid{
		Name:        createReq.Name,
		Description: createReq.Description,
		TenderID:    createReq.TenderID,
		AuthorType:  createReq.AuthorType,
		AuthorID:    createReq.AuthorID,
	}
	bid, err := b.bidService.Create(ctx, bidMod)
	if err == nil {
		resp := &domain.CreateBidResp{
			Id:          bid.ID,
			Name:        bid.Name,
			Description: bid.Description,
			AuthorType:  bid.AuthorType,
			AuthorID:    bid.AuthorID,
			Version:     bid.Version,
			CreatedAt:   bid.CreatedAt,
		}
		return resp, nil
	}

	return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
}

func (b *BidController) GetByUsername(ctx context.Context, rd domain.RequestData) ([]domain.GetBidResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	b.log.Info(ctx, "bid Get handler")

	var (
		offset, limit int
		err           error
	)

	username, ok := rd.Request.URL.Query()["username"]
	if !ok || len(username) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	if offsetStr, ok := rd.Request.URL.Query()["offset"]; ok {
		if offset, err = strconv.Atoi(offsetStr[0]); err != nil && offset < 0 {
			offset = 0
		}
	}

	if limitStr, ok := rd.Request.URL.Query()["limit"]; ok {
		if limit, err = strconv.Atoi(limitStr[0]); err != nil && limit < 0 {
			limit = 0
		}
	}

	bids, err := b.bidService.GetByUserId(ctx, offset, limit, rd.UserId)
	resp := make([]domain.GetBidResp, len(bids))
	for i := range bids {
		resp[i] = domain.GetBidResp{
			Id:         bids[i].ID,
			Name:       bids[i].Name,
			Status:     bids[i].Status,
			AuthorType: bids[i].AuthorType,
			AuthorID:   bids[i].AuthorID,
			CreatedAt:  bids[i].CreatedAt,
			Version:    bids[i].Version,
		}
	}

	if err == nil {
		return resp, nil
	}

	return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
}

func (b *BidController) GetByTenderId(ctx context.Context, rd domain.RequestData) ([]domain.GetBidResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	b.log.Info(ctx, "bid GetByTenderId handler")

	tenderId, ok := mux.Vars(rd.Request)["tenderId"]
	if !ok || len(tenderId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required", Status: domain.BadRequestCode}
	}

	username, ok := rd.Request.URL.Query()["username"]
	if !ok || len(username) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	var (
		offset, limit int
		err           error
	)

	if offsetStr, ok := rd.Request.URL.Query()["offset"]; ok {
		if offset, err = strconv.Atoi(offsetStr[0]); err != nil && offset < 0 {
			offset = 0
		}
	}

	if limitStr, ok := rd.Request.URL.Query()["limit"]; ok {
		if limit, err = strconv.Atoi(limitStr[0]); err != nil && limit < 0 {
			limit = 0
		}
	}

	bids, err := b.bidService.GetByTenderId(ctx, offset, limit, tenderId)
	resp := make([]domain.GetBidResp, len(bids))
	for i := range bids {
		resp[i] = domain.GetBidResp{
			Id:         bids[i].ID,
			Name:       bids[i].Name,
			Status:     bids[i].Status,
			AuthorType: bids[i].AuthorType,
			AuthorID:   bids[i].AuthorID,
			CreatedAt:  bids[i].CreatedAt,
			Version:    bids[i].Version,
		}
	}

	if err == nil {
		return resp, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) GetStatus(ctx context.Context, rd domain.RequestData) (string, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	b.log.Info(ctx, "bid GetStatus handler")

	bidId, ok := mux.Vars(rd.Request)["bidId"]
	if !ok || len(bidId) == 0 {
		return "", &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	username, ok := rd.Request.URL.Query()["username"]
	if !ok || len(username) == 0 {
		return "", &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return "", &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	status, err := b.bidService.GetStatus(ctx, bidId, username[0])
	if err == nil {
		return status, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return "", &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	default:
		return "", &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) SetStatus(ctx context.Context, rd domain.RequestData) (*domain.SetStatusBidResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_id", rd.UserId)
	b.log.Info(ctx, "bid SetStatus handler")

	bidId, ok := mux.Vars(rd.Request)["bidId"]
	if !ok || len(bidId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	username, ok := rd.Request.URL.Query()["username"]
	if !ok || len(username) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	status, ok := rd.Request.URL.Query()["status"]
	if !ok || len(username) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "status is required query", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.SetStatus(ctx, bidId, rd.UserId, status[0])
	if err == nil {
		resp := &domain.SetStatusBidResp{
			Id:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			CreatedAt:  bid.CreatedAt,
			Version:    bid.Version,
		}
		return resp, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrNotBidAuthor):
		return nil, &domain.HTTPError{Cause: err, Reason: "You must be the author to edit status", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrForbiddenApproval):
		return nil, &domain.HTTPError{Cause: err, Reason: "You can not self approve", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) Edit(ctx context.Context, req domain.EditBidReq, rd domain.RequestData) (*domain.EditBidResp, *domain.HTTPError) {
	bidId, ok := mux.Vars(rd.Request)["bidId"]
	b.log.Info(ctx, "bid Edit handler")

	if !ok || len(bidId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required query", Status: domain.BadRequestCode}
	}

	var username []string
	if username, ok = rd.Request.URL.Query()["username"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	bidToUpd := &model.Bid{
		ID:          bidId,
		Name:        req.Name,
		Description: req.Description,
	}

	bid, err := b.bidService.Edit(ctx, rd.UserId, bidToUpd)
	if err == nil {
		bidDom := &domain.EditBidResp{
			Id:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			CreatedAt:  bid.CreatedAt,
			Version:    bid.Version,
		}
		return bidDom, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrNotBidAuthor):
		return nil, &domain.HTTPError{Cause: err, Reason: "You must be the author to edit", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) SubmitDecision(ctx context.Context, rd domain.RequestData) (*domain.SubmitDesBidResp, *domain.HTTPError) {
	bidId, _ := mux.Vars(rd.Request)["bidId"]
	b.log.Info(ctx, "bid SubmitDecision handler")

	var (
		username []string
		ok       bool
	)
	if username, ok = rd.Request.URL.Query()["username"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	var decisionStr []string
	if decisionStr, ok = rd.Request.URL.Query()["decision"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "decision is required query", Status: domain.BadRequestCode}
	}

	decision := model.Decision(decisionStr[0])
	if decision != model.Approved && decision != model.Rejected {
		return nil, &domain.HTTPError{Cause: nil, Reason: "invalid decision", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.SubmitDecision(ctx, rd.UserId, bidId, decision)
	if err == nil {
		bidDom := &domain.SubmitDesBidResp{
			Id:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			CreatedAt:  bid.CreatedAt,
			Version:    bid.Version,
		}
		return bidDom, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrBidIsNotPublished):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid is not published", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrTenderIsNotPublished):
		return nil, &domain.HTTPError{Cause: err, Reason: "Tender is not published", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "You must be from tender's organization to submit decision", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) SubmitFeedback(ctx context.Context, rd domain.RequestData) (*domain.FeedbackBidResp, *domain.HTTPError) {
	bidId, ok := mux.Vars(rd.Request)["bidId"]
	b.log.Info(ctx, "bid Feedback handler")

	if !ok || len(bidId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required query", Status: domain.BadRequestCode}
	}

	var username []string
	if username, ok = rd.Request.URL.Query()["username"]; !ok || len(username) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	var feedbackStr []string
	if feedbackStr, ok = rd.Request.URL.Query()["feedback"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "feedback is required query", Status: domain.BadRequestCode}
	}

	feedback := &model.Feedback{
		Content:  feedbackStr[0],
		BidID:    bidId,
		AuthorID: rd.UserId,
	}
	bid, err := b.bidService.SubmitFeedback(ctx, feedback)
	if err == nil {
		bidDom := &domain.FeedbackBidResp{
			Id:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			CreatedAt:  bid.CreatedAt,
			Version:    bid.Version,
		}
		return bidDom, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "You are not responsible for tender organization", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) Rollback(ctx context.Context, rd domain.RequestData) (*domain.RollbackBidResp, *domain.HTTPError) {
	b.log.Info(ctx, "tender Rollback handler")

	bidId, ok := mux.Vars(rd.Request)["bidId"]
	if !ok || len(bidId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required query", Status: domain.BadRequestCode}
	}
	versionStr, ok := mux.Vars(rd.Request)["version"]
	if !ok || len(versionStr) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version is required query", Status: domain.BadRequestCode}
	}

	var username []string
	if username, ok = rd.Request.URL.Query()["username"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if rd.Claims.Username != username[0] {
		return nil, &domain.HTTPError{
			Cause:  domain.ErrClient,
			Reason: "username does not match with token's username",
			Status: domain.UnauthorizedCode}
	}

	var (
		version int
		err     error
	)
	if version, err = strconv.Atoi(versionStr); err != nil || version < 1 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version must be positive integer", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.Rollback(ctx, rd.UserId, bidId, version)
	if err == nil {
		resp := &domain.RollbackBidResp{
			Id:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			CreatedAt:  bid.CreatedAt,
			Version:    bid.Version,
		}
		return resp, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrNotBidAuthor):
		return nil, &domain.HTTPError{Cause: err, Reason: "You must be the bid author to roll it back", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) Reviews(ctx context.Context, rd domain.RequestData) ([]domain.ReviewResp, *domain.HTTPError) {
	b.log.Info(ctx, "tender Reviews handler")

	tenderId, ok := mux.Vars(rd.Request)["tenderId"]
	if !ok || len(tenderId) == 0 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required", Status: domain.BadRequestCode}
	}

	var authorUsername []string
	if authorUsername, ok = rd.Request.URL.Query()["authorUsername"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "authorUsername is required query", Status: domain.BadRequestCode}
	}

	var requesterUsername []string
	if requesterUsername, ok = rd.Request.URL.Query()["requesterUsername"]; !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "requesterUsername is required query", Status: domain.BadRequestCode}
	}

	var (
		offset, limit int
		err           error
	)

	if offsetStr, ok := rd.Request.URL.Query()["offset"]; ok {
		if offset, err = strconv.Atoi(offsetStr[0]); err != nil && offset < 0 {
			offset = 0
		}
	}

	if limitStr, ok := rd.Request.URL.Query()["limit"]; ok {
		if limit, err = strconv.Atoi(limitStr[0]); err != nil && limit < 0 {
			limit = 0
		}
	}

	reviews, err := b.bidService.Reviews(ctx, requesterUsername[0], authorUsername[0], tenderId, offset, limit)
	if err == nil {
		resp := make([]domain.ReviewResp, len(reviews))
		for i := range reviews {
			resp[i] = domain.ReviewResp{
				Id:          reviews[i].ID,
				Description: reviews[i].Content,
				CreatedAt:   reviews[i].CreatedAt,
			}
		}
		return resp, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "Bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrAuthorIsIncorrect):
		return nil, &domain.HTTPError{Cause: err, Reason: "Specified author is not author of the tender", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "You are not responsible for tender organization", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server unavailable", Status: domain.ServerFailureCode}
	}
}
