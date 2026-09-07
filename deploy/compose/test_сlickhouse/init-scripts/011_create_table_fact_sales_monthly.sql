-- Месячный агрегат продаж в схеме FACT.
-- Строится из FACT.fact_sales, свёрнутого по месяцу (year_month из DIMS.calendar)
-- и разрезам товар / клиент / канал / способ оплаты.
-- Грань строки: year_month × product_id × customer_id × sales_channel × payment_method.
-- Выполняется один раз при инициализации пустого volume; INSERT защищён
-- проверкой на пустоту, повторный ручной прогон ничего не дублирует.

CREATE DATABASE IF NOT EXISTS FACT;

CREATE TABLE IF NOT EXISTS FACT.fact_sales_monthly
(
    year_month      Int32,
    product_id      Int64,
    customer_id     Int64,
    sales_channel   LowCardinality(String),
    payment_method  LowCardinality(String),
    quantity        Int64,
    gross_amount    Decimal(18, 2),
    discount_amount Decimal(18, 2),
    net_amount      Decimal(18, 2)
)
ENGINE = MergeTree
ORDER BY (year_month, product_id, customer_id, sales_channel, payment_method);

INSERT INTO FACT.fact_sales_monthly
SELECT
    d.year_month                AS year_month,
    s.product_id                AS product_id,
    s.customer_id               AS customer_id,
    s.sales_channel             AS sales_channel,
    s.payment_method            AS payment_method,
    sum(s.quantity)             AS quantity,
    sum(s.gross_amount)         AS gross_amount,
    sum(s.discount_amount)      AS discount_amount,
    sum(s.net_amount)           AS net_amount
FROM `FACT`.fact_sales AS s
LEFT JOIN `DIMS`.calendar AS d ON s.date_id = d.date_id
WHERE (SELECT count() FROM `FACT`.fact_sales_monthly) = 0
GROUP BY
    year_month,
    product_id,
    customer_id,
    sales_channel,
    payment_method;
