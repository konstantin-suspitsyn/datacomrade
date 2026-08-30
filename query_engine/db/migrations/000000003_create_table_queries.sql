-- +goose up

CREATE TABLE IF NOT EXISTS queries.queries (
	id bigserial NOT NULL, -- Счетчик запросов
	query bytea NOT NULL, -- Бинарное представление запроса
	request jsonb NOT NULL, -- Исходный запрос в формате JSON
	database_id int8 NOT NULL, -- Идентификатор базы данных
	created_at timestamp DEFAULT now() NOT NULL, -- Создано
	updated_at timestamp DEFAULT now() NOT NULL, -- Обновлено
	created_by int8 NOT NULL, -- Пользователь, создавший запись
	CONSTRAINT queries_pk PRIMARY KEY (id)
);

-- +goose down

DROP TABLE IF EXISTS queries.queries;
