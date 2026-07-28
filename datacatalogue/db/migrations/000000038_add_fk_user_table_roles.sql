-- +goose up
ALTER TABLE dc.user_table_roles ADD CONSTRAINT user_table_roles_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.user_table_roles ADD CONSTRAINT user_table_roles_table_roles_fk FOREIGN KEY (table_roles_id) REFERENCES dc.table_roles(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.user_table_roles DROP CONSTRAINT user_table_roles_user_fk;
ALTER TABLE dc.user_table_roles DROP CONSTRAINT user_table_roles_table_roles_fk;
