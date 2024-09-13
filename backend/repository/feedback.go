package repository

import (
	"avito/db"
	"avito/db/model"
	"avito/domain"
	"avito/log"
	"avito/repository/cache"
	"context"
	"github.com/pkg/errors"
)

type FeedbackRep struct {
	cli                  db.DB
	logger               log.Logger
	usernameIdMatchCache *cache.Storage
}

func NewFeedbackRep(logger log.Logger, cli db.DB, usernameIdMatchCache *cache.Storage) *FeedbackRep {
	return &FeedbackRep{
		logger:               logger,
		cli:                  cli,
		usernameIdMatchCache: usernameIdMatchCache,
	}
}

func (rep *FeedbackRep) SaveFeedback(ctx context.Context, feedback *model.Feedback) error {
	_, err := rep.cli.Exec(ctx, `INSERT INTO feedback(bid_id, content, author_id, receiver_id) 
									   VALUES($1, $2, $3, $4)`,
		feedback.BidId, feedback.Content, feedback.AuthorId, feedback.ReceiverId)

	if err != nil {
		return errors.WithMessage(err, "Repository.Feedback.SaveFeedback with id: "+feedback.BidId)
	}

	return nil
}

func (rep *FeedbackRep) Reviews(ctx context.Context, authorName string, offset, limit int) ([]model.Feedback, error) {
	userId, found := rep.usernameIdMatchCache.Get(authorName)
	if !found {
		return nil, domain.ErrUserWithNameNotFound
	}

	query := `SELECT id, content, created_at FROM feedback 
			  WHERE receiver_id = $1 OFFSET $2`

	if offset < 0 {
		offset = 0
	}

	var (
		err     error
		reviews []model.Feedback
	)
	if limit > 0 {
		query += ` LIMIT $3`
		err = rep.cli.Select(ctx, &reviews, query, userId, offset, limit)
	} else {
		err = rep.cli.Select(ctx, &reviews, query, userId, offset)
	}

	if err != nil {
		return reviews, errors.WithMessage(err, "Repository.Feedback.Reviews with author username: "+authorName)
	}

	return reviews, nil
}
