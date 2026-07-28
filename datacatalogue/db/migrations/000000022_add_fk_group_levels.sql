-- +goose up
ALTER TABLE dc.group_levels ADD CONSTRAINT group_levels_column_cat_fk FOREIGN KEY (column_id) REFERENCES dc.column_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.group_levels ADD CONSTRAINT group_levels_column_cat_parent_fk FOREIGN KEY (parent_column_id) REFERENCES dc.column_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.group_levels ADD CONSTRAINT group_levels_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.group_levels DROP CONSTRAINT group_levels_column_cat_fk;
ALTER TABLE dc.group_levels DROP CONSTRAINT group_levels_column_cat_parent_fk;
ALTER TABLE dc.group_levels DROP CONSTRAINT group_levels_user_fk;