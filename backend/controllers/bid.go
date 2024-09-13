//nolint:cyclop,gosimple
package controllers

import (
	"avito/domain"
	"avito/log"
	"context"
	"github.com/pkg/errors"
	"strconv"
)

type BidService interface {
	Create(ctx context.Context, tender *domain.CreateBidReq) (*domain.CreateBidResp, error)
	GetByUsername(ctx context.Context, offset, limit int, username string) ([]domain.GetBidResp, error)
	GetByTenderId(ctx context.Context, offset, limit int, tenderId string) ([]domain.GetBidResp, error)
	GetStatus(ctx context.Context, bidId, username string) (string, error)
	SetStatus(ctx context.Context, bidId, username, status string) (*domain.SetStatusBidResp, error)
	Edit(ctx context.Context, username, bidId string, bid *domain.EditBidReq) (*domain.EditBidResp, error)
	SubmitDecision(ctx context.Context, username, bidId string, decision string) (*domain.SubmitDecisionBidResp, error)
	SubmitFeedback(ctx context.Context, content, bidId, authorUsername string) (*domain.FeedbackBidResp, error)
	Rollback(ctx context.Context, username, tenderId string, version int) (*domain.RollbackBidResp, error)
	Reviews(ctx context.Context, requesterName, authorName, tenderId string, offset, limit int) ([]domain.ReviewResp, error)
}

type BidController struct {
	log        log.Logger
	bidService BidService
}

func NewBidController(log log.Logger, bidService BidService) *BidController {
	return &BidController{log: log, bidService: bidService}
}

func (b *BidController) Create(ctx context.Context, createReq domain.CreateBidReq, rd domain.RequestData) (*domain.CreateBidResp, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "tenderId", createReq.TenderId)
	b.log.Info(ctx, "bid create handler")

	bid, err := b.bidService.Create(ctx, &createReq)

	if err == nil {
		return bid, nil
	}

	return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
}

func (b *BidController) GetByUsername(ctx context.Context, rd domain.RequestData) ([]domain.GetBidResp, *domain.HTTPError) {
	var (
		username            string
		ok                  bool
		offsetStr, limitStr string
		offset, limit       int
		err                 error
	)

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "username", username)
	b.log.Info(ctx, "bid GetByUsername handler")

	offsetStr, _ = ExtractQuery(rd.Request, "offset", "0")
	if offset, err = strconv.Atoi(offsetStr); err != nil && offset < 0 {
		offset = 0
	}

	limitStr, _ = ExtractQuery(rd.Request, "limit", "0")
	if limit, err = strconv.Atoi(limitStr); err != nil && limit < 0 {
		limit = 0
	}

	bids, err := b.bidService.GetByUsername(ctx, offset, limit, username)

	if err == nil {
		return bids, nil
	}

	switch {
	case errors.Is(err, domain.ErrUserWithNameNotFound):
		return nil, &domain.HTTPError{Cause: err, Reason: "user was not found", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) GetByTenderId(ctx context.Context, rd domain.RequestData) ([]domain.GetBidResp, *domain.HTTPError) {
	var (
		username, tenderId  string
		offsetStr, limitStr string
		offset, limit       int
		ok                  bool
		err                 error
	)

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "username", username)
	b.log.Info(ctx, "bid GetByTenderId handler")

	if tenderId, ok = ExtractParam(rd.Request, "tenderId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required", Status: domain.BadRequestCode}
	}

	offsetStr, _ = ExtractQuery(rd.Request, "offset", "0")
	if offset, err = strconv.Atoi(offsetStr); err != nil && offset < 0 {
		offset = 0
	}

	limitStr, _ = ExtractQuery(rd.Request, "limit", "0")
	if limit, err = strconv.Atoi(limitStr); err != nil && limit < 0 {
		limit = 0
	}

	bids, err := b.bidService.GetByTenderId(ctx, offset, limit, tenderId)

	if err == nil {
		return bids, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) GetStatus(ctx context.Context, rd domain.RequestData) (string, *domain.HTTPError) {
	var (
		username, bidId string
		ok              bool
	)

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return "", &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "username", username)
	b.log.Info(ctx, "bid GetStatus handler")

	if bidId, ok = ExtractParam(rd.Request, "bidId", ""); !ok {
		return "", &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	status, err := b.bidService.GetStatus(ctx, bidId, username)
	if err == nil {
		return status, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return "", &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	default:
		return "", &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) SetStatus(ctx context.Context, rd domain.RequestData) (*domain.SetStatusBidResp, *domain.HTTPError) {
	var (
		username, bidId, status string
		ok                      bool
	)

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "username", username)
	b.log.Info(ctx, "bid SetStatus handler")

	if bidId, ok = ExtractParam(rd.Request, "bidId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	if status, ok = ExtractQuery(rd.Request, "status", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "status is required query", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.SetStatus(ctx, bidId, username, status)
	if err == nil {
		return bid, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrNotBidAuthor):
		return nil, &domain.HTTPError{Cause: err, Reason: "you must be the author to edit status", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrForbiddenApproval):
		return nil, &domain.HTTPError{Cause: err, Reason: "you can not self approve", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) Edit(ctx context.Context, req domain.EditBidReq, rd domain.RequestData) (*domain.EditBidResp, *domain.HTTPError) {
	var (
		username, bidId string
		ok              bool
	)

	if bidId, ok = ExtractParam(rd.Request, "bidId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "bidId", bidId)
	b.log.Info(ctx, "bid edit handler")

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.Edit(ctx, username, bidId, &req)
	if err == nil {
		return bid, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrNotBidAuthor):
		return nil, &domain.HTTPError{Cause: err, Reason: "you must be the author to edit", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) SubmitDecision(ctx context.Context, rd domain.RequestData) (*domain.SubmitDecisionBidResp, *domain.HTTPError) {
	var (
		username, bidId, decisionStr string
		ok                           bool
	)

	if bidId, ok = ExtractParam(rd.Request, "bidId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "bidId", bidId)
	b.log.Info(ctx, "bid SubmitDecision handler")

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if decisionStr, ok = ExtractQuery(rd.Request, "decision", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "decision is required query", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.SubmitDecision(ctx, username, bidId, decisionStr)
	if err == nil {
		return bid, nil
	}

	switch {
	case errors.Is(err, domain.ErrInvalidDecision):
		return nil, &domain.HTTPError{Cause: err, Reason: "decision is invalid", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrBidIsNotPublished):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid is not published", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrTenderIsNotPublished):
		return nil, &domain.HTTPError{Cause: err, Reason: "tender is not published", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "you must be from tender's organization to submit decision",
			Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) SubmitFeedback(ctx context.Context, rd domain.RequestData) (*domain.FeedbackBidResp, *domain.HTTPError) {
	var (
		username, bidId, feedbackContent string
		ok                               bool
	)

	if bidId, ok = ExtractParam(rd.Request, "bidId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "bidId", bidId)
	b.log.Info(ctx, "bid SubmitFeedback handler")

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	if feedbackContent, ok = ExtractQuery(rd.Request, "feedback", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "feedback is required query", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.SubmitFeedback(ctx, feedbackContent, bidId, username)
	if err == nil {
		return bid, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "you are not responsible for tender organization", Status: domain.BadRequestCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) Rollback(ctx context.Context, rd domain.RequestData) (*domain.RollbackBidResp, *domain.HTTPError) {
	var (
		username, bidId, versionStr string
		ok                          bool
		version                     int
		err                         error
	)

	if bidId, ok = ExtractParam(rd.Request, "bidId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "bidId is required", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "bidId", bidId)
	b.log.Info(ctx, "tender Rollback handler")

	if versionStr, ok = ExtractParam(rd.Request, "version", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version is required", Status: domain.BadRequestCode}
	}

	if version, err = strconv.Atoi(versionStr); err != nil || version < 1 {
		return nil, &domain.HTTPError{Cause: nil, Reason: "version must be positive integer", Status: domain.BadRequestCode}
	}

	if username, ok = ExtractQuery(rd.Request, "username", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "username is required query", Status: domain.BadRequestCode}
	}

	bid, err := b.bidService.Rollback(ctx, username, bidId, version)
	if err == nil {
		return bid, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrNotBidAuthor):
		return nil, &domain.HTTPError{Cause: err, Reason: "you must be the bid author to roll it back", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}

func (b *BidController) Reviews(ctx context.Context, rd domain.RequestData) ([]domain.ReviewResp, *domain.HTTPError) {
	var (
		ok                  bool
		tenderId            string
		authorUsername      string
		requesterUsername   string
		offsetStr, limitStr string
		offset, limit       int
		err                 error
	)

	if tenderId, ok = ExtractParam(rd.Request, "tenderId", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "tenderId is required", Status: domain.BadRequestCode}
	}

	ctx = log.AddKeyVal(ctx, "tenderId", tenderId)
	b.log.Info(ctx, "tender Reviews handler")

	if authorUsername, ok = ExtractQuery(rd.Request, "authorUsername", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "authorUsername is required query", Status: domain.BadRequestCode}
	}

	if requesterUsername, ok = ExtractQuery(rd.Request, "requesterUsername", ""); !ok {
		return nil, &domain.HTTPError{Cause: nil, Reason: "requesterUsername is required query", Status: domain.BadRequestCode}
	}

	offsetStr, _ = ExtractQuery(rd.Request, "offset", "0")
	if offset, err = strconv.Atoi(offsetStr); err != nil && offset < 0 {
		offset = 0
	}

	limitStr, _ = ExtractQuery(rd.Request, "limit", "0")
	if limit, err = strconv.Atoi(limitStr); err != nil && limit < 0 {
		limit = 0
	}

	reviews, err := b.bidService.Reviews(ctx, requesterUsername, authorUsername, tenderId, offset, limit)
	if err == nil {
		return reviews, nil
	}

	switch {
	case errors.Is(err, domain.ErrBidDoesNotExist):
		return nil, &domain.HTTPError{Cause: err, Reason: "bid with this id does not exist", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrAuthorIsIncorrect):
		return nil, &domain.HTTPError{Cause: err, Reason: "specified author is not the author of the tender", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrUserNotResponsible):
		return nil, &domain.HTTPError{Cause: err, Reason: "you are not responsible for tender organization", Status: domain.ForbiddenCode}
	default:
		return nil, &domain.HTTPError{Cause: err, Reason: "server error", Status: domain.ServerFailureCode}
	}
}
