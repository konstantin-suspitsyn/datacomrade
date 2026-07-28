-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.has_to_group (
	id bigserial NOT NULL,
	column_id_a int8 NOT NULL,
	column_id_b int8 NOT NULL,
	description varchar(1000) NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT has_to_group_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.has_to_group;
