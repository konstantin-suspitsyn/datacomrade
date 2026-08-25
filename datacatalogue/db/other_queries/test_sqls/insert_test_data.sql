-- =========================================================================
-- Тестовые данные для dc.*
-- =========================================================================
-- Соглашения (важны для скрипта удаления delete_test_data.sql):
--
--   1. Все именованные сущности создаются с префиксом `test_`.
--      Удаление идёт строго по LIKE 'test\_%' (обратный слеш экранирует `_`,
--      иначе он был бы wildcard-ом).
--   2. Сущности без собственного имени (dc.column_cat, dc.group_levels,
--      dc.has_to_group, связки ролей, dc.*_calculation) удаляются только
--      через связь с родителем, помеченным префиксом `test_`.
--   3. Скрипт НИЧЕГО не вставляет в справочники, засеянные миграциями:
--      dc.calculation_type, dc.table_type, dc.column_type, а также сами
--      роли can_read/can_write/can_grant в dc.domain_roles и dc.table_roles —
--      новые строки в эти две таблицы добавлять нельзя, используются только
--      существующие. Читаются по имени, привязки к доменам/таблицам/
--      пользователям — обычные тестовые данные (см. секции 14–15).
--      dc.user_domain_roles и dc.user_table_roles хранят домен/таблицу
--      прямо в строке (domain_id / table_id), отдельных связок
--      domain <-> role и table <-> role в схеме больше нет.
--   4. Нигде не используются захардкоженные id: все FK разрешаются
--      подзапросами по уникальным именам, поэтому скрипт не зависит
--      от текущих значений sequence.
--   5. Скрипт идемпотентен: повторный запуск ничего не дублирует
--      (ON CONFLICT DO NOTHING там, где есть unique-констрейнт,
--      NOT EXISTS — там, где его нет).
--
-- Владелец всех тестовых записей — dc."user".name = 'test_user_admin',
-- чтобы удаление тестовых пользователей не упиралось в FK ON DELETE RESTRICT.
--
-- -------------------------------------------------------------------------
-- Объём набора (проверяется секцией 16, до COMMIT):
--
--   * не меньше 300 тестовых таблиц (14 именованных + 286 сгенерированных);
--   * не меньше 30 колонок в каждой тестовой таблице;
--   * не меньше 10 колонок с show_in_ui = true в каждой тестовой таблице.
--
-- Бизнес-правила dc.column_cat (реализованы в ColumnCatGuard,
-- datacatalogue/internal/api/tablesapiv1/column_cat_guard.go) — набор им
-- подчиняется, иначе через API эти же строки завести было бы нельзя:
--
--   * show_in_ui = true требует alias_id: колонке без алиаса нечего
--     показывать в интерфейсе. column_type_id у колонки NOT NULL на уровне
--     схемы, так что у видимой колонки заполнены оба поля;
--   * если alias_id заполнен, dc.alias.column_type_id обязан совпадать
--     с dc.column_cat.column_type_id — иначе одна колонка описана двумя
--     разными типами. Поэтому у каждого test_alias_* ровно один тип
--     (секция 7), и колонка берёт alias и column_type из одной строки
--     шаблона (секции 9 и 9.1);
--   * alias_id необязателен (миграция 000000046) — набор специально
--     содержит колонки без алиаса, но все они скрытые (show_in_ui = false).
-- =========================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- 1. Пользователи
-- -------------------------------------------------------------------------
-- test_user_admin     — владелец всех тестовых объектов; can_read + can_write
--                       на домен sales, can_read на таблицу payments
-- test_user_reader    — can_read на домен sales, can_read на таблицу payments
--                       (payments — из домена finance, доступ только через
--                       таблицу, см. секции 14–15)
-- test_user_writer    — can_write на домен sales; к finance/payments доступа
--                       нет ни доменного, ни табличного
-- test_user_no_roles  — без единой роли (негативные кейсы прав доступа)
INSERT INTO dc."user" ("name", external_id, created_at, updated_at, is_deleted)
VALUES ('test_user_admin', 'aaaaaaaa-0000-4000-8000-000000000001', now(), now(), false),
       ('test_user_reader', 'aaaaaaaa-0000-4000-8000-000000000002', now(), now(), false),
       ('test_user_writer', 'aaaaaaaa-0000-4000-8000-000000000003', now(), now(), false),
       ('test_user_no_roles', 'aaaaaaaa-0000-4000-8000-000000000004', now(), now(), false)
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 2. Хосты
-- -------------------------------------------------------------------------
INSERT INTO dc.host ("name", description, host_env, port_env, username_env, password_env,
                     is_deleted, created_at, updated_at, user_id)
SELECT v.name, v.description, v.host_env, v.port_env, v.username_env, v.password_env,
       false, now(), now(), u.id
FROM (VALUES ('test_host_pg', 'Тестовый хост PostgreSQL',
              'TEST_PG_HOST', 'TEST_PG_PORT', 'TEST_PG_USER', 'TEST_PG_PASSWORD'),
             ('test_host_ch', 'Тестовый хост ClickHouse',
              'TEST_CH_HOST', 'TEST_CH_PORT', 'TEST_CH_USER', 'TEST_CH_PASSWORD')
     ) AS v(name, description, host_env, port_env, username_env, password_env)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 3. Типы БД
-- -------------------------------------------------------------------------
INSERT INTO dc.database_type ("name", db_version, is_deleted, created_at, updated_at, user_id)
SELECT v.name, v.db_version, false, now(), now(), u.id
FROM (VALUES ('test_dbtype_postgres', '16.4'),
             ('test_dbtype_clickhouse', '24.3')
     ) AS v(name, db_version)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 4. Базы данных
-- -------------------------------------------------------------------------
INSERT INTO dc.database_cat ("name", host_id, database_type_id, description,
                             is_deleted, created_at, updated_at, user_id)
SELECT v.name, h.id, dt.id, v.description, false, now(), now(), u.id
FROM (VALUES ('test_db_oltp', 'test_host_pg', 'test_dbtype_postgres', 'Тестовая OLTP-база'),
             ('test_db_dwh', 'test_host_ch', 'test_dbtype_clickhouse', 'Тестовое хранилище')
     ) AS v(name, host_name, dbtype_name, description)
         JOIN dc.host h ON h."name" = v.host_name
         JOIN dc.database_type dt ON dt."name" = v.dbtype_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 5. Схемы (в dc.schema_cat нет unique-констрейнта -> защита через NOT EXISTS)
-- -------------------------------------------------------------------------
INSERT INTO dc.schema_cat (database_id, "name", is_deleted, created_at, updated_at, user_id)
SELECT d.id, v.name, false, now(), now(), u.id
FROM (VALUES ('test_schema_public', 'test_db_oltp'),
             ('test_schema_marts', 'test_db_dwh')
     ) AS v(name, db_name)
         JOIN dc.database_cat d ON d."name" = v.db_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
WHERE NOT EXISTS (SELECT 1
                  FROM dc.schema_cat s
                  WHERE s.database_id = d.id
                    AND s."name" = v.name);

-- -------------------------------------------------------------------------
-- 6. Домены
-- -------------------------------------------------------------------------
INSERT INTO dc.domain_cat (domain_name, is_deleted, created_at, updated_at, user_id)
SELECT v.domain_name, false, now(), now(), u.id
FROM (VALUES ('test_domain_sales'),
             ('test_domain_finance')
     ) AS v(domain_name)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 7. Алиасы
-- -------------------------------------------------------------------------
-- У алиаса есть собственный column_type_id (миграция 000000045), и он задаёт
-- тип для каждой колонки, которая на этот алиас ссылается: правило
-- «alias.column_type_id = column_cat.column_type_id» проверяется в
-- ColumnCatGuard. Поэтому здесь у каждого test_alias_* ровно один тип,
-- а шаблоны колонок ниже берут alias и column_type парой из одной строки.
--
-- Типы колонок ('int', 'bigint', 'decimal', ...) читаются из dc.column_type,
-- засеянного миграцией 000000041; новых строк туда скрипт не добавляет.
-- Набор из 32 алиасов покрывает все 10 типов колонок.
--
-- ON CONFLICT ... DO UPDATE, а не DO NOTHING: если алиас с таким именем уже
-- лежит в базе с другим типом (например, от прежней версии скрипта, где
-- column_type_id брался из DEFAULT), его надо починить — иначе колонки,
-- ссылающиеся на него, нарушат правило совпадения типов.
INSERT INTO dc.alias ("name", description, column_type_id, created_at, updated_at, is_deleted, user_id)
SELECT v.name, v.description, ct.id, now(), now(), false, u.id
FROM (VALUES ('test_alias_id', 'Идентификатор', 'bigint'),
             ('test_alias_geo', 'География', 'bigint'),
             ('test_alias_name', 'Наименование', 'string'),
             ('test_alias_code', 'Код', 'string'),
             ('test_alias_status', 'Статус', 'string'),
             ('test_alias_currency', 'Валюта', 'string'),
             ('test_alias_country', 'Страна', 'string'),
             ('test_alias_city', 'Город', 'string'),
             ('test_alias_category', 'Категория', 'string'),
             ('test_alias_type', 'Тип', 'string'),
             ('test_alias_version', 'Версия', 'string'),
             ('test_alias_tag', 'Тег', 'string'),
             ('test_alias_channel', 'Канал', 'string'),
             ('test_alias_source', 'Источник', 'string'),
             ('test_alias_email', 'Email', 'string'),
             ('test_alias_phone', 'Телефон', 'string'),
             ('test_alias_amount', 'Сумма', 'decimal'),
             ('test_alias_price', 'Цена за единицу', 'decimal'),
             ('test_alias_qty', 'Количество (штуки)', 'int'),
             ('test_alias_score', 'Оценка', 'int'),
             ('test_alias_duration', 'Длительность', 'int'),
             ('test_alias_priority', 'Приоритет', 'int'),
             ('test_alias_percent', 'Процент', 'float'),
             ('test_alias_weight', 'Вес', 'float'),
             ('test_alias_volume', 'Объём', 'float'),
             ('test_alias_rating', 'Рейтинг', 'float'),
             ('test_alias_flag', 'Признак (флаг)', 'bool'),
             ('test_alias_url', 'Ссылка (URL)', 'text'),
             ('test_alias_comment', 'Комментарий', 'text'),
             ('test_alias_dt', 'Дата', 'date'),
             ('test_alias_created_at', 'Дата и время создания', 'datetime'),
             ('test_alias_event_ts', 'Момент события с часовым поясом', 'datetime_tz')
     ) AS v(name, description, coltype_name)
         JOIN dc.column_type ct ON ct."name" = v.coltype_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT alias_name_unique
    DO UPDATE SET description    = EXCLUDED.description,
                  column_type_id = EXCLUDED.column_type_id,
                  updated_at     = now();

-- -------------------------------------------------------------------------
-- 8. Таблицы
-- -------------------------------------------------------------------------
-- Типы таблиц ('fact', 'dimension') берутся из миграции 000000038, не создаются заново.
INSERT INTO dc.table_cat ("name", description, schema_id, table_type_id, domain_id,
                          is_deleted, is_get_dict, created_at, updated_at, user_id)
SELECT v.name, v.description, s.id, tt.id, dom.id, false, v.is_get_dict, now(), now(), u.id
FROM (VALUES ('test_tbl_orders', 'Тестовые заказы', 'test_schema_marts', 'fact', 'test_domain_sales', false),
             ('test_tbl_clients', 'Тестовые клиенты', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_geo', 'Тестовый справочник географии', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_payments', 'Тестовые платежи', 'test_schema_public', 'fact', 'test_domain_finance', false),
             -- +10 таблиц для проверки разных прав доступа по секции 15
             ('test_tbl_products', 'Тестовые товары', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_suppliers', 'Тестовые поставщики', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_inventory', 'Тестовые остатки на складе', 'test_schema_marts', 'fact', 'test_domain_sales', false),
             ('test_tbl_shipments', 'Тестовые отгрузки', 'test_schema_marts', 'fact', 'test_domain_sales', false),
             ('test_tbl_returns', 'Тестовые возвраты товаров', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_invoices', 'Тестовые счета', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_refunds', 'Тестовые возвраты платежей', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_employees', 'Тестовые сотрудники', 'test_schema_public', 'dimension', 'test_domain_sales', true),
             ('test_tbl_campaigns', 'Тестовые маркетинговые кампании', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_regions_ext', 'Тестовый расширенный справочник регионов', 'test_schema_marts', 'dimension', 'test_domain_sales', true)
     ) AS v(name, description, schema_name, table_type_name, domain_name, is_get_dict)
         JOIN dc.schema_cat s ON s."name" = v.schema_name
         JOIN dc.table_type tt ON tt."name" = v.table_type_name
         JOIN dc.domain_cat dom ON dom.domain_name = v.domain_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 8.1. Сгенерированные таблицы (test_tbl_gen_001 .. test_tbl_gen_286)
-- -------------------------------------------------------------------------
-- 14 именованных таблиц выше + 286 сгенерированных = 300 таблиц.
-- Раскладка по схемам/типам/доменам детерминирована (зависит только от номера),
-- поэтому повторный прогон даёт ровно те же строки, а ON CONFLICT DO NOTHING
-- гасит дубли по (schema_id, name).
INSERT INTO dc.table_cat ("name", description, schema_id, table_type_id, domain_id,
                          is_deleted, is_get_dict, created_at, updated_at, user_id)
SELECT v.name, v.description, s.id, tt.id, dom.id, false, v.is_get_dict, now(), now(), u.id
FROM (SELECT 'test_tbl_gen_' || to_char(g.i, 'FM000')                                      AS name,
             'Тестовая сгенерированная таблица ' || to_char(g.i, 'FM000')                  AS description,
             CASE WHEN g.i % 2 = 0 THEN 'test_schema_marts' ELSE 'test_schema_public' END  AS schema_name,
             CASE WHEN g.i % 2 = 0 THEN 'dimension' ELSE 'fact' END                        AS table_type_name,
             CASE WHEN g.i % 3 = 0 THEN 'test_domain_finance' ELSE 'test_domain_sales' END AS domain_name,
             (g.i % 2 = 0)                                                                 AS is_get_dict
      FROM generate_series(1, 286) AS g(i)) AS v
         JOIN dc.schema_cat s ON s."name" = v.schema_name
         JOIN dc.table_type tt ON tt."name" = v.table_type_name
         JOIN dc.domain_cat dom ON dom.domain_name = v.domain_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 9. Колонки именованных таблиц (бизнес-специфика)
-- -------------------------------------------------------------------------
-- Здесь только те колонки, на которые ссылаются секции 10, 11 и 13
-- (иерархия группировок, связи между таблицами, доступные расчёты).
-- Добор до 30+ колонок в каждой таблице делает шаблон из секции 9.1.
--
-- Типы расчёта ('sum', 'count', 'count_distinct', ...) берутся из миграции 000000034.
-- Типы колонок ('bigint', 'string', 'date', 'decimal', ...) — из миграции 000000041.
-- Пара (алиас, тип колонки) в каждой строке согласована с секцией 7: алиас
-- test_alias_id всегда идёт с 'bigint', test_alias_amount — с 'decimal' и т. д.
--
-- Секция идёт ПЕРЕД шаблоном: при совпадении имён (например, orders.amount
-- есть и здесь, и в шаблоне) в таблице останется строка с бизнес-описанием,
-- а шаблонный дубль погасит ON CONFLICT DO NOTHING в секции 9.1.
INSERT INTO dc.column_cat (table_id, "name", alias_id, column_type_id, description,
                           calculation_type_id, is_deleted, show_in_ui,
                           created_at, updated_at, user_id)
SELECT t.id, v.name, a.id, ct.id, v.description, calc.id, false, v.show_in_ui, now(), now(), u.id
FROM (VALUES ('test_tbl_orders', 'order_id', 'Идентификатор заказа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_orders', 'client_id', 'Идентификатор клиента', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_orders', 'order_dt', 'Дата заказа', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_orders', 'amount', 'Сумма заказа', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_clients', 'client_id', 'Идентификатор клиента', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_clients', 'client_name', 'Имя клиента', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_clients', 'city_id', 'Идентификатор города', 'test_alias_geo', 'bigint', 'count_distinct', false),
             ('test_tbl_clients', 'client_type', 'Тип клиента', 'test_alias_type', 'string', 'count', true),
             ('test_tbl_geo', 'city_id', 'Идентификатор города', 'test_alias_geo', 'bigint', 'count_distinct', true),
             ('test_tbl_geo', 'region_id', 'Идентификатор региона', 'test_alias_geo', 'bigint', 'count_distinct', true),
             ('test_tbl_geo', 'geo_name', 'Название географии', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_geo', 'geo_code', 'Код географии', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_payments', 'payment_id', 'Идентификатор платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_payments', 'order_id', 'Идентификатор заказа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_payments', 'amount', 'Сумма платежа', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_payments', 'payment_dt', 'Дата платежа', 'test_alias_dt', 'date', 'max', true),
             -- Колонки для 10 таблиц из секции 8 (products..regions_ext)
             ('test_tbl_products', 'product_id', 'Идентификатор товара', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_products', 'product_name', 'Название товара', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_products', 'supplier_id', 'Идентификатор поставщика', 'test_alias_id', 'bigint', 'count_distinct', false),
             ('test_tbl_suppliers', 'supplier_id', 'Идентификатор поставщика', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_suppliers', 'supplier_name', 'Название поставщика', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_inventory', 'inventory_id', 'Идентификатор записи остатка', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_inventory', 'product_id', 'Идентификатор товара', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_inventory', 'qty', 'Количество на складе', 'test_alias_qty', 'int', 'sum', true),
             ('test_tbl_inventory', 'inventory_dt', 'Дата остатка', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_inventory', 'warehouse_code', 'Код склада', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_shipments', 'shipment_id', 'Идентификатор отгрузки', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_shipments', 'order_id', 'Идентификатор заказа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_shipments', 'shipment_dt', 'Дата отгрузки', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_shipments', 'amount', 'Сумма отгрузки', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_returns', 'return_id', 'Идентификатор возврата', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_returns', 'order_id', 'Идентификатор заказа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_returns', 'amount', 'Сумма возврата', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_returns', 'return_dt', 'Дата возврата', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_returns', 'reason', 'Причина возврата', 'test_alias_comment', 'text', 'count', true),
             ('test_tbl_invoices', 'invoice_id', 'Идентификатор счёта', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_invoices', 'payment_id', 'Идентификатор платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_invoices', 'amount', 'Сумма счёта', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_invoices', 'invoice_dt', 'Дата счёта', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_refunds', 'refund_id', 'Идентификатор возврата платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_refunds', 'payment_id', 'Идентификатор платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_refunds', 'amount', 'Сумма возврата платежа', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_refunds', 'refund_dt', 'Дата возврата платежа', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_employees', 'employee_id', 'Идентификатор сотрудника', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_employees', 'employee_name', 'Имя сотрудника', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_campaigns', 'campaign_id', 'Идентификатор кампании', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_campaigns', 'campaign_name', 'Название кампании', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_campaigns', 'campaign_dt', 'Дата старта кампании', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_campaigns', 'budget', 'Бюджет кампании', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_regions_ext', 'region_id', 'Идентификатор региона', 'test_alias_geo', 'bigint', 'count_distinct', true),
             ('test_tbl_regions_ext', 'region_name', 'Расширенное название региона', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_regions_ext', 'region_code', 'Код региона', 'test_alias_code', 'string', 'count', true)
     ) AS v(table_name, name, description, alias_name, coltype_name, calc_name, show_in_ui)
         JOIN dc.table_cat t ON t."name" = v.table_name
         JOIN dc.alias a ON a."name" = v.alias_name
         JOIN dc.column_type ct ON ct."name" = v.coltype_name
         JOIN dc.calculation_type calc ON calc."name" = v.calc_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 9.1. Базовый шаблон колонок для всех тестовых таблиц
-- -------------------------------------------------------------------------
-- Шаблон из 32 позиций разворачивается на КАЖДУЮ таблицу с именем test_%,
-- то есть и на 14 именованных, и на 286 сгенерированных. Отсюда берётся
-- требуемый объём:
--
--   * 32 колонки в каждой таблице (минимум по требованию — 30);
--   * позиции j = 1..20 идут с show_in_ui = true — 20 видимых колонок
--     в каждой таблице (минимум по требованию — 10);
--   * позиции j = 21..27 — скрытые, но с алиасом;
--   * позиции j = 28..32 — скрытые и без алиаса (alias_id IS NULL,
--     разрешено миграцией 000000046). Правило «show_in_ui = true требует
--     alias_id» не нарушается: без алиаса тут только скрытые колонки.
--
-- Пара (alias_name, coltype_name) в каждой строке шаблона согласована
-- с секцией 7 — иначе сломается правило совпадения типов колонки и алиаса.
-- У позиций без алиаса column_type_id всё равно заполнен: он NOT NULL
-- на уровне схемы и не зависит от наличия алиаса.
--
-- Имена колонок совпадают у всех таблиц, поэтому у именованных таблиц часть
-- шаблонных позиций гасится ON CONFLICT DO NOTHING — там уже стоит строка
-- из секции 9 с тем же именем. Итоговое число колонок именованной таблицы =
-- 32 + (её собственные колонки, которых нет в шаблоне), то есть всегда > 30.
INSERT INTO dc.column_cat (table_id, "name", alias_id, column_type_id, description,
                           calculation_type_id, is_deleted, show_in_ui,
                           created_at, updated_at, user_id)
SELECT t.id,
       tpl.col_name,
       a.id,
       ct.id,
       tpl.col_description || ' (' || t."name" || ')',
       calc.id,
       false,
       tpl.j <= 20,
       now(), now(), u.id
FROM dc.table_cat t
         CROSS JOIN (VALUES (1, 'id', 'Идентификатор', 'test_alias_id', 'bigint', 'count_distinct'),
                            (2, 'name', 'Наименование', 'test_alias_name', 'string', 'count'),
                            (3, 'code', 'Код', 'test_alias_code', 'string', 'count'),
                            (4, 'status', 'Статус', 'test_alias_status', 'string', 'count'),
                            (5, 'amount', 'Сумма', 'test_alias_amount', 'decimal', 'sum'),
                            (6, 'dt', 'Дата', 'test_alias_dt', 'date', 'max'),
                            (7, 'qty', 'Количество', 'test_alias_qty', 'int', 'sum'),
                            (8, 'price', 'Цена за единицу', 'test_alias_price', 'decimal', 'avg'),
                            (9, 'is_active', 'Признак активности', 'test_alias_flag', 'bool', 'count'),
                            (10, 'email', 'Email', 'test_alias_email', 'string', 'count'),
                            (11, 'phone', 'Телефон', 'test_alias_phone', 'string', 'count'),
                            (12, 'created_at', 'Дата и время создания', 'test_alias_created_at', 'datetime', 'max'),
                            (13, 'event_ts', 'Момент события', 'test_alias_event_ts', 'datetime_tz', 'max'),
                            (14, 'percent', 'Процент', 'test_alias_percent', 'float', 'avg'),
                            (15, 'weight', 'Вес', 'test_alias_weight', 'float', 'sum'),
                            (16, 'volume', 'Объём', 'test_alias_volume', 'float', 'sum'),
                            (17, 'currency', 'Валюта', 'test_alias_currency', 'string', 'count'),
                            (18, 'country', 'Страна', 'test_alias_country', 'string', 'count'),
                            (19, 'city', 'Город', 'test_alias_city', 'string', 'count'),
                            (20, 'category', 'Категория', 'test_alias_category', 'string', 'count'),
                            -- Скрытые колонки с алиасом
                            (21, 'type', 'Тип', 'test_alias_type', 'string', 'count'),
                            (22, 'score', 'Оценка', 'test_alias_score', 'int', 'avg'),
                            (23, 'rating', 'Рейтинг', 'test_alias_rating', 'float', 'avg'),
                            (24, 'duration', 'Длительность', 'test_alias_duration', 'int', 'sum'),
                            (25, 'version', 'Версия', 'test_alias_version', 'string', 'count'),
                            (26, 'geo_id', 'Идентификатор географии', 'test_alias_geo', 'bigint', 'count_distinct'),
                            (27, 'url', 'Ссылка', 'test_alias_url', 'text', 'count'),
                            -- Скрытые колонки без алиаса (alias_id IS NULL)
                            (28, 'comment', 'Комментарий', NULL, 'text', 'count'),
                            (29, 'tag', 'Тег', NULL, 'string', 'count'),
                            (30, 'priority', 'Приоритет', NULL, 'int', 'max'),
                            (31, 'channel', 'Канал', NULL, 'string', 'count'),
                            (32, 'source', 'Источник', NULL, 'string', 'count')
                    ) AS tpl(j, col_name, col_description, alias_name, coltype_name, calc_name)
         LEFT JOIN dc.alias a ON a."name" = tpl.alias_name
         JOIN dc.column_type ct ON ct."name" = tpl.coltype_name
         JOIN dc.calculation_type calc ON calc."name" = tpl.calc_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
WHERE t."name" LIKE 'test\_%'
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 10. Иерархия группировок (dc.group_levels)
-- -------------------------------------------------------------------------
-- region_id — корень (ссылается сам на себя, т.к. parent_column_id NOT NULL),
-- city_id — уровень ниже.
INSERT INTO dc.group_levels (column_id, parent_column_id, "level", description,
                             created_at, updated_at, is_deleted, user_id)
SELECT c.id, p.id, v.level, v.description, now(), now(), false, u.id
FROM (VALUES ('test_tbl_geo', 'region_id', 'test_tbl_geo', 'region_id', 1::int2, 'Тестовый уровень: регион'),
             ('test_tbl_geo', 'city_id', 'test_tbl_geo', 'region_id', 2::int2, 'Тестовый уровень: город внутри региона')
     ) AS v(table_name, column_name, parent_table_name, parent_column_name, level, description)
         JOIN dc.table_cat t ON t."name" = v.table_name
         JOIN dc.column_cat c ON c.table_id = t.id AND c."name" = v.column_name
         JOIN dc.table_cat pt ON pt."name" = v.parent_table_name
         JOIN dc.column_cat p ON p.table_id = pt.id AND p."name" = v.parent_column_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
WHERE NOT EXISTS (SELECT 1
                  FROM dc.group_levels gl
                  WHERE gl.column_id = c.id
                    AND gl.parent_column_id = p.id);

-- -------------------------------------------------------------------------
-- 11. Связи колонок между таблицами (dc.has_to_group)
-- -------------------------------------------------------------------------
INSERT INTO dc.has_to_group (column_id_a, column_id_b, description,
                             is_deleted, created_at, updated_at, user_id)
SELECT ca.id, cb.id, v.description, false, now(), now(), u.id
FROM (VALUES ('test_tbl_orders', 'client_id', 'test_tbl_clients', 'client_id', 'Заказы -> клиенты'),
             ('test_tbl_clients', 'city_id', 'test_tbl_geo', 'city_id', 'Клиенты -> география'),
             ('test_tbl_payments', 'order_id', 'test_tbl_orders', 'order_id', 'Платежи -> заказы'),
             -- Связи для 10 таблиц из секции 8
             ('test_tbl_products', 'supplier_id', 'test_tbl_suppliers', 'supplier_id', 'Товары -> поставщики'),
             ('test_tbl_inventory', 'product_id', 'test_tbl_products', 'product_id', 'Остатки -> товары'),
             ('test_tbl_shipments', 'order_id', 'test_tbl_orders', 'order_id', 'Отгрузки -> заказы'),
             ('test_tbl_returns', 'order_id', 'test_tbl_orders', 'order_id', 'Возвраты -> заказы'),
             ('test_tbl_invoices', 'payment_id', 'test_tbl_payments', 'payment_id', 'Счета -> платежи'),
             ('test_tbl_refunds', 'payment_id', 'test_tbl_payments', 'payment_id', 'Возвраты платежей -> платежи'),
             ('test_tbl_regions_ext', 'region_id', 'test_tbl_geo', 'region_id', 'Расширенные регионы -> география')
     ) AS v(table_a, column_a, table_b, column_b, description)
         JOIN dc.table_cat ta ON ta."name" = v.table_a
         JOIN dc.column_cat ca ON ca.table_id = ta.id AND ca."name" = v.column_a
         JOIN dc.table_cat tb ON tb."name" = v.table_b
         JOIN dc.column_cat cb ON cb.table_id = tb.id AND cb."name" = v.column_b
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
WHERE NOT EXISTS (SELECT 1
                  FROM dc.has_to_group h
                  WHERE h.column_id_a = ca.id
                    AND h.column_id_b = cb.id);

-- -------------------------------------------------------------------------
-- 12. Доступные расчёты для баз (dc.database_calculation)
-- -------------------------------------------------------------------------
INSERT INTO dc.database_calculation (database_cat_id, calculation_type_id,
                                     created_at, updated_at, is_deleted, user_id)
SELECT d.id, calc.id, now(), now(), false, u.id
FROM (VALUES ('test_db_dwh', 'sum'),
             ('test_db_dwh', 'count'),
             ('test_db_dwh', 'count_distinct'),
             ('test_db_dwh', 'avg'),
             ('test_db_dwh', 'max'),
             ('test_db_dwh', 'min'),
             ('test_db_oltp', 'sum'),
             ('test_db_oltp', 'count'),
             ('test_db_oltp', 'max'),
             ('test_db_oltp', 'min')
     ) AS v(db_name, calc_name)
         JOIN dc.database_cat d ON d."name" = v.db_name
         JOIN dc.calculation_type calc ON calc."name" = v.calc_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
WHERE NOT EXISTS (SELECT 1
                  FROM dc.database_calculation dcalc
                  WHERE dcalc.database_cat_id = d.id
                    AND dcalc.calculation_type_id = calc.id);

-- -------------------------------------------------------------------------
-- 13. Доступные расчёты для колонок (dc.following_calculation)
-- -------------------------------------------------------------------------
INSERT INTO dc.following_calculation (column_cat_id, calculation_type_id,
                                      created_at, updated_at, is_deleted, user_id)
SELECT c.id, calc.id, now(), now(), false, u.id
FROM (VALUES ('test_tbl_orders', 'amount', 'sum'),
             ('test_tbl_orders', 'amount', 'avg'),
             ('test_tbl_orders', 'amount', 'max'),
             ('test_tbl_payments', 'amount', 'sum'),
             -- Расчёты для колонок 10 таблиц из секции 8
             ('test_tbl_inventory', 'qty', 'sum'),
             ('test_tbl_inventory', 'qty', 'avg'),
             ('test_tbl_shipments', 'amount', 'sum'),
             ('test_tbl_returns', 'amount', 'sum'),
             ('test_tbl_invoices', 'amount', 'sum'),
             ('test_tbl_refunds', 'amount', 'sum')
     ) AS v(table_name, column_name, calc_name)
         JOIN dc.table_cat t ON t."name" = v.table_name
         JOIN dc.column_cat c ON c.table_id = t.id AND c."name" = v.column_name
         JOIN dc.calculation_type calc ON calc."name" = v.calc_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
WHERE NOT EXISTS (SELECT 1
                  FROM dc.following_calculation fc
                  WHERE fc.column_cat_id = c.id
                    AND fc.calculation_type_id = calc.id);

-- -------------------------------------------------------------------------
-- 14. Привязка ролей доменов к пользователям
-- -------------------------------------------------------------------------
-- В dc.domain_roles новые строки не создаются — используются только
-- can_read/can_write/can_grant, засеянные миграцией 000000039.
--
-- dc.user_domain_roles хранит домен прямо в строке (domain_id), поэтому
-- роль всегда выдаётся на конкретный домен и отдельной связки
-- domain <-> role (бывшая dc.domains_domain_roles) больше не нужно.
-- test_domain_finance намеренно НЕ получает здесь ни одной строки:
-- единственный путь к её данным — через test_tbl_payments (см. секцию 15).
-- updated_by_id (миграции 000000044) — кто выдал/изменил роль. Роли выдаёт
-- админ, поэтому здесь всегда test_user_admin, а не сам получатель роли.
INSERT INTO dc.user_domain_roles (user_id, domain_roles_id, domain_id, created_at, updated_at, is_deleted, updated_by_id)
SELECT u.id, r.id, dom.id, now(), now(), false, admin.id
FROM (VALUES ('test_user_admin', 'can_read', 'test_domain_sales'),
             ('test_user_admin', 'can_write', 'test_domain_sales'),
             ('test_user_reader', 'can_read', 'test_domain_sales'),
             ('test_user_writer', 'can_write', 'test_domain_sales')
     ) AS v(user_name, role_name, domain_name)
         JOIN dc."user" u ON u."name" = v.user_name
         JOIN dc.domain_roles r ON r."name" = v.role_name
         JOIN dc.domain_cat dom ON dom.domain_name = v.domain_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') admin
WHERE NOT EXISTS (SELECT 1
                  FROM dc.user_domain_roles udr
                  WHERE udr.user_id = u.id
                    AND udr.domain_roles_id = r.id
                    AND udr.domain_id = dom.id);

-- -------------------------------------------------------------------------
-- 15. Привязка ролей таблиц к пользователям
-- -------------------------------------------------------------------------
-- В dc.table_roles новые строки не создаются — используются только
-- can_read/can_write/can_grant, засеянные миграцией 000000040.
--
-- dc.user_table_roles хранит таблицу прямо в строке (table_id), поэтому
-- отдельная связка table <-> role (бывшая dc.tables_table_roles) не нужна.
-- test_tbl_payments — единственная точка входа в test_domain_finance
-- (см. комментарий в секции 14).
-- test_user_no_roles намеренно не получает ни одной роли.
-- test_user_writer намеренно не получает роль на test_tbl_payments —
-- у него нет ни доменного, ни табличного пути к test_domain_finance.
-- updated_by_id (миграции 000000043) — кто выдал/изменил роль, см. комментарий в секции 14.
-- 10 таблиц из секции 8 (products..regions_ext) разбросаны по трём пользователям
-- вперемешку can_read/can_write, чтобы проверять табличные права независимо
-- от доменных (test_user_writer имеет доменный can_write только на sales,
-- но здесь получает и табличные роли на finance-таблицы, и наоборот).
INSERT INTO dc.user_table_roles (user_id, table_roles_id, table_id, created_at, updated_at, is_deleted, updated_by_id)
SELECT u.id, r.id, t.id, now(), now(), false, admin.id
FROM (VALUES ('test_user_admin', 'can_read', 'test_tbl_payments'),
             ('test_user_reader', 'can_read', 'test_tbl_payments'),
             ('test_user_reader', 'can_read', 'test_tbl_products'),
             ('test_user_writer', 'can_write', 'test_tbl_suppliers'),
             ('test_user_admin', 'can_read', 'test_tbl_inventory'),
             ('test_user_reader', 'can_read', 'test_tbl_shipments'),
             ('test_user_writer', 'can_write', 'test_tbl_returns'),
             ('test_user_admin', 'can_write', 'test_tbl_invoices'),
             ('test_user_reader', 'can_read', 'test_tbl_refunds'),
             ('test_user_writer', 'can_write', 'test_tbl_employees'),
             ('test_user_admin', 'can_read', 'test_tbl_campaigns'),
             ('test_user_reader', 'can_write', 'test_tbl_regions_ext')
     ) AS v(user_name, role_name, table_name)
         JOIN dc."user" u ON u."name" = v.user_name
         JOIN dc.table_roles r ON r."name" = v.role_name
         JOIN dc.table_cat t ON t."name" = v.table_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') admin
WHERE NOT EXISTS (SELECT 1
                  FROM dc.user_table_roles utr
                  WHERE utr.user_id = u.id
                    AND utr.table_roles_id = r.id
                    AND utr.table_id = t.id);

-- -------------------------------------------------------------------------
-- 16. Проверка инвариантов набора (до COMMIT)
-- -------------------------------------------------------------------------
-- Блок ничего не пишет, только считает и падает с понятным сообщением, если
-- набор не соответствует требованиям. Он внутри той же транзакции, поэтому
-- при нарушении откатывается весь скрипт — полусобранных данных в базе
-- не остаётся.
DO
$$
    DECLARE
        v_tables            bigint;
        v_tables_no_columns bigint;
        v_min_columns       bigint;
        v_min_visible       bigint;
        v_ui_without_alias  bigint;
        v_type_mismatch     bigint;
    BEGIN
        -- 1. Не меньше 300 тестовых таблиц
        SELECT count(*) INTO v_tables FROM dc.table_cat WHERE "name" LIKE 'test\_%';
        IF v_tables < 300 THEN
            RAISE EXCEPTION 'Тестовых таблиц %, требуется не меньше 300', v_tables;
        END IF;

        -- 2. Таблиц без колонок быть не должно: они не попадут в GROUP BY ниже
        --    и молча пройдут проверку минимума
        SELECT count(*)
        INTO v_tables_no_columns
        FROM dc.table_cat t
        WHERE t."name" LIKE 'test\_%'
          AND NOT EXISTS (SELECT 1 FROM dc.column_cat c WHERE c.table_id = t.id);
        IF v_tables_no_columns > 0 THEN
            RAISE EXCEPTION 'Тестовых таблиц без единой колонки: %', v_tables_no_columns;
        END IF;

        -- 3. Не меньше 30 колонок и не меньше 10 видимых в каждой таблице
        SELECT min(q.cnt), min(q.visible)
        INTO v_min_columns, v_min_visible
        FROM (SELECT count(*)                                   AS cnt,
                     count(*) FILTER (WHERE c.show_in_ui)       AS visible
              FROM dc.column_cat c
                       JOIN dc.table_cat t ON t.id = c.table_id
              WHERE t."name" LIKE 'test\_%'
              GROUP BY c.table_id) q;

        IF v_min_columns < 30 THEN
            RAISE EXCEPTION 'Минимум колонок в тестовой таблице — %, требуется не меньше 30', v_min_columns;
        END IF;

        IF v_min_visible < 10 THEN
            RAISE EXCEPTION 'Минимум видимых (show_in_ui) колонок в тестовой таблице — %, требуется не меньше 10', v_min_visible;
        END IF;

        -- 4. show_in_ui = true требует alias_id (ColumnCatGuard)
        SELECT count(*)
        INTO v_ui_without_alias
        FROM dc.column_cat c
                 JOIN dc.table_cat t ON t.id = c.table_id
        WHERE t."name" LIKE 'test\_%'
          AND c.show_in_ui
          AND c.alias_id IS NULL;
        IF v_ui_without_alias > 0 THEN
            RAISE EXCEPTION 'Колонок с show_in_ui = true и пустым alias_id: %', v_ui_without_alias;
        END IF;

        -- 5. Тип колонки обязан совпадать с типом её алиаса (ColumnCatGuard)
        SELECT count(*)
        INTO v_type_mismatch
        FROM dc.column_cat c
                 JOIN dc.table_cat t ON t.id = c.table_id
                 JOIN dc.alias a ON a.id = c.alias_id
        WHERE t."name" LIKE 'test\_%'
          AND a.column_type_id <> c.column_type_id;
        IF v_type_mismatch > 0 THEN
            RAISE EXCEPTION 'Колонок с типом, не совпадающим с типом алиаса: %', v_type_mismatch;
        END IF;

        RAISE NOTICE 'Тестовый набор: % таблиц, минимум % колонок и % видимых колонок на таблицу',
            v_tables, v_min_columns, v_min_visible;
    END
$$;

COMMIT;

-- =========================================================================
-- Проверка вставленного (запускать отдельно при необходимости):
--
-- SELECT 'user'          AS entity, count(*) FROM dc."user"        WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'host',          count(*) FROM dc.host          WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'database_type', count(*) FROM dc.database_type WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'database_cat',  count(*) FROM dc.database_cat  WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'schema_cat',    count(*) FROM dc.schema_cat    WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'domain_cat',    count(*) FROM dc.domain_cat    WHERE domain_name LIKE 'test\_%'
-- UNION ALL SELECT 'table_cat',     count(*) FROM dc.table_cat     WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'alias',         count(*) FROM dc.alias         WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'column_cat',    count(*) FROM dc.column_cat c
--     WHERE EXISTS (SELECT 1 FROM dc.table_cat t WHERE t.id = c.table_id AND t."name" LIKE 'test\_%')
-- ORDER BY 1;
--
-- Ожидаемо: alias = 32, table_cat = 300 (14 именованных + 286 сгенерированных),
-- column_cat = 9645 = 300 * 32 шаблонных - 7 (столько шаблонных имён уже
-- заняты колонками из секции 9: amount у шести таблиц и qty у inventory)
-- + 52 собственные колонки именованных таблиц.
--
-- Разрез по таблицам (ожидаемо: 300 строк, min_cols = 32, max_cols = 36,
-- min_visible = 20):
--
-- SELECT count(*)     AS tables,
--        min(cnt)     AS min_cols,
--        max(cnt)     AS max_cols,
--        min(visible) AS min_visible
-- FROM (SELECT c.table_id,
--              count(*)                             AS cnt,
--              count(*) FILTER (WHERE c.show_in_ui) AS visible
--       FROM dc.column_cat c
--                JOIN dc.table_cat t ON t.id = c.table_id
--       WHERE t."name" LIKE 'test\_%'
--       GROUP BY c.table_id) q;
-- =========================================================================
