-- Справочник недель в схеме DIMS + словарь поверх него.
-- Таблица DIMS.weeks — ровно результат запроса
--   SELECT DISTINCT year_num, week_of_year, year_week FROM DIMS.calendar
-- (year_num здесь календарный, toYear). year_week = ISO_год*100 + ISO_неделя,
-- поэтому на границе ISO-года один year_week встречается с двумя year_num
-- (напр. 202501 — это 2024-12-30..31 и весь январь 2025): 106 строк, 105 year_week.
--
-- Словарь DIMS.weeks_dict требует уникальный ключ year_week, поэтому его
-- источник берёт year_num как ISO-год прямо из ключа (intDiv(year_week, 100))
-- и сворачивает до 105 строк. Так год в словаре согласован с неделей.
-- Доступ через dictGet('DIMS.weeks_dict', <attr>, toUInt64(year_week)).
--
-- Выполняется один раз при инициализации пустого volume; INSERT защищён
-- проверкой на пустоту, CREATE OR REPLACE DICTIONARY — идемпотентен.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE TABLE IF NOT EXISTS DIMS.weeks
(
    year_num     Int32,
    week_of_year Int32,
    year_week    Int32
)
ENGINE = MergeTree
ORDER BY (year_week, year_num);

INSERT INTO DIMS.weeks
SELECT DISTINCT
    year_num,
    week_of_year,
    year_week
FROM `DIMS`.calendar
WHERE (SELECT count() FROM DIMS.weeks) = 0;

CREATE OR REPLACE DICTIONARY DIMS.weeks_dict
(
    year_week    UInt64,
    year_num     Int32,
    week_of_year Int32
)
PRIMARY KEY year_week
SOURCE(CLICKHOUSE(
    HOST 'localhost' PORT 9000 USER 'query_engine_user' PASSWORD 'pass' DB 'DIMS'
    QUERY 'SELECT year_week, toInt32(intDiv(year_week, 100)) AS year_num, toInt32(year_week % 100) AS week_of_year FROM DIMS.weeks GROUP BY year_week'
))
LIFETIME(MIN 300 MAX 600)
LAYOUT(HASHED());
