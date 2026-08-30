-- +goose up

CREATE SCHEMA IF NOT EXISTS qe;

-- +goose down

DROP SCHEMA qe;
