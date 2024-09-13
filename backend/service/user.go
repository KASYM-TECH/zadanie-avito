package service

import (
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

func (u UserService) Signup(ctx context.Context, signupRequest *domain.SignupRequest) (string, error) {
	employee := model.Employee{
		Username:  signupRequest.Username,
		LastName:  signupRequest.LastName,
		FirstName: signupRequest.FirstName}

	id, err := u.userRep.Insert(ctx, employee)
	if err != nil {
		return "", errors.WithMessage(err, "Service.Signup insert user")
	}

	return id, nil
}
