
-- +migrate Up
ALTER TABLE products ALTER COLUMN "weight" type NUMERIC(10,3);
-- +migrate Down
ALTER TABLE products ALTER COLUMN "weight" type BIGINT;