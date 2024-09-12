package assembly

import (
	"avito/config"
	"avito/controllers"
	"avito/db"
	"avito/db/transaction"
	"avito/log"
	"avito/migrations"
	"avito/repository"
	"avito/repository/cache"
	"avito/server"
	"avito/service"
	"context"
	"net/http"
)

type Close func() error

type Assembler struct {
	logger  log.Logger
	closers []Close
}

func NewAssembler(logger log.Logger) *Assembler {
	return &Assembler{
		logger: logger,
	}
}

func (a *Assembler) Assemble(ctx context.Context, conf *config.Config) (http.Handler, error) {
	cli, err := db.Open(ctx, conf.Dsn())
	if err != nil {
		a.logger.Fatal(ctx, err.Error())
		return nil, err
	}

	err = cli.CreateSchema(conf.DbSchema)
	if err != nil {
		a.logger.Fatal(ctx, err.Error())
		return nil, err
	}

	err = cli.SwitchSchema(conf.DbSchema)
	if err != nil {
		a.logger.Fatal(ctx, err.Error())
		return nil, err
	}

	a.closers = append(a.closers, cli.Close)

	mgRunner := migrations.NewRunner(migrations.DialectPostgreSQL, conf.MigrationDir, a.logger)
	err = mgRunner.Run(ctx, cli.DB.DB)
	if err != nil {
		a.logger.Fatal(ctx, err.Error())
		return nil, err
	}

	var (
		bidIdStorage    = cache.NewSet()
		tenderIdStorage = cache.NewSet()
		txManager       = transaction.NewManager(cli, a.logger, bidIdStorage, tenderIdStorage)
	)
	dummyController := controllers.NewDummyController(a.logger)

	orgRep := repository.NewOrganizationRep(a.logger, cli)
	orgService := service.NewOrganizationService(orgRep)
	orgController := controllers.NewOrganizationController(a.logger, orgService)

	userRep := repository.NewUserRep(a.logger, cli)
	userService := service.NewUserService(userRep)
	userController := controllers.NewUserController(a.logger, userService)

	tenderRep := repository.NewTenderRep(a.logger, cli, tenderIdStorage)
	tenderService := service.NewTenderService(tenderRep, orgRep)
	tenderController := controllers.NewTenderController(a.logger, tenderService)

	bidRep := repository.NewBidRep(a.logger, cli, bidIdStorage)
	feedbackRep := repository.NewFeedbackRep(a.logger, cli)
	bidService := service.NewBidService(bidRep, feedbackRep, tenderRep, orgRep, txManager)
	bidController := controllers.NewBidController(a.logger, bidService)

	r := server.NewRouter(a.logger)
	middlewares := server.NewMiddleware(a.logger)
	r.AddRoutes(middlewares, server.Controllers{
		DummyCnt:  dummyController,
		UserCnt:   userController,
		TenderCnt: tenderController,
		OrgCnt:    orgController,
		BidCnt:    bidController})

	return r.Router, nil
}

func (a *Assembler) Close(ctx context.Context) error {
	for _, closeFunc := range a.closers {
		err := closeFunc()
		if err != nil {
			return err
		}
	}
	return nil
}
