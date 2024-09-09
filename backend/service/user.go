package service

import (
	"avito/auth"
	"avito/auth/jwt"
	"avito/db/model"
	"avito/domain"
	"context"
	"github.com/pkg/errors"
)

type UserRep interface {
	Insert(ctx context.Context, login, hashedPassword string, roleId int) (string, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

type UserService struct {
	UserRep UserRep
}

func NewUserService(userRep UserRep) UserService {
	return UserService{UserRep: userRep}
}

func (u UserService) Signup(ctx context.Context, login, password string) (string, error) {
	hashedPwd := auth.GenerateHashedPassword(password)
	if hashedPwd == "" {
		return "", errors.WithMessage(domain.BusinessErr, "hashing password failed")
	}

	id, err := u.UserRep.Insert(ctx, login, hashedPwd, 1)
	if err != nil {
		return "", errors.WithMessage(err, "insert user")
	}

	return id, nil
}

func (u UserService) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := u.UserRep.GetUserByLogin(ctx, username)
	if err != nil {
		return "", "", errors.WithMessage(domain.ClientErr, err.Error())
	}
	if user == nil {
		return "", "", errors.WithMessage(domain.ClientErr, "could not find user")
	}

	if !auth.IsPasswordCorrect(user.HashedPassword, password) {
		return "", "", errors.WithMessage(domain.ClientErr, "the password is not correct")
	}

	refreshToken, genErr := jwt.NewRefreshToken(user.Id, user.RoleID)
	if !genErr {
		return "", "", errors.WithMessage(domain.ClientErr, "could not generate token")
	}

	accessToken, genErr := jwt.NewAccessToken(user.Id, user.RoleID)
	if !genErr {
		return "", "", errors.WithMessage(domain.InternalErr, "could not generate token")
	}

	return refreshToken, accessToken, nil
}
