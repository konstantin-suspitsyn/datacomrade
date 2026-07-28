-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.table_type (
	id bigserial NOT NULL,
	"name" varchar(128) NOT NULL,
	description varchar(1000) NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT table_type_name_unique UNIQUE (name),
	CONSTRAINT table_type_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.table_type;
