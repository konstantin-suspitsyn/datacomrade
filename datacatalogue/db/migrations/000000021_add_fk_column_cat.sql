-- +goose up
ALTER TABLE dc.column_cat ADD CONSTRAINT column_cat_table_cat_fk FOREIGN KEY (table_id) REFERENCES dc.table_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.column_cat ADD CONSTRAINT column_cat_alias_fk FOREIGN KEY (alias_id) REFERENCES dc.alias(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.column_cat ADD CONSTRAINT column_cat_column_type_fk FOREIGN KEY (column_type_id) REFERENCES dc.column_type(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.column_cat ADD CONSTRAINT column_cat_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.column_cat ADD CONSTRAINT column_cat_pk PRIMARY KEY (id);


-- +goose down
ALTER TABLE dc.column_cat DROP CONSTRAINT column_cat_table_cat_fk;
ALTER TABLE dc.column_cat DROP CONSTRAINT column_cat_alias_fk;
ALTER TABLE dc.column_cat DROP CONSTRAINT column_cat_column_type_fk;
ALTER TABLE dc.column_cat DROP CONSTRAINT column_cat_user_fk;
ALTER TABLE dc.column_cat DROP CONSTRAINT column_cat_pk;