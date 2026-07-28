-- +goose up
ALTER TABLE dc.column_cat ADD CONSTRAINT column_cat_calculation_type_fk FOREIGN KEY (calculation_type_id) REFERENCES dc.calculation_type(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.column_cat DROP CONSTRAINT column_cat_calculation_type_fk;