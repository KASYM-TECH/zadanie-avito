package model

import (
	"time"
)

type Feedback struct {
	Id         string    `db:"id"`
	Content    string    `db:"content"`
	BidId      string    `db:"bid_id"`
	AuthorId   string    `db:"author_id"`
	ReceiverId string    `db:"receiver_id"`
	CreatedAt  time.Time `db:"created_at"`
}
