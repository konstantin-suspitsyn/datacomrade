-- +goose up
CREATE TABLE IF NOT EXISTS 
dc.alias (
	id bigserial NOT NULL, -- Счетчик alias
	"name" varchar(255) NOT NULL, -- Имя alias должно быть уникально
	description varchar(1000) NOT NULL, -- Описание alias полей
	created_at timestamp DEFAULT now() NOT NULL, -- Создано
	updated_at timestamp DEFAULT now() NULL, -- Обновлено
	is_deleted bool DEFAULT false NOT NULL, -- Флаг удаления записи
	user_id int8 NOT NULL, -- Пользователь
	CONSTRAINT alias_name_unique UNIQUE (name),
	CONSTRAINT alias_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS dc.alias;
