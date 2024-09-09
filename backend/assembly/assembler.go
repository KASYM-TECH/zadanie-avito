package assembly

import (
	"avito/config"
	"avito/controllers"
	"avito/db"
	"avito/log"
	"avito/migrations"
	"avito/repositories"
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
	a.closers = append(a.closers, cli.Close)

	mgRunner := migrations.NewRunner(migrations.DialectPostgreSQL, "migrations", a.logger)
	err = mgRunner.Run(ctx, cli.DB.DB)
	if err != nil {
		a.logger.Fatal(ctx, err.Error())
		return nil, err
	}

	userRep := repositories.NewUserRep(a.logger, cli)
	userService := service.NewUserService(userRep)
	userController := controllers.NewUserController(a.logger, userService)

	_ = service.NewBannerService()
	bannerController := controllers.NewBannerController(a.logger)

	r := server.NewRouter(a.logger)
	middlewares := server.NewMiddleware(a.logger)
	r.AddRoutes(middlewares, userController, bannerController)

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
