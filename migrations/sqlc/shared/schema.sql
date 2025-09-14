CREATE SCHEMA IF NOT EXISTS shared;

CREATE TABLE shared."domain" (
	id bigserial NOT NULL,
	"name" varchar NOT NULL,
	description varchar NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	user_id int8 DEFAULT '-1'::integer NOT NULL,
	CONSTRAINT domain_name_unique UNIQUE (name),
	CONSTRAINT domain_pk PRIMARY KEY (id)
);
