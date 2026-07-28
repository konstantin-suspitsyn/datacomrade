-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.database_cat (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	host_id int8 NOT NULL,
	database_type_id int8 NOT NULL,
	description varchar(1000) NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT database_cat_name_unique UNIQUE (name),
	CONSTRAINT database_cat_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.database_cat;
