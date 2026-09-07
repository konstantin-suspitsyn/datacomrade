-- Справочник товаров в схеме DIMS.
-- product_id совпадает с FACT.fact_sales.product_id (диапазон 1..50) — это ключ
-- связи факта и измерения (join / dictGet; внешних ключей в ClickHouse нет).
-- Атрибуты выводятся из cityHash64(product_id, k) — детерминированно,
-- повторный прогон даёт тот же набор. INSERT защищён проверкой на пустоту.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE TABLE IF NOT EXISTS DIMS.products
(
    product_id   Int64,
    product_name String,
    category     LowCardinality(String),
    brand        LowCardinality(String)
)
ENGINE = MergeTree
ORDER BY product_id;

INSERT INTO DIMS.products
SELECT
    number + 1                                                                        AS product_id,
    concat('Товар ', leftPad(toString(number + 1), 2, '0'))                           AS product_name,
    ['Электроника', 'Одежда', 'Дом', 'Спорт', 'Книги'][(cityHash64(number + 1, 1) % 5) + 1] AS category,
    ['Brand A', 'Brand B', 'Brand C', 'Brand D', 'Brand E'][(cityHash64(number + 1, 2) % 5) + 1] AS brand
FROM numbers(50)
WHERE (SELECT count() FROM DIMS.products) = 0;
