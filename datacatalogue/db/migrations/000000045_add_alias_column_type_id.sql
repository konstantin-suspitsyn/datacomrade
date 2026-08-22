-- +goose up

ALTER TABLE dc.alias ADD column_type_id bigint DEFAULT -1 NOT NULL;
ALTER TABLE dc.alias ADD CONSTRAINT alias_column_type_id_fk FOREIGN KEY (column_type_id) REFERENCES dc.column_type(id);

-- +goose down

ALTER TABLE dc.alias DROP CONSTRAINT IF EXISTS alias_column_type_id_fk;
ALTER TABLE dc.alias DROP COLUMN IF EXISTS column_type_id;
