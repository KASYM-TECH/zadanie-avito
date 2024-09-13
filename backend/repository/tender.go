package repository

import (
	"avito/db"
	"avito/db/model"
	"avito/domain"
	"avito/log"
	"avito/repository/cache"
	"context"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"strings"
)

type TenderRep struct {
	cli                  db.DB
	logger               log.Logger
	idsCache             *cache.Set
	usernameIdMatchCache *cache.Storage
}

func NewTenderRep(logger log.Logger, cli db.DB, idsCache *cache.Set, usernameIdMatchCache *cache.Storage) *TenderRep {
	tenderRep := &TenderRep{
		logger:               logger,
		cli:                  cli,
		idsCache:             idsCache,
		usernameIdMatchCache: usernameIdMatchCache,
	}
	ctx := context.Background()
	ids, err := tenderRep.GetTenderIds(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Get tenders ids error: %v", err))
	}
	idsCache.WarmUp(ids)

	return tenderRep
}

func (rep *TenderRep) Insert(ctx context.Context, newTender *model.Tender, username string) (string, error) {
	userId, found := rep.usernameIdMatchCache.Get(username)
	if !found {
		return "", domain.ErrUserWithNameNotFound
	}

	var tenderId string
	err := rep.cli.SelectRow(ctx, &tenderId,
		`WITH tender_id_t AS (INSERT INTO tender (status, organization_id, user_id) VALUES ($1, $2, $3) RETURNING id)
			   INSERT INTO tender_content(name, description, service_type, tender_id) VALUES ($4, $5, $6, 
				(SELECT id FROM tender_id_t)) RETURNING (SELECT id FROM tender_id_t)`,
		newTender.Status, newTender.OrganizationId, userId,
		newTender.Name, newTender.Description, newTender.ServiceType)

	pgErr := &pgconn.PgError{}
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return "", domain.ErrOrganizationDoesNotExist
		}
	}

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Tender.Insert with name: "+newTender.Name)
	}

	rep.idsCache.Add(tenderId)

	return tenderId, nil
}

func (rep *TenderRep) GetPublished(ctx context.Context, offset, limit int, types []string) ([]model.Tender, error) {
	query := `SELECT t.id, c.name, c.description, c.service_type, t.status, t.version, t.created_at FROM tender t
 				JOIN tender_content c ON t.id = c.tender_id and t.version = c.version 
 				WHERE status = 'Published'`

	if len(types) != 0 {
		query += ` WHERE c.service_type IN ('` + strings.Join(types, "', '") + `')`
	}

	var (
		tenders []model.Tender
		err     error
	)

	if offset < 0 {
		offset = 0
	}

	query += ` ORDER BY c.name OFFSET $1`
	if limit > 0 {
		query += ` LIMIT $2`
		err = rep.cli.Select(ctx, &tenders, query, offset, limit)
	} else {
		err = rep.cli.Select(ctx, &tenders, query, offset)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "Repository.Tender.Get")
	}

	return tenders, nil
}

func (rep *TenderRep) GetByUsername(ctx context.Context, offset, limit int, username string) ([]model.Tender, error) {
	userId, found := rep.usernameIdMatchCache.Get(username)
	if !found {
		return nil, domain.ErrUserWithNameNotFound
	}
	query := `
			SELECT t.id,
				c.name,
				c.description,
				c.service_type,
				t.status,
				t.version,
				t.created_at
 			FROM tender t JOIN tender_content c ON t.id = c.tender_id and t.version = c.version
            WHERE user_id = $1 ORDER BY name OFFSET $2`

	var (
		tenders []model.Tender
		err     error
	)

	if offset < 0 {
		offset = 0
	}

	if limit > 0 {
		query += ` LIMIT $3`
		err = rep.cli.Select(ctx, &tenders, query, userId, offset, limit)
	} else {
		err = rep.cli.Select(ctx, &tenders, query, userId, offset)
	}

	if err != nil {
		return nil, errors.WithMessage(err, "Repository.Tender.GetByUsername")
	}

	return tenders, nil
}

func (rep *TenderRep) GetTenderStatus(ctx context.Context, tenderId string) (string, error) {
	if !rep.idsCache.Exists(tenderId) {
		return "", domain.ErrTenderDoesNotExist
	}

	var tenderStatus string
	err := rep.cli.SelectRow(ctx, &tenderStatus,
		`SELECT status FROM tender WHERE id = $1`, tenderId)

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Tender.GetBidStatus with id: "+tenderId)
	}

	return tenderStatus, nil
}

func (rep *TenderRep) SetTenderStatus(ctx context.Context, tenderId, status string) error {
	_, err := rep.cli.Exec(ctx, `UPDATE tender SET status = $1 WHERE id = $2`, status, tenderId)

	if err != nil {
		return errors.WithMessage(err, "Repository.Tender.GetTenderStatus with id: "+tenderId)
	}

	return nil
}

func (rep *TenderRep) SetTenderStatusIfOpen(ctx context.Context, tenderId, status string) error {
	if !rep.idsCache.Exists(tenderId) {
		return domain.ErrTenderDoesNotExist
	}

	_, err := rep.cli.Exec(ctx, `UPDATE tender SET status = $1 WHERE id = $2 AND status = 'Published'`, status, tenderId)

	if err != nil {
		return errors.WithMessage(err, "Repository.Tender.GetTenderStatus with id: "+tenderId)
	}

	return nil
}

func (rep *TenderRep) UpdateById(ctx context.Context, tender *model.Tender) error {
	if !rep.idsCache.Exists(tender.Id) {
		return domain.ErrTenderDoesNotExist
	}

	_, err := rep.cli.Exec(ctx,
		`WITH version_table AS (
    		INSERT INTO tender_content(name, description, service_type, version, tender_id) 
    		VALUES($1, $2, $3, (SELECT MAX(version) FROM tender_content WHERE tender_id = $4)+1, $4) RETURNING version) 
			UPDATE tender SET version=(SELECT version from version_table) WHERE id = $4`,
		tender.Name, tender.Description, tender.ServiceType, tender.Id)

	if err != nil {
		return errors.WithMessage(err, "Repository.Tender.UpdateById with id: "+tender.Id)
	}

	return nil
}

func (rep *TenderRep) GetById(ctx context.Context, tenderId string) (*model.Tender, error) {
	if !rep.idsCache.Exists(tenderId) {
		return nil, domain.ErrTenderDoesNotExist
	}

	var tender model.Tender
	err := rep.cli.SelectRow(ctx, &tender,
		`SELECT t.id, c.name, c.description, t.status, c.service_type, t.version, t.created_at FROM tender t 
    			JOIN tender_content c ON t.id = c.tender_id AND t.version = c.version WHERE id = $1`, tenderId)

	if err != nil {
		return nil, errors.WithMessage(err, "Repository.Tender.GetById with id: "+tenderId)
	}

	return &tender, nil
}

func (rep *TenderRep) Rollback(ctx context.Context, tenderId string, version int) error {
	_, err := rep.cli.Exec(ctx,
		`WITH 
					last_version AS (
    					INSERT INTO tender_content(name, description, service_type, tender_id, version) 
    					SELECT name, description, service_type, tender_id, 
							(SELECT MAX(version) FROM tender_content WHERE tender_id = $1) + 1
						FROM tender_content
						WHERE tender_id = $1 and version = $2
						RETURNING version)
				UPDATE tender SET version=(SELECT version FROM last_version) WHERE id = $1`,
		tenderId, version)

	if err != nil {
		return errors.WithMessage(err, "Repository.Tender.Rollback with id: "+tenderId)
	}

	return nil
}

func (rep *TenderRep) AuthorByTenderId(ctx context.Context, tenderId string) (string, error) {
	if !rep.idsCache.Exists(tenderId) {
		return "", domain.ErrTenderDoesNotExist
	}

	var authorUsername string
	err := rep.cli.SelectRow(ctx, &authorUsername,
		`SELECT e.username FROM tender t JOIN employee e ON e.id = t.user_id WHERE t.id = $1`, tenderId)

	if err != nil {
		return "", errors.WithMessage(err, "Repository.Tender.AuthorByTenderId with id: "+tenderId)
	}

	return authorUsername, nil
}

func (rep *TenderRep) UsernameBelongsToTenderOrg(ctx context.Context, username, tenderId string) (bool, error) {
	userId, found := rep.usernameIdMatchCache.Get(username)
	if !found {
		return false, domain.ErrUserWithNameNotFound
	}
	var belongs bool
	err := rep.cli.SelectRow(ctx, &belongs,
		`SELECT EXISTS (
    WITH organization_id_t AS (
    	SELECT organization_id FROM tender WHERE id = $1
    )
    SELECT 1 FROM organization_responsible 
    WHERE organization_id = (SELECT organization_id FROM organization_id_t)
    AND user_id = $2);`, tenderId, userId)

	if err != nil {
		return false, errors.WithMessage(err, "Repository.Tender.UsernameBelongsToTenderOrg with id: "+tenderId)
	}

	return belongs, nil
}

func (rep *TenderRep) GetTenderIds(ctx context.Context) ([]string, error) {
	var ids []string
	err := rep.cli.Select(ctx, &ids, `SELECT id FROM tender`)

	if err != nil {
		return nil, err
	}

	return ids, nil
}
