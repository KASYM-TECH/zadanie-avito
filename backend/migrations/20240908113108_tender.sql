-- +goose Up
CREATE TYPE tender_service_type AS ENUM (
  'Construction',
  'Delivery',
  'Manufacture'
);

CREATE TYPE tender_status AS ENUM (
  'Created',
  'Published',
  'Closed'
);

CREATE TABLE IF NOT EXISTS tender (
    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL CHECK (char_length(name) <= 100),
    description TEXT NOT NULL CHECK (char_length(name) <= 500),
    service_type tender_service_type NOT NULL,
    status tender_status NOT NULL,
    organization_id BIGINT NOT NULL REFERENCES organization(id),
    version INTEGER CHECK (version >= 1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE tender CASCADE;
DROP TYPE tender_service_type CASCADE;
DROP TYPE tender_status CASCADE;