//nolint:gochecknoglobals
package model

import (
	"time"
)

type BidStatus string

const (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusApproved  BidStatus = "Approved"
	BidStatusRejected  BidStatus = "Rejected"
)

type BidAuthorType string

const (
	BidAuthorTypeOrganization BidAuthorType = "Organization"
	BidAuthorTypeUser         BidAuthorType = "User"
)

type Decision string

var (
	Approved = Decision("Approved")
	Rejected = Decision("Rejected")
)

type Bid struct {
	Id          string        `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	Status      BidStatus     `db:"status"`
	TenderId    string        `db:"tender_id"`
	AuthorType  BidAuthorType `db:"author_type"`
	AuthorId    string        `db:"author_id"`
	Version     int           `db:"version"`
	CreatedAt   time.Time     `db:"created_at"`
}
