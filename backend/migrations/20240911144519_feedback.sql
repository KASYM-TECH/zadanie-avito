-- +goose Up
CREATE TABLE IF NOT EXISTS feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL CHECK (char_length(content) <= 1000),
    bid_id UUID NOT NULL REFERENCES bid(id),
    author_id UUID NOT NULL REFERENCES employee(id),
    receiver_id UUID NOT NULL REFERENCES employee(id),
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '3 hours')
);

-- +goose Down
DROP TABLE feedback CASCADE;