-- +goose up
CREATE TABLE IF NOT EXISTS 
dc."user" (
	id bigserial NOT NULL,
	"name" varchar(512) NOT NULL, -- /microservice/username
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	is_deleted bool DEFAULT false NOT NULL,
	CONSTRAINT user_name_unique UNIQUE (name),
	CONSTRAINT user_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc."user";
