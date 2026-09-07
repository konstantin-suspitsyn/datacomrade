-- Таблица-словарь календаря в схеме DIMS.
-- Колонки повторяют метаданные dc.column_cat для test_dim_calendar
-- (datacatalogue/db/other_queries/test_sqls/insert_meaningful_data.sql, секция 10):
--   date_id       - суррогатный ключ вида YYYYMMDD (Int64, как bigint в каталоге);
--   week_of_year  - номер ISO-недели 1..53;
--   week_of_month - номер недели внутри месяца 1..5 (7-дневными блоками от 1-го числа);
--   year_week     - год-неделя YYYYWW (ISO-год * 100 + ISO-неделя);
--   year_month    - год-месяц YYYYMM;
--   day_of_week   - 1..7, понедельник = 1 (совпадает с toDayOfWeek в ClickHouse).
-- Выполняется один раз при инициализации пустого volume; INSERT защищён
-- проверкой на пустоту, поэтому повторный ручной прогон ничего не дублирует.

CREATE DATABASE IF NOT EXISTS DIMS;

CREATE TABLE IF NOT EXISTS DIMS.calendar
(
    date_id       Int64,
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
    is_weekend    Bool,
    is_holiday    Bool,
    holiday_name  String
)
ENGINE = MergeTree
ORDER BY date_id;

-- Тестовые данные: сплошной диапазон дат 2024-01-01 .. 2025-12-31 (731 день).
-- Праздники — фиксированные нерабочие дни РФ по (месяц, день).
INSERT INTO DIMS.calendar
SELECT
    toInt64(formatDateTime(d, '%Y%m%d'))                                                       AS date_id,
    d                                                                                          AS calendar_date,
    toYear(d)                                                                                  AS year_num,
    toQuarter(d)                                                                               AS quarter_num,
    toMonth(d)                                                                                 AS month_num,
    ['January','February','March','April','May','June',
     'July','August','September','October','November','December'][toMonth(d)]                  AS month_name,
    toISOWeek(d)                                                                               AS week_of_year,
    intDiv(toDayOfMonth(d) - 1, 7) + 1                                                         AS week_of_month,
    toISOYear(d) * 100 + toISOWeek(d)                                                          AS year_week,
    toInt32(toYYYYMM(d))                                                                       AS year_month,
    toDayOfMonth(d)                                                                            AS day_of_month,
    toDayOfWeek(d)                                                                             AS day_of_week,
    ['Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday'][toDayOfWeek(d)]   AS day_name,
    toDayOfWeek(d) >= 6                                                                        AS is_weekend,
    holiday_name != ''                                                                         AS is_holiday,
    holiday_name
FROM
(
    SELECT
        toDate('2024-01-01') + number AS d,
        multiIf(
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (1, 1),   'Новый год',
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (1, 7),   'Рождество Христово',
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (2, 23),  'День защитника Отечества',
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (3, 8),   'Международный женский день',
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (5, 1),   'Праздник Весны и Труда',
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (5, 9),   'День Победы',
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (6, 12),  'День России',
            (toMonth(toDate('2024-01-01') + number), toDayOfMonth(toDate('2024-01-01') + number)) = (11, 4),  'День народного единства',
            ''
        ) AS holiday_name
    FROM numbers(731)
)
WHERE (SELECT count() FROM DIMS.calendar) = 0;
