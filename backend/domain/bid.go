package domain

import (
	"avito/db/model"
	"time"
)

type CreateBidReq struct {
	Name        string              `validate:"required,lte=100" json:"name"`
	Description string              `validate:"required,lte=500" json:"description"`
	TenderID    string              `validate:"required" json:"tenderId"`
	AuthorType  model.BidAuthorType `validate:"required" json:"authorType"`
	AuthorID    string              `validate:"required" json:"authorId"`
}

type CreateBidResp struct {
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	AuthorType  model.BidAuthorType `json:"authorType"`
	AuthorID    string              `json:"authorId"`
	Version     int                 `json:"version"`
	CreatedAt   time.Time           `json:"createdAt"`
}

type GetBidResp struct {
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Status     model.BidStatus     `json:"status"`
	AuthorType model.BidAuthorType `json:"authorType"`
	AuthorID   string              `json:"authorId"`
	Version    int                 `json:"version"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type SetStatusBidResp struct {
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Status     model.BidStatus     `json:"status"`
	AuthorType model.BidAuthorType `json:"authorType"`
	AuthorID   string              `json:"authorId"`
	Version    int                 `json:"version"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type EditBidReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EditBidResp struct {
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Status     model.BidStatus     `json:"status"`
	AuthorType model.BidAuthorType `json:"authorType"`
	AuthorID   string              `json:"authorId"`
	Version    int                 `json:"version"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type SubmitDecisionBidResp struct {
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Status     model.BidStatus     `json:"status"`
	AuthorType model.BidAuthorType `json:"authorType"`
	AuthorID   string              `json:"authorId"`
	Version    int                 `json:"version"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type FeedbackBidResp struct {
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Status     model.BidStatus     `json:"status"`
	AuthorType model.BidAuthorType `json:"authorType"`
	AuthorID   string              `json:"authorId"`
	Version    int                 `json:"version"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type RollbackBidResp struct {
	Id         string              `json:"id"`
	Name       string              `json:"name"`
	Status     model.BidStatus     `json:"status"`
	AuthorType model.BidAuthorType `json:"authorType"`
	AuthorID   string              `json:"authorId"`
	Version    int                 `json:"version"`
	CreatedAt  time.Time           `json:"createdAt"`
}
