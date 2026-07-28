-- +goose up
ALTER TABLE dc.table_cat ADD CONSTRAINT table_cat_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.table_cat ADD CONSTRAINT table_cat_schema_cat_fk FOREIGN KEY (schema_id) REFERENCES dc.schema_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.table_cat ADD CONSTRAINT table_cat_table_type_fk FOREIGN KEY (table_type_id) REFERENCES dc.table_type(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.table_cat ADD CONSTRAINT table_cat_domain_cat_fk FOREIGN KEY (domain_id) REFERENCES dc.domain_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.table_cat DROP CONSTRAINT table_cat_user_fk;
ALTER TABLE dc.table_cat DROP CONSTRAINT table_cat_schema_cat_fk;
ALTER TABLE dc.table_cat DROP CONSTRAINT table_cat_table_type_fk;
ALTER TABLE dc.table_cat DROP CONSTRAINT table_cat_domain_cat_fk;