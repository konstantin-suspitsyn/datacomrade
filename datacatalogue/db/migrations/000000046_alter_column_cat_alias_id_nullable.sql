-- +goose up

-- alias_id становится необязательным: колонка каталога может не иметь алиаса,
-- и тогда сверять её column_type_id не с чем. Внешний ключ остаётся —
-- Postgres не проверяет его для NULL.
ALTER TABLE dc.column_cat ALTER COLUMN alias_id DROP NOT NULL;

-- +goose down

-- Откат возможен только если ни одной строки без алиаса не осталось:
-- подставлять сюда произвольный alias_id значило бы приписать колонке
-- чужой смысл, поэтому миграция падает и требует разобраться руками.
ALTER TABLE dc.column_cat ALTER COLUMN alias_id SET NOT NULL;
