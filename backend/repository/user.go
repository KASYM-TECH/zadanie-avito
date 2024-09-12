package repository

import (
	"avito/db"
	"avito/db/model"
	"avito/log"
	"context"
	"github.com/pkg/errors"
)

type UserRep struct {
	cli    db.DB
	logger log.Logger
}

func NewUserRep(logger log.Logger, cli db.DB) *UserRep {
	return &UserRep{
		logger: logger, cli: cli,
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

	return id, nil
}

func (rep UserRep) GetIdByUsername(ctx context.Context, username string) (string, error) {
	var userID string
	err := rep.cli.SelectRow(ctx, &userID, "SELECT id FROM employee WHERE username = $1", username)

	if err != nil {
		return "", errors.WithMessage(err, "user.GetUserByLogin with login: "+username)
	}

	return userID, nil
}
