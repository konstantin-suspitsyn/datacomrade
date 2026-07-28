-- +goose up
CREATE TABLE IF NOT EXISTS dc.following_calculation (
	id bigserial NOT NULL,
	column_cat_id bigint NOT NULL,
	calculation_type_id bigint NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	user_id bigint NOT NULL,
	CONSTRAINT following_calculation_pk PRIMARY KEY (id)
);

-- +goose down
DROP TABLE IF EXISTS dc.following_calculation;