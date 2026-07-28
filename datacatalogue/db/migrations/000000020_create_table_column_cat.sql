-- +goose up
CREATE TABLE IF NOT EXISTS
dc.column_cat (
	id bigserial NOT NULL,
	table_id bigint NOT NULL,
	"name" varchar(256) NOT NULL,
	alias_id bigint NOT NULL,
	column_type_id bigint NOT NULL,
	description varchar(1000) NOT NULL,
	calculation_type_id bigint NOT NULL, -- Здесь должно быть ограничение на стороне backend, чтобы сравнивать с БД, к которой привязана колонка
	is_deleted bool DEFAULT false NOT NULL,
	show_in_ui bool NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id bigint NOT NULL,
	CONSTRAINT column_cat_name_unique UNIQUE (table_id, "name")
);

-- +goose down
DROP TABLE IF EXISTS dc.column_cat;
