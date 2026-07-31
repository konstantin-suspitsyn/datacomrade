-- +goose up

ALTER TABLE dc.user_table_roles ADD updated_by_id bigint NOT NULL;
ALTER TABLE dc.user_table_roles ADD CONSTRAINT user_table_roles_creator_user_fk FOREIGN KEY (updated_by_id) REFERENCES dc."user"(id);

-- +goose down

ALTER TABLE dc.user_table_roles DROP CONSTRAINT IF EXISTS user_table_roles_creator_user_fk;
ALTER TABLE dc.user_table_roles DROP COLUMN IF EXISTS updated_by_id;
