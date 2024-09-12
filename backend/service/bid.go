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
	GetByUserId(ctx context.Context, offset, limit int, userId string) ([]model.Bid, error)
	GetByTenderId(ctx context.Context, offset, limit int, userId string) ([]model.Bid, error)
	GetBidStatus(ctx context.Context, bidId string) (string, error)
	SetBidStatus(ctx context.Context, bidId, status string) error
	UpdateById(ctx context.Context, bid *model.Bid) error
	SetBidStatusIfOpen(ctx context.Context, bidId string, status string) (string, error)
	Rollback(ctx context.Context, bidId string, version int) error
	GetAuthorId(ctx context.Context, bidId string) (string, error)
	GetOrgIdByBidId(ctx context.Context, bidId string) (string, error)
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
	orgRep      OrgRep
	txMan       TxManager
}

func NewBidService(bidRep BidRep, feedbackRep FeedbackRep, tenderRep TenderRep, orgRep OrgRep, txMan TxManager) BidService {
	return BidService{bidRep: bidRep, feedbackRep: feedbackRep, tenderRep: tenderRep, orgRep: orgRep, txMan: txMan}
}

func (s BidService) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	bidId, err := s.bidRep.Insert(ctx, bid)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid insert")
	}

	bidNew, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid insert")
	}

	return bidNew, nil
}

func (s BidService) GetByUserId(ctx context.Context, offset, limit int, userId string) ([]model.Bid, error) {
	bids, err := s.bidRep.GetByUserId(ctx, offset, limit, userId)
	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (s BidService) GetByTenderId(ctx context.Context, offset, limit int, tenderId string) ([]model.Bid, error) {
	bids, err := s.bidRep.GetByTenderId(ctx, offset, limit, tenderId)
	if err != nil {
		return nil, err
	}

	return bids, nil
}

func (s BidService) GetStatus(ctx context.Context, bidId, username string) (string, error) {
	bid, err := s.bidRep.GetBidStatus(ctx, bidId)
	if err != nil {
		return "", err
	}

	return bid, nil
}

func (s BidService) SetStatus(ctx context.Context, bidId, userId, status string) (*model.Bid, error) {
	authorId, err := s.bidRep.GetAuthorId(ctx, bidId)
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

	return bid, nil
}

func (s BidService) Edit(ctx context.Context, userId string, bid *model.Bid) (*model.Bid, error) {
	authorId, err := s.bidRep.GetAuthorId(ctx, bid.ID)
	if err != nil {
		return nil, err
	}
	if authorId != userId {
		return nil, domain.ErrNotBidAuthor
	}

	err = s.bidRep.UpdateById(ctx, bid)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid edit")
	}

	bidUpdated, err := s.bidRep.GetById(ctx, bid.ID)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	return bidUpdated, nil
}

func (s BidService) SubmitDecision(ctx context.Context, userId, bidId string, decision model.Decision) (*model.Bid, error) {
	orgId, err := s.bidRep.GetOrgIdByBidId(ctx, bidId)
	if err != nil {
		return nil, err
	}

	isResponsible, err := s.orgRep.EmpBelongs(ctx, userId, orgId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	err = s.txMan.DecisionTransaction(ctx, func(ctx context.Context, tx DecisionTransaction) error {
		var (
			tenderStatus = model.TenderStatusClosed
			bidStatus    = model.BidStatusRejected
		)

		if decision == model.Approved {
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

	bidUpdated, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	return bidUpdated, nil
}

func (s BidService) SubmitFeedback(ctx context.Context, feedback *model.Feedback) (*model.Bid, error) {
	orgId, err := s.bidRep.GetOrgIdByBidId(ctx, feedback.BidID)
	if err != nil {
		return nil, err
	}

	isResponsible, err := s.orgRep.EmpBelongs(ctx, feedback.AuthorID, orgId)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, domain.ErrUserNotResponsible
	}

	err = s.feedbackRep.SaveFeedback(ctx, feedback)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid submit feedback")
	}

	bidUpdated, err := s.bidRep.GetById(ctx, feedback.BidID)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	return bidUpdated, nil
}

func (s BidService) Rollback(ctx context.Context, userId, bidId string, version int) (*model.Bid, error) {
	authorId, err := s.bidRep.GetAuthorId(ctx, bidId)
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

	bidUpdated, err := s.bidRep.GetById(ctx, bidId)
	if err != nil {
		return nil, errors.WithMessage(err, "Service.Bid get")
	}

	return bidUpdated, nil
}

func (s BidService) Reviews(ctx context.Context, requesterName, authorName, tenderId string, offset, limit int) ([]model.Feedback, error) {
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

	return reviews, nil
}
