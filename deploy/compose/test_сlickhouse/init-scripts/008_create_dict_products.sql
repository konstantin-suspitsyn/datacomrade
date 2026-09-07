-- Словарь товаров: DIMS.products_dict поверх таблицы DIMS.products (файл 005).
-- Ключ product_id (= FACT.fact_sales.product_id). Доступ через
-- dictGet('DIMS.products_dict', <attr>, toUInt64(product_id)).
-- SOURCE — сам этот ClickHouse; креды совпадают с deploy/compose/test_сlickhouse/.env.
-- CREATE OR REPLACE: повторный прогон пересобирает определение словаря.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE OR REPLACE DICTIONARY DIMS.products_dict
(
    product_id   UInt64,
    product_name String,
    category     String,
    brand        String
)
PRIMARY KEY product_id
SOURCE(CLICKHOUSE(HOST 'localhost' PORT 9000 USER 'query_engine_user' PASSWORD 'pass' DB 'DIMS' TABLE 'products'))
LIFETIME(MIN 300 MAX 600)
LAYOUT(HASHED());
