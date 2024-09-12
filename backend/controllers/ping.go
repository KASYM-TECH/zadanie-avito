package controllers

import (
	"avito/domain"
	"avito/log"
	"context"
)

type DummyController struct {
	log log.Logger
}

func NewDummyController(log log.Logger) *DummyController {
	return &DummyController{log: log}
}

func (u *DummyController) Ping(ctx context.Context, rd domain.RequestData) (string, *domain.HTTPError) {
	u.log.Info(ctx, "ping-pong")

	return "pong", nil
}
