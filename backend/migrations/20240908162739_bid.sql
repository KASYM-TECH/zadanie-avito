-- +goose Up
CREATE TYPE bid_status AS ENUM (
  'Created',
  'Published',
  'Canceled',
  'Approved',
  'Rejected'
);

CREATE TYPE bid_author_type AS ENUM (
  'Organization',
  'User'
);

CREATE TABLE IF NOT EXISTS bid (
    id UUId PRIMARY KEY DEFAULT gen_random_uuid(),
    status bid_status NOT NULL DEFAULT 'Created',
    tender_id UUId NOT NULL REFERENCES tender(id),
    author_type bid_author_type NOT NULL,
    author_id UUId NOT NULL REFERENCES employee(id),
    version INTEGER CHECK (version >= 1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '3 hours')
);

CREATE TABLE IF NOT EXISTS bid_content (
    bid_id UUId REFERENCES bid(id),
    version INTEGER CHECK (version >= 1) NOT NULL DEFAULT 1,
    name TEXT NOT NULL CHECK (char_length(name) <= 100),
    description TEXT NOT NULL CHECK (char_length(name) <= 500),
    PRIMARY KEY(bid_id, version)
);

-- +goose Down
DROP TABLE bid CASCADE;
DROP TYPE bid_author_type CASCADE;
DROP TYPE bid_status CASCADE;