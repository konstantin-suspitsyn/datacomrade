-- +goose up

CREATE SCHEMA IF NOT EXISTS queries;

-- +goose down

DROP SCHEMA IF EXISTS queries;
