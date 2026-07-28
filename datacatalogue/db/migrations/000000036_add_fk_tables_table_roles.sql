-- +goose up
ALTER TABLE dc.tables_table_roles ADD CONSTRAINT tables_table_roles_table_cat_fk FOREIGN KEY (table_cat_id) REFERENCES dc.table_cat(id) ON UPDATE CASCADE;
ALTER TABLE dc.tables_table_roles ADD CONSTRAINT tables_table_roles_table_roles_fk FOREIGN KEY (table_roles_id) REFERENCES dc.table_roles(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.tables_table_roles DROP CONSTRAINT tables_table_roles_table_cat_fk;
ALTER TABLE dc.tables_table_roles DROP CONSTRAINT tables_table_roles_table_roles_fk;