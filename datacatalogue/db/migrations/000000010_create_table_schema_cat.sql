-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.schema_cat (
	id bigserial NOT NULL,
	database_id int8 NOT NULL,
	"name" varchar(128) NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT schema_cat_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.schema_cat;
