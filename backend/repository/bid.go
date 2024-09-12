package repository

import (
	"avito/db"
	"avito/db/model"
	"avito/domain"
	"avito/log"
	"avito/repository/cache"
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type BidRep struct {
	cli      db.DB
	logger   log.Logger
	idsCache *cache.Set
}

func NewBidRep(logger log.Logger, cli db.DB, idsCache *cache.Set) *BidRep {
	bidRep := &BidRep{
		logger:   logger,
		cli:      cli,
		idsCache: idsCache,
	}
	ctx := context.Background()
	ids, err := bidRep.GetBidIds(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Get bid ids error: %v", err))
	}
	idsCache.WarmUp(ids)

	return bidRep
}

func (rep *BidRep) Insert(ctx context.Context, newBid *model.Bid) (string, error) {
	var bidId string
	err := rep.cli.SelectRow(ctx, &bidId,
		`WITH bid_id_t AS (INSERT INTO bid (tender_id, author_type, author_id) VALUES ($1, $2, $3) RETURNING id)
			   INSERT INTO bid_content(name, description, bid_id) VALUES ($4, $5, (SELECT id FROM bid_id_t)) 
               RETURNING (SELECT id FROM bid_id_t)`,
		newBid.TenderID, newBid.AuthorType, newBid.AuthorID,
		newBid.Name, newBid.Description)

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Bid.Insert with name: "+newBid.Name)
	}

	rep.idsCache.Add(bidId)

	return bidId, nil
}

func (rep *BidRep) GetById(ctx context.Context, bidId string) (*model.Bid, error) {
	if !rep.idsCache.Exists(bidId) {
		return nil, domain.ErrBidDoesNotExist
	}

	var bid model.Bid
	err := rep.cli.SelectRow(ctx, &bid,
		`SELECT b.id, c.name, b.status, b.author_type, b.author_id, b.version, b.created_at FROM bid b
    			JOIN bid_content c ON b.id = c.bid_id AND c.version = b.version WHERE id = $1`, bidId)

	if err != nil {
		return nil, errors.WithMessage(err, "Repository.Bid.GetById with id: "+bidId)
	}

	return &bid, nil
}

func (rep *BidRep) GetByUserId(ctx context.Context, offset, limit int, userId string) ([]model.Bid, error) {
	query := `SELECT b.id,
			c.name,
			b.status,
			b.author_type,
			b.author_id,
			b.version,
			b.created_at
 			FROM bid b JOIN bid_content c ON b.id = c.bid_id and b.version = c.version
            WHERE author_id = $1 ORDER BY name OFFSET $2`

	var (
		bid []model.Bid
		err error
	)

	if offset < 0 {
		offset = 0
	}

	if limit > 0 {
		query += ` LIMIT $3`
		err = rep.cli.Select(ctx, &bid, query, userId, offset, limit)
	} else {
		err = rep.cli.Select(ctx, &bid, query, userId, offset)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "Repository.Bid.Get")
	}

	return bid, nil
}

func (rep *BidRep) GetByTenderId(ctx context.Context, offset, limit int, tenderId string) ([]model.Bid, error) {
	query := `SELECT b.id,
			c.name,
			b.status,
			b.author_type,
			b.author_id,
			b.version,
			b.created_at
 			FROM bid b JOIN bid_content c ON b.id = c.bid_id and b.version = c.version
            WHERE tender_id = $1 AND status != 'Created' ORDER BY name OFFSET $2`

	var (
		bid []model.Bid
		err error
	)

	if offset < 0 {
		offset = 0
	}

	if limit > 0 {
		query += ` LIMIT $3`
		err = rep.cli.Select(ctx, &bid, query, tenderId, offset, limit)
	} else {
		err = rep.cli.Select(ctx, &bid, query, tenderId, offset)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "Repository.Bid.GetByTenderId")
	}

	return bid, nil
}

func (rep *BidRep) GetBidStatus(ctx context.Context, bidId string) (string, error) {
	if !rep.idsCache.Exists(bidId) {
		return "", domain.ErrBidDoesNotExist
	}

	var status string
	err := rep.cli.SelectRow(ctx, &status,
		`SELECT status FROM bid WHERE id = $1`, bidId)

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Bid.GetBidStatus with id: "+bidId)
	}

	return status, nil
}

func (rep *BidRep) SetBidStatusById(ctx context.Context, bidId, status string) error {
	if !rep.idsCache.Exists(bidId) {
		return domain.ErrBidDoesNotExist
	}

	_, err := rep.cli.Exec(ctx, `UPDATE bid SET status = $1 WHERE id = $2`, status, bidId)

	if err != nil {
		return errors.WithMessage(err, "Repository.Bid.SetBidStatusById with id: "+bidId)
	}

	return nil
}

func (rep *BidRep) UpdateById(ctx context.Context, bid *model.Bid) error {
	if !rep.idsCache.Exists(bid.ID) {
		return domain.ErrBidDoesNotExist
	}

	_, err := rep.cli.Exec(ctx,
		`WITH version_t AS (
					INSERT INTO bid_content(name, description, version, bid_id) 
					VALUES($1, $2, (SELECT MAX(version) FROM bid_content WHERE bid_id = $3)+1, $3) RETURNING version)
			UPDATE bid SET version=(SELECT version FROM version_t) WHERE id = $3`,
		bid.Name, bid.Description, bid.ID)

	if err != nil {
		return errors.WithMessage(err, "Repository.Bid.UpdateById with id: "+bid.ID)
	}

	return nil
}

func (rep *BidRep) GetAuthorId(ctx context.Context, bidId string) (string, error) {
	if !rep.idsCache.Exists(bidId) {
		return "", domain.ErrBidDoesNotExist
	}

	var authorId string
	err := rep.cli.SelectRow(ctx, &authorId, `SELECT author_id FROM bid WHERE id = $1`, bidId)

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Bid.UpdateById with id: "+bidId)
	}

	return authorId, nil
}

func (rep *BidRep) SetBidStatus(ctx context.Context, bidId string, status string) error {
	if !rep.idsCache.Exists(bidId) {
		return domain.ErrBidDoesNotExist
	}

	res, err := rep.cli.Exec(ctx,
		`UPDATE bid SET status = $1 WHERE id = $2`, status, bidId)

	if err != nil {
		return errors.WithMessage(err, "Repository.Bid.SetBidStatus with id: "+bidId)
	}

	if num, err := res.RowsAffected(); err != nil || num == 0 {
		return errors.WithMessage(domain.ErrPublishedBidNotFound, "Repository.Bid.SetBidStatusIfOpen with id: "+bidId)
	}

	return nil
}

func (rep *BidRep) SetBidStatusIfOpen(ctx context.Context, bidId string, status string) (string, error) {
	if !rep.idsCache.Exists(bidId) {
		return "", domain.ErrBidDoesNotExist
	}

	var tenderId string
	err := rep.cli.SelectRow(ctx, &tenderId, `UPDATE bid SET status = $1 WHERE id = $2
                             						AND status = 'Published' RETURNING tender_id`, status, bidId)

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Bid.SetBidStatusIfOpen with id: "+bidId)
	}

	return tenderId, nil
}

func (rep *BidRep) Rollback(ctx context.Context, bidId string, version int) error {
	if !rep.idsCache.Exists(bidId) {
		return domain.ErrBidDoesNotExist
	}

	_, err := rep.cli.Exec(ctx,
		`WITH 
					last_version AS (
    					INSERT INTO bid_content(name, description, bid_id, version) 
    					SELECT name, description, bid_id, 
							(SELECT MAX(version) FROM bid_content WHERE bid_id = $1) + 1
						FROM bid_content
						WHERE bid_id = $1 and version = $2
						RETURNING version)
				UPDATE bid SET version=(SELECT version FROM last_version) WHERE id = $1`,
		bidId, version)

	if err != nil {
		return errors.WithMessage(err, "Repository.Bid.Rollback with id: "+bidId)
	}

	return nil
}

func (rep *BidRep) GetOrgIdByBidId(ctx context.Context, bidId string) (string, error) {
	if !rep.idsCache.Exists(bidId) {
		return "", domain.ErrBidDoesNotExist
	}

	var orgId string
	err := rep.cli.SelectRow(ctx, &orgId, `SELECT t.organization_id FROM bid b
            										JOIN tender t ON t.id = b.tender_id WHERE b.id = $1`, bidId)

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Bid.UpdateById with id: "+bidId)
	}

	return orgId, nil
}

func (rep *BidRep) GetBidIds(ctx context.Context) ([]string, error) {
	var ids []string
	err := rep.cli.Select(ctx, &ids, `SELECT id FROM bid`)

	if err != nil {
		return nil, err
	}

	return ids, nil
}
