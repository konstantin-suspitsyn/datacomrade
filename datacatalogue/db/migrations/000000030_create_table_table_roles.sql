-- +goose up

CREATE TABLE dc.table_roles (
	id bigserial NOT NULL,
	name varchar(128) NOT NULL,
	description varchar(2000) NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	CONSTRAINT table_roles_pk PRIMARY KEY (id),
	CONSTRAINT table_roles_name_unique UNIQUE (name)
);

-- +goose down

DROP TABLE IF EXISTS dc.table_roles;