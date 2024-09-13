package domain

import (
	"net/http"
)

type RequestData struct {
	Request *http.Request
}
