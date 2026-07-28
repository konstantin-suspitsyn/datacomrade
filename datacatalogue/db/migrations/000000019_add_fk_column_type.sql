-- +goose up
ALTER TABLE dc.column_type ADD CONSTRAINT column_type_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.column_type DROP CONSTRAINT column_type_user_fk;