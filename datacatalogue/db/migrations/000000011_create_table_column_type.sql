-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.column_type (
	id bigserial NOT NULL,
	"name" varchar(128) NOT NULL,
	description varchar(1000) NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT column_type_name_unique UNIQUE (name),
	CONSTRAINT column_type_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.column_type;
