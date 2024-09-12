//nolint:gochecknoglobals
package domain

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrDb            = errors.New("Database error")
	ErrInternal      = errors.New("Internal error")
	ErrBusiness      = errors.New("Business error")
	ErrClient        = errors.New("Client error")
	ErrInvalidConfig = errors.New("Invalid config")
	ErrJsonParse     = errors.New("Invalid json")

	ErrUserNotResponsible   = errors.New("User not responsible for org")
	ErrPublishedBidNotFound = errors.New("Published bid with this id does not exist")
	ErrAuthorIsIncorrect    = errors.New("Specified author is not author of the tender")
	ErrNotBidAuthor         = errors.New("You must be the author of the bid")
	ErrBidIsNotPublished    = errors.New("Bid is not published")
	ErrTenderIsNotPublished = errors.New("Tender is not published")
	ErrForbiddenApproval    = errors.New("You do not have permission for self approval")
	ErrBidDoesNotExist      = errors.New("Bid with this id does not exist")
	ErrTenderDoesNotExist   = errors.New("Tender with this id does not exist")
)

type StatusCode int

const (
	SuccessCode       StatusCode = 200
	BadRequestCode    StatusCode = 400
	UnauthorizedCode  StatusCode = 401
	ForbiddenCode     StatusCode = 403
	ServerFailureCode StatusCode = 500
)

type HTTPError struct {
	Cause  error      `body:"-"`
	Reason string     `body:"Reason"`
	Status StatusCode `body:"-"`
}

func NewHTTPError(err error, detail string, status StatusCode) HTTPError {
	return HTTPError{err, detail, status}
}

func (e HTTPError) String() string {
	if e.Cause == nil {
		return e.Reason
	}
	return "[Reason] " + e.Reason + "\n" + "[Cause] " + e.Cause.Error()
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
