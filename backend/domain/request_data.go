package domain

import (
	"avito/auth"
	"net/http"
)

type RequestData struct {
	Claims  auth.Claims
	Request *http.Request
	UserId  string
}
