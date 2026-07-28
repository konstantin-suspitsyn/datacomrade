-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.domain_cat (
	id bigserial NOT NULL,
	domain_name varchar(100) NOT NULL, -- Имя домена, к которому относится таблица
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT domain_cat_name_unique UNIQUE (domain_name),
	CONSTRAINT domain_cat_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.domain_cat;
