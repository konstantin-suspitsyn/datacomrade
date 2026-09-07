-- Словарь клиентов: DIMS.customers_dict поверх таблицы DIMS.customers (файл 006).
-- Ключ customer_id (= FACT.fact_sales.customer_id). Доступ через
-- dictGet('DIMS.customers_dict', <attr>, toUInt64(customer_id)).
-- SOURCE — сам этот ClickHouse; креды совпадают с deploy/compose/test_сlickhouse/.env.
-- CREATE OR REPLACE: повторный прогон пересобирает определение словаря.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE OR REPLACE DICTIONARY DIMS.customers_dict
(
    customer_id     UInt64,
    customer_name   String,
    segment         String,
    city            String,
    registered_date Date
)
PRIMARY KEY customer_id
SOURCE(CLICKHOUSE(HOST 'localhost' PORT 9000 USER 'query_engine_user' PASSWORD 'pass' DB 'DIMS' TABLE 'customers'))
LIFETIME(MIN 300 MAX 600)
LAYOUT(HASHED());
