package service

import (
	"avito/auth/jwt"
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/pkg/errors"
)

type UserRep interface {
	Insert(ctx context.Context, employee model.Employee) (string, error)
	GetIdByUsername(ctx context.Context, username string) (string, error)
}

type UserService struct {
	userRep UserRep
}

func NewUserService(userRep UserRep) UserService {
	return UserService{userRep: userRep}
}

func (u UserService) Signup(ctx context.Context, employee model.Employee) (string, error) {
	id, err := u.userRep.Insert(ctx, employee)
	if err != nil {
		return "", errors.WithMessage(err, "Service.Signup insert user")
	}

	return id, nil
}

func (u UserService) Login(ctx context.Context, username string) (string, string, error) {
	userID, err := u.userRep.GetIdByUsername(ctx, username)
	if err != nil {
		return "", "", errors.WithMessage(err, "Service.Login could not fetch user from db")
	}

	refreshToken, genErr := jwt.NewRefreshToken(userID, username)
	if !genErr {
		return "", "", errors.WithMessage(domain.ErrInternal, "could not generate token")
	}

	accessToken, genErr := jwt.NewAccessToken(userID, username)
	if !genErr {
		return "", "", errors.WithMessage(domain.ErrInternal, "could not generate token")
	}

	return refreshToken, accessToken, nil
}
