-- +goose Up
CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE IF NOT EXISTS organization (
    id UUId PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '3 hours'),
    updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '3 hours')
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUId PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUId REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUId REFERENCES employee(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE organization_responsible CASCADE;
DROP TABLE organization CASCADE;
DROP TYPE organization_type CASCADE;
