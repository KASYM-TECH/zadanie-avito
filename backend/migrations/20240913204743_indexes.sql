-- +goose Up
CREATE INDEX IF NOT EXISTS tender_version_btree_idx ON tender (version);
CREATE INDEX IF NOT EXISTS bid_version_btree_idx ON bid (version);

-- +goose Down
DROP INDEX tender_version_btree_idx;
DROP INDEX bid_version_btree_idx;
