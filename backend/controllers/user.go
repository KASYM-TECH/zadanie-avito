package controllers

import (
	"avito/domain"
	"avito/log"
	"context"
	"errors"
)

type UserService interface {
	Signup(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, string, error)
}

type UserController struct {
	log         log.Logger
	userService UserService
}

func NewUserController(log log.Logger, userService UserService) *UserController {
	return &UserController{log: log, userService: userService}
}

func (u *UserController) Signup(ctx context.Context, signupRequest domain.SignupRequest, rd domain.RequestData) (string, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_login", signupRequest.Login)
	u.log.Info(ctx, "signup handler")

	userId, err := u.userService.Signup(ctx, signupRequest.Login, signupRequest.Password)
	switch {
	case errors.Is(err, domain.BusinessErr):
		return "", &domain.HTTPError{Cause: err, Detail: "could not signup user", Status: domain.ServerFailureCode}
	case errors.Is(err, domain.DbErr):
		return "", &domain.HTTPError{Cause: err, Detail: "could not signup user", Status: domain.BadRequestCode}

	default:
		return userId, nil
	}
}

func (u *UserController) Login(ctx context.Context, loginRequest domain.LoginRequest, rd domain.RequestData) (*domain.LoginResponse, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_login", loginRequest.Login)
	u.log.Info(ctx, "login handler")

	refreshToken, accessToken, err := u.userService.Login(ctx, loginRequest.Login, loginRequest.Password)
	switch {
	case errors.Is(err, domain.ClientErr):
		return nil, &domain.HTTPError{Cause: err, Detail: "could not login user", Status: domain.BadRequestCode}
	case errors.Is(err, domain.InternalErr):
		return nil, &domain.HTTPError{Cause: err, Detail: "could not login user", Status: domain.ServerFailureCode}

	default:
		return &domain.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
	}
}
