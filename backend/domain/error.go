package domain

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
)

var (
	DbErr       = errors.New("Database error")
	InternalErr = errors.New("Internal error")
	BusinessErr = errors.New("Business error")
	ClientErr   = errors.New("Client error")
)

type StatusCode int

const (
	SuccessCode       StatusCode = 200
	BadRequestCode    StatusCode = 400
	UnauthorizedCode  StatusCode = 501
	ForbiddenCode     StatusCode = 503
	ServerFailureCode StatusCode = 500
)

type HTTPError struct {
	Cause  error      `body:"-"`
	Detail string     `body:"detail"`
	Status StatusCode `body:"-"`
}

func NewHTTPError(err error, detail string, status StatusCode) HTTPError {
	return HTTPError{err, detail, status}
}

func (e HTTPError) String() string {
	if e.Cause == nil {
		return e.Detail
	}
	return "[Detail] " + e.Detail + "\n" + "[Cause] " + e.Cause.Error()
}

type ErrInContext struct{}

var (
	ErrInContextKey = ErrInContext{}
)

func SetError(r *http.Request, err HTTPError) {
	*r = *r.WithContext(context.WithValue(r.Context(), ErrInContextKey, err))
}

func GetError(r *http.Request) *HTTPError {
	if val, ok := r.Context().Value(ErrInContextKey).(HTTPError); ok {
		return &val
	}
	return nil
}
