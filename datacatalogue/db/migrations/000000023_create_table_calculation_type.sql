-- +goose up
CREATE TABLE IF NOT EXISTS dc.calculation_type (
	id bigserial NOT NULL,
	"name" varchar(52) NOT NULL,
	description varchar(1000) NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	CONSTRAINT calculation_type_pk PRIMARY KEY (id),
	CONSTRAINT calculation_type_name_unique UNIQUE ("name")
);

-- +goose down
DROP TABLE IF EXISTS dc.calculation_type;