//nolint:interfacebloat
package service

import (
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/pkg/errors"
)

type BidRep interface {
	Insert(ctx context.Context, newTender *model.Bid) (string, error)
	GetById(ctx context.Context, bidId string) (*model.Bid, error)
	GetByUsername(ctx context.Context, offset, limit int, username string) ([]model.Bid, error)
	GetVisibleByTenderId(ctx context.Context, offset, limit int, tenderId string) ([]model.Bid, error)
	GetBidStatus(ctx context.Context, bidId string) (string, error)
	SetBidStatus(ctx context.Context, bidId, status string) error
	UpdateById(ctx context.Context, bid *model.Bid) error
	SetBidStatusIfOpen(ctx context.Context, bidId string, status string) (string, error)
	Rollback(ctx context.Context, bidId string, version int) error
	GetAuthorId(ctx context.Context, bidId string) (string, error)
	GetOrgIdByBidId(ctx context.Context, bidId string) (string, error)
	GetUserIdByName(ctx context.Context, username string) (string, error)
}

type FeedbackRep interface {
	SaveFeedback(ctx context.Context, feedback *model.Feedback) error
	Reviews(ctx context.Context, authorName string, offset, limit int) ([]model.Feedback, error)
}

type DecisionTransaction interface {
	SetBidStatusIfOpen(ctx context.Context, bidId string, status string) (string, error)
	SetTenderStatusIfOpen(ctx context.Context, tenderId, status string) error
}

type TxManager interface {
	DecisionTransaction(ctx context.Context, pTx func(ctx context.Context, tx DecisionTransaction) error) error
}

type BidService struct {
	bidRep      BidRep
	feedbackRep FeedbackRep
	tenderRep   TenderRep
	orgRep      OrganizationRep
	txMan       TxManager
}

func NewBidService(bidRep BidRep, feedbackRep FeedbackRep, tenderRep TenderRep, orgRep OrganizationRep, txMan TxManager) BidService {
	return BidService{bidRep: bidRep, feedbackRep: feedbackRep, tenderRep: tenderRep, orgRep: orgRep, txMan: txMan}
}

func (s BidService) Create(ctx context.Context, bidDom *domain.CreateBidReq) (*domain.CreateBidResp, error) {
	bidMod := &model.Bid{
		Name:        bidDom.Name,
		Description: bidDom.Description,
		TenderId:    bidDom.TenderId,
		AuthorType:  bidDom.AuthorType,
		AuthorId:    bidDom.AuthorId,
	}

	bidId, err := s.bidRep.Insert(ctx, bidMod)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid insert")
	}

	bid, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid insert")
	}

	resp := &domain.CreateBidResp{
		Id:          bid.Id,
		Name:        bid.Name,
		Description: bid.Description,
		AuthorType:  bid.AuthorType,
		AuthorId:    bid.AuthorId,
		Version:     bid.Version,
		CreatedAt:   bid.CreatedAt,
	}

	return resp, nil
}

func (s BidService) GetByUsername(ctx context.Context, offset, limit int, username string) ([]domain.GetBidResp, error) {
	bids, err := s.bidRep.GetByUsername(ctx, offset, limit, username)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.GetBidResp, len(bids))
	for i := range bids {
		resp[i] = domain.GetBidResp{
			Id:         bids[i].Id,
			Name:       bids[i].Name,
			Status:     bids[i].Status,
			AuthorType: bids[i].AuthorType,
			AuthorId:   bids[i].AuthorId,
			CreatedAt:  bids[i].CreatedAt,
			Version:    bids[i].Version,
		}
	}

	return resp, nil
}

func (s BidService) GetByTenderId(ctx context.Context, offset, limit int, tenderId string) ([]domain.GetBidResp, error) {
	bids, err := s.bidRep.GetVisibleByTenderId(ctx, offset, limit, tenderId)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.GetBidResp, len(bids))
	for i := range bids {
		resp[i] = domain.GetBidResp{
			Id:         bids[i].Id,
			Name:       bids[i].Name,
			Status:     bids[i].Status,
			AuthorType: bids[i].AuthorType,
			AuthorId:   bids[i].AuthorId,
			CreatedAt:  bids[i].CreatedAt,
			Version:    bids[i].Version,
		}
	}

	return resp, nil
}

func (s BidService) GetStatus(ctx context.Context, bidId, username string) (string, error) {
	authorId, err := s.bidRep.GetAuthorId(ctx, bidId)
	if err != nil {
		return "", err
	}

	userId, err := s.bidRep.GetUserIdByName(ctx, username)
	if err != nil {
		return "", err
	}

	status, err := s.bidRep.GetBidStatus(ctx, bidId)
	if err != nil {
		return "", err
	}

	if authorId != userId && status == string(model.BidStatusCreated) {
		return "", domain.ErrBidDoesNotExist
	}

	return status, nil
}

func (s BidService) SetStatus(ctx context.Context, bidId, username, status string) (*domain.SetStatusBidResp, error) {
	authorId, err := s.bidRep.GetAuthorId(ctx, bidId)
	if err != nil {
		return nil, err
	}

	userId, err := s.bidRep.GetUserIdByName(ctx, username)
	if err != nil {
		return nil, err
	}

	if authorId != userId {
		return nil, domain.ErrNotBidAuthor
	}

	if status == string(model.BidStatusApproved) || status == string(model.BidStatusRejected) {
		return nil, domain.ErrForbiddenApproval
	}

	err = s.bidRep.SetBidStatus(ctx, bidId, status)
	if err != nil {
		return nil, err
	}

	bid, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, err
	}

	resp := &domain.SetStatusBidResp{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		CreatedAt:  bid.CreatedAt,
		Version:    bid.Version,
	}

	return resp, nil
}

func (s BidService) Edit(ctx context.Context, username, bidId string, editBid *domain.EditBidReq) (*domain.EditBidResp, error) {
	authorId, err := s.bidRep.GetAuthorId(ctx, bidId)
	if err != nil {
		return nil, err
	}

	userId, err := s.bidRep.GetUserIdByName(ctx, username)
	if err != nil {
		return nil, err
	}

	if authorId != userId {
		return nil, domain.ErrNotBidAuthor
	}

	bidToUpd := &model.Bid{
		Id:          bidId,
		Name:        editBid.Name,
		Description: editBid.Description,
	}

	err = s.bidRep.UpdateById(ctx, bidToUpd)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid edit")
	}

	bid, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	bidDom := &domain.EditBidResp{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		CreatedAt:  bid.CreatedAt,
		Version:    bid.Version,
	}

	return bidDom, nil
}

func (s BidService) SubmitDecision(ctx context.Context, username, bidId, decision string) (*domain.SubmitDecisionBidResp, error) {
	orgId, err := s.bidRep.GetOrgIdByBidId(ctx, bidId)
	if err != nil {
		return nil, err
	}

	isResponsible, err := s.orgRep.EmpBelongs(ctx, username, orgId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	if decision != string(model.Approved) && decision != string(model.Rejected) {
		return nil, domain.ErrInvalidDecision
	}

	err = s.txMan.DecisionTransaction(ctx, func(ctx context.Context, tx DecisionTransaction) error {
		var (
			tenderStatus = model.TenderStatusClosed
			bidStatus    = model.BidStatusRejected
		)

		if decision == string(model.Approved) {
			bidStatus = model.BidStatusApproved
		}

		tenderId, err := tx.SetBidStatusIfOpen(ctx, bidId, string(bidStatus))
		if err != nil {
			return domain.ErrBidIsNotPublished
		}

		err = tx.SetTenderStatusIfOpen(ctx, tenderId, string(tenderStatus))
		if err != nil {
			return domain.ErrTenderIsNotPublished
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	bid, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	bidDom := &domain.SubmitDecisionBidResp{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		CreatedAt:  bid.CreatedAt,
		Version:    bid.Version,
	}

	return bidDom, nil
}

func (s BidService) SubmitFeedback(ctx context.Context, content, bidId, authorUsername string) (*domain.FeedbackBidResp, error) {
	orgId, err := s.bidRep.GetOrgIdByBidId(ctx, bidId)
	if err != nil {
		return nil, err
	}

	isResponsible, err := s.orgRep.EmpBelongs(ctx, authorUsername, orgId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	authorId, err := s.bidRep.GetUserIdByName(ctx, authorUsername)
	if err != nil {
		return nil, err
	}

	receiverId, err := s.bidRep.GetAuthorId(ctx, bidId)
	if err != nil {
		return nil, err
	}

	feedback := &model.Feedback{
		Content:    content,
		BidId:      bidId,
		AuthorId:   authorId,
		ReceiverId: receiverId,
	}

	err = s.feedbackRep.SaveFeedback(ctx, feedback)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid submit feedback")
	}

	bid, err := s.bidRep.GetById(ctx, feedback.BidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	bidDom := &domain.FeedbackBidResp{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		CreatedAt:  bid.CreatedAt,
		Version:    bid.Version,
	}

	return bidDom, nil
}

func (s BidService) Rollback(ctx context.Context, username, bidId string, version int) (*domain.RollbackBidResp, error) {
	authorId, err := s.bidRep.GetAuthorId(ctx, bidId)
	if err != nil {
		return nil, err
	}

	userId, err := s.bidRep.GetUserIdByName(ctx, username)
	if err != nil {
		return nil, err
	}

	if authorId != userId {
		return nil, domain.ErrNotBidAuthor
	}

	err = s.bidRep.Rollback(ctx, bidId, version)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid rollback")
	}

	bid, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	resp := &domain.RollbackBidResp{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		CreatedAt:  bid.CreatedAt,
		Version:    bid.Version,
	}

	return resp, nil
}

func (s BidService) Reviews(ctx context.Context, requesterName, authorName, tenderId string, offset, limit int) ([]domain.ReviewResp, error) {
	if realAuthor, err := s.tenderRep.AuthorByTenderId(ctx, tenderId); err != nil || realAuthor != authorName {
		return nil, domain.ErrAuthorIsIncorrect
	}
	if isResponsible, err := s.tenderRep.UsernameBelongsToTenderOrg(ctx, requesterName, tenderId); err != nil || !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	reviews, err := s.feedbackRep.Reviews(ctx, authorName, offset, limit)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid reviews")
	}

	resp := make([]domain.ReviewResp, len(reviews))
	for i := range reviews {
		resp[i] = domain.ReviewResp{
			Id:          reviews[i].Id,
			Description: reviews[i].Content,
			CreatedAt:   reviews[i].CreatedAt,
		}
	}

	return resp, nil
}
