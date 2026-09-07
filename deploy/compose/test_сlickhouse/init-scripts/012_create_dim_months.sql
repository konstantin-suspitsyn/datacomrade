-- Справочник месяцев в схеме DIMS + словарь поверх него.
-- Источник — DISTINCT-срез DIMS.calendar по атрибутам месяца.
-- Ключ year_month (YYYYMM) — тот же, что в FACT.fact_sales_monthly.
-- Доступ через dictGet('DIMS.months_dict', <attr>, toUInt64(year_month)).
-- Выполняется один раз при инициализации пустого volume; INSERT защищён
-- проверкой на пустоту, CREATE OR REPLACE DICTIONARY — идемпотентен.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE TABLE IF NOT EXISTS DIMS.months
(
    year_month  Int32,
    year_num    Int32,
    quarter_num Int32,
    month_num   Int32,
    month_name  String
)
ENGINE = MergeTree
ORDER BY year_month;

INSERT INTO DIMS.months
SELECT DISTINCT
    year_month,
    year_num,
    quarter_num,
    month_num,
    month_name
FROM `DIMS`.calendar
WHERE (SELECT count() FROM DIMS.months) = 0;

CREATE OR REPLACE DICTIONARY DIMS.months_dict
(
    year_month  UInt64,
    year_num    Int32,
    quarter_num Int32,
    month_num   Int32,
    month_name  String
)
PRIMARY KEY year_month
SOURCE(CLICKHOUSE(HOST 'localhost' PORT 9000 USER 'query_engine_user' PASSWORD 'pass' DB 'DIMS' TABLE 'months'))
LIFETIME(MIN 300 MAX 600)
LAYOUT(HASHED());
