-- +goose up
ALTER TABLE dc.schema_cat ADD CONSTRAINT schema_cat_database_cat_fk FOREIGN KEY (database_id) REFERENCES dc.database_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.schema_cat ADD CONSTRAINT schema_cat_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.schema_cat DROP CONSTRAINT schema_cat_database_cat_fk;
ALTER TABLE dc.schema_cat DROP CONSTRAINT schema_cat_user_fk;