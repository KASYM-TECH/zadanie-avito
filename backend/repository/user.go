package repository

import (
	"avito/db"
	"avito/db/model"
	"avito/log"
	"avito/repository/cache"
	"context"
	"github.com/pkg/errors"
)

type UserRep struct {
	cli                  db.DB
	logger               log.Logger
	usernameIdMatchCache *cache.Storage
}

func NewUserRep(logger log.Logger, cli db.DB, usernameIdMatchCache *cache.Storage) *UserRep {
	return &UserRep{
		logger:               logger,
		cli:                  cli,
		usernameIdMatchCache: usernameIdMatchCache,
	}
}

func (rep UserRep) Insert(ctx context.Context, employee model.Employee) (string, error) {
	var id string
	err := rep.cli.SelectRow(ctx, &id,
		"insert into employee(username, first_name, last_name) values ($1, $2, $3) returning id",
		employee.Username, employee.FirstName, employee.LastName)

	if err != nil {
		return id, errors.WithMessage(err, "Repository.User.Insert with username: "+employee.Username)
	}

	rep.usernameIdMatchCache.Add(employee.Username, id)

	return id, nil
}

func (rep UserRep) GetIdByUsername(ctx context.Context, username string) (string, error) {
	var userId string
	err := rep.cli.SelectRow(ctx, &userId, "SELECT id FROM employee WHERE username = $1", username)

	if err != nil {
		return "", errors.WithMessage(err, "user.GetUserByLogin with login: "+username)
	}

	return userId, nil
}

func (rep UserRep) LoadUsernameId(ctx context.Context) ([]cache.KeyValue, error) {
	usernameIdMatch := make([]cache.KeyValue, 0)

	err := rep.cli.Select(ctx, &usernameIdMatch, "SELECT username as Key, id as Value FROM employee")

	if err != nil {
		return nil, errors.WithMessage(err, "user.LoadUsernameId")
	}

	return usernameIdMatch, nil
}
