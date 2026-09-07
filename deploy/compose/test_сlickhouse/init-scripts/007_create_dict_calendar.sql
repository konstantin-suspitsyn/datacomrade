-- Словарь календаря: DIMS.calendar_dict поверх таблицы DIMS.calendar (файл 003).
-- Ключ date_id (YYYYMMDD). Доступ через dictGet('DIMS.calendar_dict', <attr>, toUInt64(date_id)).
-- LAYOUT(HASHED): ключи разрежены (~20 млн), FLAT выделил бы огромный массив.
-- SOURCE — сам этот ClickHouse; креды совпадают с deploy/compose/test_сlickhouse/.env
-- (TEST_CLICKHOUSE_USER / TEST_CLICKHOUSE_PASSWORD). Self-referencing источник
-- требует явной аутентификации.
-- CREATE OR REPLACE: повторный прогон пересобирает определение словаря.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE OR REPLACE DICTIONARY DIMS.calendar_dict
(
    date_id       UInt64,
    calendar_date Date,
    year_num      Int32,
    quarter_num   Int32,
    month_num     Int32,
    month_name    String,
    week_of_year  Int32,
    week_of_month Int32,
    year_week     Int32,
    year_month    Int32,
    day_of_month  Int32,
    day_of_week   Int32,
    day_name      String,
    is_weekend    UInt8,
    is_holiday    UInt8,
    holiday_name  String
)
PRIMARY KEY date_id
SOURCE(CLICKHOUSE(HOST 'localhost' PORT 9000 USER 'query_engine_user' PASSWORD 'pass' DB 'DIMS' TABLE 'calendar'))
LIFETIME(MIN 300 MAX 600)
LAYOUT(HASHED());
