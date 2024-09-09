package domain

import (
	"avito/auth"
	http "github.com/julienschmidt/httprouter"
)

type RequestData struct {
	Params http.Params
	Claims auth.Claims
	RoleId int
}
