-- +goose up

-- calculation_type_id становится необязательным: у колонки-атрибута измерения
-- (год, месяц, флаг выходного, элемент словаря) агрегата нет, и указывать
-- тип расчёта нечем. Внешний ключ column_cat_calculation_type_fk остаётся —
-- Postgres не проверяет его для NULL (как уже сделано с alias_id, миграция
-- 000000046).
ALTER TABLE dc.column_cat ALTER COLUMN calculation_type_id DROP NOT NULL;

-- +goose down

-- Откат возможен только если ни одной строки с пустым calculation_type_id
-- не осталось: подставлять произвольный тип расчёта значило бы приписать
-- колонке агрегат, которого у неё нет, поэтому миграция падает и требует
-- разобраться руками.
ALTER TABLE dc.column_cat ALTER COLUMN calculation_type_id SET NOT NULL;
