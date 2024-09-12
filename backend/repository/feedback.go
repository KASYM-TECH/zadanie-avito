package repository

import (
	"avito/db"
	"avito/db/model"
	"avito/log"
	"context"
	"github.com/pkg/errors"
)

type FeedbackRep struct {
	cli    db.DB
	logger log.Logger
}

func NewFeedbackRep(logger log.Logger, cli db.DB) *FeedbackRep {
	return &FeedbackRep{
		logger: logger, cli: cli,
	}
}

func (rep *FeedbackRep) SaveFeedback(ctx context.Context, feedback *model.Feedback) error {
	_, err := rep.cli.Exec(ctx, `INSERT INTO feedback(bid_id, content, author_id, receiver_id) 
									   VALUES($1, $2, $3, (SELECT author_id FROM bid WHERE id = $1))`,
		feedback.BidID, feedback.Content, feedback.AuthorID)

	if err != nil {
		return errors.WithMessage(err, "Repository.Feedback.SaveFeedback with id: "+feedback.BidID)
	}

	return nil
}

func (rep *FeedbackRep) Reviews(ctx context.Context, authorName string, offset, limit int) ([]model.Feedback, error) {
	query := `WITH author_id_t AS (SELECT id FROM employee WHERE username = $1)
			  SELECT id, content, created_at FROM feedback 
			  WHERE author_id = (SELECT author_id FROM author_id_t) OFFSET $2`

	if offset < 0 {
		offset = 0
	}

	var (
		err     error
		reviews []model.Feedback
	)
	if limit > 0 {
		query += ` LIMIT $3`
		err = rep.cli.Select(ctx, &reviews, query, authorName, offset, limit)
	} else {
		err = rep.cli.Select(ctx, &reviews, query, authorName, offset)
	}

	if err != nil {
		return reviews, errors.WithMessage(err, "Repository.Feedback.Reviews with author username: "+authorName)
	}

	return reviews, nil
}
