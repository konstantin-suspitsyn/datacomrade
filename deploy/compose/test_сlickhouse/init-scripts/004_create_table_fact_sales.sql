-- Фактовая таблица продаж в схеме FACT.
-- Колонки повторяют метаданные dc.column_cat для test_fact_sales
-- (datacatalogue/db/other_queries/test_sqls/insert_meaningful_data.sql, секция 13):
--   sale_id  - суррогатный ключ продажи;
--   date_id  - YYYYMMDD, ссылается на DIMS.calendar.date_id;
--   меры quantity / gross_amount / discount_amount / net_amount;
--   net_amount = gross_amount - discount_amount (согласованно на уровне строки).
-- Типы: string-атрибуты -> LowCardinality(String), decimal -> Decimal(18, 2).
-- Выполняется один раз при инициализации пустого volume; INSERT защищён
-- проверкой на пустоту, повторный ручной прогон ничего не дублирует.

CREATE DATABASE IF NOT EXISTS FACT;

CREATE TABLE IF NOT EXISTS FACT.fact_sales
(
    sale_id         Int64,
    date_id         Int64,
    sale_date       Date,
    product_id      Int64,
    customer_id     Int64,
    sales_channel   LowCardinality(String),
    payment_method  LowCardinality(String),
    quantity        Int32,
    gross_amount    Decimal(18, 2),
    discount_amount Decimal(18, 2),
    net_amount      Decimal(18, 2)
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(sale_date)
ORDER BY (sale_date, product_id, customer_id);

-- Тестовые данные: 20 000 продаж, детерминированно раскиданных по диапазону
-- дат DIMS.calendar (2024-01-01 .. 2025-12-31). Значения выводятся из
-- cityHash64(number, k) с разными k — так каждый столбец получает свою
-- независимую псевдослучайную величину (обычный rand() ClickHouse схлопнул
-- бы одинаковые вызовы в один столбец), и повторный прогон даёт тот же набор.
INSERT INTO FACT.fact_sales
SELECT
    n + 1                                                          AS sale_id,
    toInt64(formatDateTime(sale_date, '%Y%m%d'))                   AS date_id,
    sale_date,
    product_id,
    customer_id,
    ['retail', 'online', 'partner'][channel_idx + 1]               AS sales_channel,
    ['card', 'cash', 'transfer'][payment_idx + 1]                  AS payment_method,
    quantity,
    gross_amount,
    discount_amount,
    gross_amount - discount_amount                                 AS net_amount
FROM
(
    SELECT
        n,
        sale_date,
        quantity,
        channel_idx,
        payment_idx,
        product_id,
        customer_id,
        toDecimal64(quantity * unit_price, 2)                              AS gross_amount,
        toDecimal64(round(quantity * unit_price * discount_pct, 2), 2)     AS discount_amount
    FROM
    (
        SELECT
            number                                              AS n,
            toDate('2024-01-01') + (cityHash64(number, 1) % 731) AS sale_date,
            toInt32((cityHash64(number, 2) % 20) + 1)            AS quantity,
            toInt64((cityHash64(number, 3) % 4901) + 100)        AS unit_price,
            (cityHash64(number, 4) % 21) / 100.0                 AS discount_pct,
            toInt64((cityHash64(number, 5) % 50) + 1)            AS product_id,
            toInt64((cityHash64(number, 6) % 200) + 1)           AS customer_id,
            cityHash64(number, 7) % 3                            AS channel_idx,
            cityHash64(number, 8) % 3                            AS payment_idx
        FROM numbers(20000)
    )
)
WHERE (SELECT count() FROM FACT.fact_sales) = 0;
