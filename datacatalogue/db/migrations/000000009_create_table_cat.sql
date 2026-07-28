-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.table_cat (
	id bigserial NOT NULL,
	"name" varchar(128) NOT NULL,
	description varchar(2000) NOT NULL,
	schema_id int8 NOT NULL,
	table_type_id int8 NOT NULL,
	domain_id int8 NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
    is_get_dict bool DEFAULT false NOT NULL, -- это только для ClickHouse. Чтобы использовать getDict вместо join
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT table_cat_name_unique UNIQUE (schema_id, name),
	CONSTRAINT table_cat_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.table_cat;
