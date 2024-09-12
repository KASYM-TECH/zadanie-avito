//nolint:lll
package controllers

import (
	"avito/db/model"
	"avito/domain"
	"avito/log"
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserService interface {
	Signup(ctx context.Context, employee model.Employee) (string, error)
	Login(ctx context.Context, username string) (string, string, error)
}

type UserController struct {
	log         log.Logger
	userService UserService
}

func NewUserController(log log.Logger, userService UserService) *UserController {
	return &UserController{log: log, userService: userService}
}

func (u *UserController) Signup(ctx context.Context, signupRequest domain.SignupRequest, rd domain.RequestData) (string, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "username", signupRequest.Username)
	u.log.Info(ctx, "signup handler")

	userId, err := u.userService.Signup(ctx, model.Employee{
		Username:  signupRequest.Username,
		LastName:  signupRequest.LastName,
		FirstName: signupRequest.FirstName})

	if err == nil {
		return userId, nil
	}

	pgErr := &pgconn.PgError{}
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return "", &domain.HTTPError{Cause: err, Reason: "username must be unique", Status: domain.BadRequestCode}
		}
	}

	return "", &domain.HTTPError{Cause: err, Reason: "could not signup user", Status: domain.ServerFailureCode}
}

func (u *UserController) Login(ctx context.Context, loginRequest domain.LoginRequest, rd domain.RequestData) (*domain.LoginResponse, *domain.HTTPError) {
	ctx = log.AddKeyVal(ctx, "user_login", loginRequest.Username)
	u.log.Info(ctx, "login handler")

	refreshToken, accessToken, err := u.userService.Login(ctx, loginRequest.Username)
	switch {
	case errors.Is(err, domain.ErrClient):
		return nil, &domain.HTTPError{Cause: err, Reason: "could not login user", Status: domain.BadRequestCode}
	case errors.Is(err, domain.ErrInternal):
		return nil, &domain.HTTPError{Cause: err, Reason: "could not login user", Status: domain.ServerFailureCode}

	default:
		return &domain.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
	}
}
