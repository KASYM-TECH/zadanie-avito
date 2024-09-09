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

CREATE TABLE bid (
    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (char_length(name) <= 100),
    description TEXT NOT NULL CHECK (char_length(name) <= 500),
    status bid_status NOT NULL,
    tender_id TEXT NOT NULL REFERENCES tender(id),
    author_type bid_author_type NOT NULL,
    author_id BIGINT NOT NULL REFERENCES employee(id),
    version INTEGER CHECK (version >= 1),
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE bid CASCADE;
DROP TYPE bid_author_type CASCADE;
DROP TYPE bid_status CASCADE;