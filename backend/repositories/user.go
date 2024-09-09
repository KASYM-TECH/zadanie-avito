package repositories

import (
	"avito/db"
	"avito/db/model"
	"avito/domain"
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

func (rep UserRep) Insert(ctx context.Context, login, hashedPassword string, roleId int) (string, error) {
	var id string
	err := rep.cli.SelectRow(ctx, &id,
		"insert into users(roleId, login, hashedPassword) values ($1, $2, $3) returning id",
		roleId, login, hashedPassword)

	if err != nil {
		return "", errors.WithMessage(domain.DbErr, "Duplicate in user.Insert with login="+login)
	}

	return id, nil
}

func (rep UserRep) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := rep.cli.SelectRow(ctx, &user, "SELECT * FROM users WHERE login = $1", login)

	if err != nil {
		return nil, errors.WithMessage(err, "user.GetUserByLogin with login="+login)
	}

	return &user, nil
}
