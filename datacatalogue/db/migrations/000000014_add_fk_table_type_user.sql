-- +goose up
ALTER TABLE dc.table_type ADD CONSTRAINT table_type_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.table_type DROP CONSTRAINT table_type_user_fk;