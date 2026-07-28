-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.group_levels (
	id bigserial NOT NULL,
	column_id int8 NOT NULL,
	parent_column_id int8 NOT NULL,
	"level" int2 NOT NULL,
	description varchar(1000) NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT group_levels_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.group_levels;
