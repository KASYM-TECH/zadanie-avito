package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func ExtractQuery(req *http.Request, key, defaultValue string) (string, bool) {
	if val := req.URL.Query().Get(key); req.URL.Query().Has(key) || val != "" {
		return val, true
	}
	return defaultValue, false
}

func ExtractQueryMany(req *http.Request, key string) ([]string, bool) {
	if vals := req.URL.Query()[key]; len(vals) > 0 {
		return vals, true
	}
	return nil, false
}

func ExtractParam(req *http.Request, key, defaultValue string) (string, bool) {
	if val, ok := mux.Vars(req)[key]; ok && val != "" {
		return val, true
	}
	return defaultValue, false
}
