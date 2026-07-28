-- +goose up
ALTER TABLE dc.alias ADD CONSTRAINT alias_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.alias DROP CONSTRAINT alias_user_fk;