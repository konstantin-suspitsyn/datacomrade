-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.database_type (
	id bigserial NOT NULL,
	"name" varchar(128) NOT NULL,
	db_version varchar(512) NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT database_type_name_unique UNIQUE (name),
	CONSTRAINT database_type_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.database_type;
