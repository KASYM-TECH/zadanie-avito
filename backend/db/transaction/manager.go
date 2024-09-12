package transaction

import (
	"avito/db"
	"avito/log"
	"avito/repository"
	"avito/repository/cache"
	"avito/service"
	"context"
)

type Manager struct {
	db               db.Transactional
	logger           log.Logger
	bidIdsStorage    *cache.Set
	tenderIdsStorage *cache.Set
}

func NewManager(db db.Transactional,
	logger log.Logger,
	bidIdsStorage *cache.Set,
	tenderIdsStorage *cache.Set) *Manager {
	return &Manager{
		db:               db,
		logger:           logger,
		tenderIdsStorage: tenderIdsStorage,
		bidIdsStorage:    bidIdsStorage,
	}
}

type decisionTx struct {
	*repository.BidRep
	*repository.TenderRep
}

func (m Manager) DecisionTransaction(ctx context.Context, pTx func(ctx context.Context, tx service.DecisionTransaction) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		bidRepo := repository.NewBidRep(m.logger, tx, m.bidIdsStorage)
		tenderRepo := repository.NewTenderRep(m.logger, tx, m.tenderIdsStorage)
		return pTx(ctx, &decisionTx{bidRepo, tenderRepo})
	})
}
