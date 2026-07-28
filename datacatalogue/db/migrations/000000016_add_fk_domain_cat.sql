-- +goose up
ALTER TABLE dc.domain_cat ADD CONSTRAINT domain_cat_user_fk FOREIGN KEY (user_id) REFERENCES dc."user"(id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- +goose down
ALTER TABLE dc.domain_cat DROP CONSTRAINT domain_cat_user_fk;