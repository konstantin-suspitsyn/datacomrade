-- +goose up
ALTER TABLE dc.database_cat ADD CONSTRAINT database_cat_host_fk FOREIGN KEY (host_id) REFERENCES dc.host(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.database_cat ADD CONSTRAINT database_cat_database_type_fk FOREIGN KEY (database_type_id) REFERENCES dc.database_type(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.database_cat ADD CONSTRAINT database_cat_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;


-- +goose down
ALTER TABLE dc.database_cat DROP CONSTRAINT database_cat_host_fk;
ALTER TABLE dc.database_cat DROP CONSTRAINT database_cat_database_type_fk;
ALTER TABLE dc.database_cat DROP CONSTRAINT database_cat_user_fk;