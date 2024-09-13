//nolint:lll
package controllers

import (
	"avito/domain"
	"avito/log"
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserService interface {
	Signup(ctx context.Context, signupRequest *domain.SignupRequest) (string, error)
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

	userId, err := u.userService.Signup(ctx, &signupRequest)

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
