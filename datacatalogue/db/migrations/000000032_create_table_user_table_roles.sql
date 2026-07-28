-- +goose up
CREATE TABLE dc.user_table_roles (
	id bigserial NOT NULL,
	user_id bigint NOT NULL,
	table_roles_id bigint NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	CONSTRAINT user_table_roles_pk PRIMARY KEY (id)
);

-- +goose down
DROP TABLE IF EXISTS dc.user_table_roles; 