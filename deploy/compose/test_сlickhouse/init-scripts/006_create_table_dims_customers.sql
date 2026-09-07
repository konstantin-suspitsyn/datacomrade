-- Справочник клиентов в схеме DIMS.
-- customer_id совпадает с FACT.fact_sales.customer_id (диапазон 1..200) — ключ
-- связи факта и измерения. Атрибуты выводятся из cityHash64(customer_id, k) —
-- детерминированно. INSERT защищён проверкой на пустоту.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE TABLE IF NOT EXISTS DIMS.customers
(
    customer_id     Int64,
    customer_name   String,
    segment         LowCardinality(String),
    city            LowCardinality(String),
    registered_date Date
)
ENGINE = MergeTree
ORDER BY customer_id;

INSERT INTO DIMS.customers
SELECT
    number + 1                                                                          AS customer_id,
    concat('Клиент ', leftPad(toString(number + 1), 3, '0'))                            AS customer_name,
    ['retail', 'smb', 'enterprise'][(cityHash64(number + 1, 1) % 3) + 1]                AS segment,
    ['Москва', 'Санкт-Петербург', 'Новосибирск', 'Екатеринбург', 'Казань'][(cityHash64(number + 1, 2) % 5) + 1] AS city,
    toDate('2022-01-01') + (cityHash64(number + 1, 3) % 730)                            AS registered_date
FROM numbers(200)
WHERE (SELECT count() FROM DIMS.customers) = 0;
