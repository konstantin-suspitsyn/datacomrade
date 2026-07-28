-- +goose up
ALTER TABLE dc.following_calculation ADD CONSTRAINT following_calculation_column_cat_fk FOREIGN KEY (column_cat_id) REFERENCES dc.column_cat(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.following_calculation ADD CONSTRAINT following_calculation_calculation_type_fk FOREIGN KEY (calculation_type_id) REFERENCES dc.calculation_type(id) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE dc.following_calculation ADD CONSTRAINT following_calculation_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.following_calculation DROP CONSTRAINT following_calculation_column_cat_fk;
ALTER TABLE dc.following_calculation DROP CONSTRAINT following_calculation_calculation_type_fk;
ALTER TABLE dc.following_calculation DROP CONSTRAINT following_calculation_user_fk;