package model

import (
	"time"
)

type Feedback struct {
	ID         string    `db:"id"`
	Content    string    `db:"content"`
	BidID      string    `db:"bid_id"`
	AuthorID   string    `db:"author_id"`
	ReceiverID string    `db:"receiver_id"`
	CreatedAt  time.Time `db:"created_at"`
}
