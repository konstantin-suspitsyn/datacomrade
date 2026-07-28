-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.host (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	description varchar(1000) NOT NULL,
	host_env varchar(255) NOT NULL, -- Название Env переменной для host
	port_env varchar(255) NOT NULL, -- Название Env переменной для порта
	username_env varchar(255) NOT NULL, -- Название Env переменной для имени пользователя
	password_env varchar(255) NOT NULL, -- Название Env переменной для пароля
	is_deleted bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	user_id int8 NOT NULL,
	CONSTRAINT host_env_unique UNIQUE (host_env),
	CONSTRAINT host_name_unique UNIQUE (name),
	CONSTRAINT host_password_env_unique UNIQUE (password_env),
	CONSTRAINT host_pk PRIMARY KEY (id),
	CONSTRAINT host_port_env_unique UNIQUE (port_env),
	CONSTRAINT host_user_name_unique UNIQUE (username_env)
);

-- +goose down

DROP TABLE IF EXISTS dc.host;
