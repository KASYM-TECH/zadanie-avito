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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status tender_status NOT NULL,
    organization_id UUID NOT NULL REFERENCES organization(id),
    version INTEGER CHECK (version >= 1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '3 hours'),
    user_id UUID NOT NULL REFERENCES employee(id)
);

CREATE TABLE IF NOT EXISTS tender_content (
    tender_id UUID REFERENCES tender(id),
    version INTEGER CHECK (version >= 1) NOT NULL DEFAULT 1,
    name TEXT NOT NULL CHECK (char_length(name) <= 100),
    description TEXT NOT NULL CHECK (char_length(name) <= 500),
    service_type tender_service_type NOT NULL,
    PRIMARY KEY(tender_id, version)
);

-- +goose Down
DROP TABLE tender CASCADE;
DROP TYPE tender_service_type CASCADE;
DROP TYPE tender_status CASCADE;