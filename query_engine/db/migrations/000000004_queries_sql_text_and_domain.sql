-- +goose up

-- query теперь хранит финальный параметризованный SQL (text), а не бинарное
-- представление AST. Старых строк в таблице нет, приведение типа безопасно.
ALTER TABLE queries.queries
	ALTER COLUMN query TYPE text USING convert_from(query, 'UTF8');

COMMENT ON COLUMN queries.queries.query IS 'Финальный параметризованный SQL';

ALTER TABLE queries.queries
	ADD COLUMN IF NOT EXISTS params jsonb NOT NULL DEFAULT '[]'::jsonb; -- Значения плейсхолдеров SQL по порядку ($1, $2, ...)

ALTER TABLE queries.queries
	ADD COLUMN IF NOT EXISTS domain_id int8 NOT NULL; -- Бизнес-домен, в котором запрос был резолвлен

-- +goose down

ALTER TABLE queries.queries
	DROP COLUMN IF EXISTS domain_id;

ALTER TABLE queries.queries
	DROP COLUMN IF EXISTS params;

ALTER TABLE queries.queries
	ALTER COLUMN query TYPE bytea USING convert_to(query, 'UTF8');

COMMENT ON COLUMN queries.queries.query IS 'Бинарное представление запроса';
