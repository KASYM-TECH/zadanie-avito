-- +goose Up
CREATE TABLE IF NOT EXISTS feedback (
    id UUId PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL CHECK (char_length(content) <= 1000),
    bid_id UUId NOT NULL REFERENCES bid(id),
    author_id UUId NOT NULL REFERENCES employee(id),
    receiver_id UUId NOT NULL REFERENCES employee(id),
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '3 hours')
);

-- +goose Down
DROP TABLE feedback CASCADE;