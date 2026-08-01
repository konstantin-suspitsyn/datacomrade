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
INSERT INTO dc.alias ("name", description, created_at, updated_at, is_deleted, user_id)
SELECT v.name, v.description, now(), now(), false, u.id
FROM (VALUES ('test_alias_id', 'Идентификатор'),
             ('test_alias_amount', 'Сумма'),
             ('test_alias_dt', 'Дата'),
             ('test_alias_name', 'Наименование'),
             ('test_alias_geo', 'География'),
             ('test_alias_qty', 'Количество (штуки)'),
             -- +24 алиаса для 36 сгенерированных таблиц (секция 8) и общего минимума в 30 алиасов
             ('test_alias_code', 'Код'),
             ('test_alias_status', 'Статус'),
             ('test_alias_flag', 'Признак (флаг)'),
             ('test_alias_email', 'Email'),
             ('test_alias_phone', 'Телефон'),
             ('test_alias_url', 'Ссылка (URL)'),
             ('test_alias_percent', 'Процент'),
             ('test_alias_price', 'Цена за единицу'),
             ('test_alias_weight', 'Вес'),
             ('test_alias_volume', 'Объём'),
             ('test_alias_currency', 'Валюта'),
             ('test_alias_country', 'Страна'),
             ('test_alias_city', 'Город'),
             ('test_alias_category', 'Категория'),
             ('test_alias_type', 'Тип'),
             ('test_alias_score', 'Оценка'),
             ('test_alias_rating', 'Рейтинг'),
             ('test_alias_duration', 'Длительность'),
             ('test_alias_version', 'Версия'),
             ('test_alias_comment', 'Комментарий'),
             ('test_alias_tag', 'Тег'),
             ('test_alias_priority', 'Приоритет'),
             ('test_alias_channel', 'Канал'),
             ('test_alias_source', 'Источник')
     ) AS v(name, description)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

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
             ('test_tbl_regions_ext', 'Тестовый расширенный справочник регионов', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             -- +36 сгенерированных таблиц для общего минимума в 50 таблиц
             ('test_tbl_gen_01', 'Тестовая сгенерированная таблица 1', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_02', 'Тестовая сгенерированная таблица 2', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_03', 'Тестовая сгенерированная таблица 3', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_gen_04', 'Тестовая сгенерированная таблица 4', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_05', 'Тестовая сгенерированная таблица 5', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_06', 'Тестовая сгенерированная таблица 6', 'test_schema_marts', 'dimension', 'test_domain_finance', true),
             ('test_tbl_gen_07', 'Тестовая сгенерированная таблица 7', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_08', 'Тестовая сгенерированная таблица 8', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_09', 'Тестовая сгенерированная таблица 9', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_gen_10', 'Тестовая сгенерированная таблица 10', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_11', 'Тестовая сгенерированная таблица 11', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_12', 'Тестовая сгенерированная таблица 12', 'test_schema_marts', 'dimension', 'test_domain_finance', true),
             ('test_tbl_gen_13', 'Тестовая сгенерированная таблица 13', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_14', 'Тестовая сгенерированная таблица 14', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_15', 'Тестовая сгенерированная таблица 15', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_gen_16', 'Тестовая сгенерированная таблица 16', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_17', 'Тестовая сгенерированная таблица 17', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_18', 'Тестовая сгенерированная таблица 18', 'test_schema_marts', 'dimension', 'test_domain_finance', true),
             ('test_tbl_gen_19', 'Тестовая сгенерированная таблица 19', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_20', 'Тестовая сгенерированная таблица 20', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_21', 'Тестовая сгенерированная таблица 21', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_gen_22', 'Тестовая сгенерированная таблица 22', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_23', 'Тестовая сгенерированная таблица 23', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_24', 'Тестовая сгенерированная таблица 24', 'test_schema_marts', 'dimension', 'test_domain_finance', true),
             ('test_tbl_gen_25', 'Тестовая сгенерированная таблица 25', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_26', 'Тестовая сгенерированная таблица 26', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_27', 'Тестовая сгенерированная таблица 27', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_gen_28', 'Тестовая сгенерированная таблица 28', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_29', 'Тестовая сгенерированная таблица 29', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_30', 'Тестовая сгенерированная таблица 30', 'test_schema_marts', 'dimension', 'test_domain_finance', true),
             ('test_tbl_gen_31', 'Тестовая сгенерированная таблица 31', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_32', 'Тестовая сгенерированная таблица 32', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_33', 'Тестовая сгенерированная таблица 33', 'test_schema_public', 'fact', 'test_domain_finance', false),
             ('test_tbl_gen_34', 'Тестовая сгенерированная таблица 34', 'test_schema_marts', 'dimension', 'test_domain_sales', true),
             ('test_tbl_gen_35', 'Тестовая сгенерированная таблица 35', 'test_schema_public', 'fact', 'test_domain_sales', false),
             ('test_tbl_gen_36', 'Тестовая сгенерированная таблица 36', 'test_schema_marts', 'dimension', 'test_domain_finance', true)
     ) AS v(name, description, schema_name, table_type_name, domain_name, is_get_dict)
         JOIN dc.schema_cat s ON s."name" = v.schema_name
         JOIN dc.table_type tt ON tt."name" = v.table_type_name
         JOIN dc.domain_cat dom ON dom.domain_name = v.domain_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT DO NOTHING;

-- -------------------------------------------------------------------------
-- 9. Колонки
-- -------------------------------------------------------------------------
-- Типы расчёта ('sum', 'count', 'count_distinct', ...) берутся из миграции 000000034.
-- Типы колонок ('bigint', 'string', 'date', 'decimal', ...) берутся из миграции 000000041.
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
             ('test_tbl_geo', 'city_id', 'Идентификатор города', 'test_alias_geo', 'bigint', 'count_distinct', true),
             ('test_tbl_geo', 'region_id', 'Идентификатор региона', 'test_alias_geo', 'bigint', 'count_distinct', true),
             ('test_tbl_geo', 'geo_name', 'Название географии', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_payments', 'payment_id', 'Идентификатор платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_payments', 'order_id', 'Идентификатор заказа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_payments', 'amount', 'Сумма платежа', 'test_alias_amount', 'decimal', 'sum', true),
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
             ('test_tbl_shipments', 'shipment_id', 'Идентификатор отгрузки', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_shipments', 'order_id', 'Идентификатор заказа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_shipments', 'shipment_dt', 'Дата отгрузки', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_shipments', 'amount', 'Сумма отгрузки', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_returns', 'return_id', 'Идентификатор возврата', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_returns', 'order_id', 'Идентификатор заказа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_returns', 'amount', 'Сумма возврата', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_invoices', 'invoice_id', 'Идентификатор счёта', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_invoices', 'payment_id', 'Идентификатор платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_invoices', 'amount', 'Сумма счёта', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_refunds', 'refund_id', 'Идентификатор возврата платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_refunds', 'payment_id', 'Идентификатор платежа', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_refunds', 'amount', 'Сумма возврата платежа', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_employees', 'employee_id', 'Идентификатор сотрудника', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_employees', 'employee_name', 'Имя сотрудника', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_campaigns', 'campaign_id', 'Идентификатор кампании', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_campaigns', 'campaign_name', 'Название кампании', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_regions_ext', 'region_id', 'Идентификатор региона', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_regions_ext', 'region_name', 'Расширенное название региона', 'test_alias_name', 'string', 'count', true),
             -- +216 колонок (по 6: id, name, code, status, amount, dt) для 36 сгенерированных таблиц,
             -- чтобы общий итог по колонкам был не меньше 250
             ('test_tbl_gen_01', 'id', 'Идентификатор (сгенерированная таблица 01)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_01', 'name', 'Наименование (сгенерированная таблица 01)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_01', 'code', 'Код (сгенерированная таблица 01)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_01', 'status', 'Статус (сгенерированная таблица 01)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_01', 'amount', 'Сумма (сгенерированная таблица 01)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_01', 'dt', 'Дата (сгенерированная таблица 01)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_02', 'id', 'Идентификатор (сгенерированная таблица 02)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_02', 'name', 'Наименование (сгенерированная таблица 02)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_02', 'code', 'Код (сгенерированная таблица 02)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_02', 'status', 'Статус (сгенерированная таблица 02)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_02', 'amount', 'Сумма (сгенерированная таблица 02)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_02', 'dt', 'Дата (сгенерированная таблица 02)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_03', 'id', 'Идентификатор (сгенерированная таблица 03)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_03', 'name', 'Наименование (сгенерированная таблица 03)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_03', 'code', 'Код (сгенерированная таблица 03)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_03', 'status', 'Статус (сгенерированная таблица 03)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_03', 'amount', 'Сумма (сгенерированная таблица 03)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_03', 'dt', 'Дата (сгенерированная таблица 03)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_04', 'id', 'Идентификатор (сгенерированная таблица 04)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_04', 'name', 'Наименование (сгенерированная таблица 04)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_04', 'code', 'Код (сгенерированная таблица 04)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_04', 'status', 'Статус (сгенерированная таблица 04)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_04', 'amount', 'Сумма (сгенерированная таблица 04)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_04', 'dt', 'Дата (сгенерированная таблица 04)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_05', 'id', 'Идентификатор (сгенерированная таблица 05)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_05', 'name', 'Наименование (сгенерированная таблица 05)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_05', 'code', 'Код (сгенерированная таблица 05)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_05', 'status', 'Статус (сгенерированная таблица 05)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_05', 'amount', 'Сумма (сгенерированная таблица 05)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_05', 'dt', 'Дата (сгенерированная таблица 05)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_06', 'id', 'Идентификатор (сгенерированная таблица 06)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_06', 'name', 'Наименование (сгенерированная таблица 06)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_06', 'code', 'Код (сгенерированная таблица 06)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_06', 'status', 'Статус (сгенерированная таблица 06)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_06', 'amount', 'Сумма (сгенерированная таблица 06)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_06', 'dt', 'Дата (сгенерированная таблица 06)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_07', 'id', 'Идентификатор (сгенерированная таблица 07)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_07', 'name', 'Наименование (сгенерированная таблица 07)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_07', 'code', 'Код (сгенерированная таблица 07)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_07', 'status', 'Статус (сгенерированная таблица 07)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_07', 'amount', 'Сумма (сгенерированная таблица 07)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_07', 'dt', 'Дата (сгенерированная таблица 07)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_08', 'id', 'Идентификатор (сгенерированная таблица 08)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_08', 'name', 'Наименование (сгенерированная таблица 08)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_08', 'code', 'Код (сгенерированная таблица 08)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_08', 'status', 'Статус (сгенерированная таблица 08)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_08', 'amount', 'Сумма (сгенерированная таблица 08)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_08', 'dt', 'Дата (сгенерированная таблица 08)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_09', 'id', 'Идентификатор (сгенерированная таблица 09)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_09', 'name', 'Наименование (сгенерированная таблица 09)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_09', 'code', 'Код (сгенерированная таблица 09)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_09', 'status', 'Статус (сгенерированная таблица 09)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_09', 'amount', 'Сумма (сгенерированная таблица 09)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_09', 'dt', 'Дата (сгенерированная таблица 09)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_10', 'id', 'Идентификатор (сгенерированная таблица 10)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_10', 'name', 'Наименование (сгенерированная таблица 10)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_10', 'code', 'Код (сгенерированная таблица 10)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_10', 'status', 'Статус (сгенерированная таблица 10)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_10', 'amount', 'Сумма (сгенерированная таблица 10)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_10', 'dt', 'Дата (сгенерированная таблица 10)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_11', 'id', 'Идентификатор (сгенерированная таблица 11)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_11', 'name', 'Наименование (сгенерированная таблица 11)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_11', 'code', 'Код (сгенерированная таблица 11)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_11', 'status', 'Статус (сгенерированная таблица 11)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_11', 'amount', 'Сумма (сгенерированная таблица 11)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_11', 'dt', 'Дата (сгенерированная таблица 11)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_12', 'id', 'Идентификатор (сгенерированная таблица 12)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_12', 'name', 'Наименование (сгенерированная таблица 12)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_12', 'code', 'Код (сгенерированная таблица 12)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_12', 'status', 'Статус (сгенерированная таблица 12)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_12', 'amount', 'Сумма (сгенерированная таблица 12)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_12', 'dt', 'Дата (сгенерированная таблица 12)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_13', 'id', 'Идентификатор (сгенерированная таблица 13)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_13', 'name', 'Наименование (сгенерированная таблица 13)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_13', 'code', 'Код (сгенерированная таблица 13)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_13', 'status', 'Статус (сгенерированная таблица 13)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_13', 'amount', 'Сумма (сгенерированная таблица 13)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_13', 'dt', 'Дата (сгенерированная таблица 13)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_14', 'id', 'Идентификатор (сгенерированная таблица 14)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_14', 'name', 'Наименование (сгенерированная таблица 14)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_14', 'code', 'Код (сгенерированная таблица 14)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_14', 'status', 'Статус (сгенерированная таблица 14)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_14', 'amount', 'Сумма (сгенерированная таблица 14)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_14', 'dt', 'Дата (сгенерированная таблица 14)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_15', 'id', 'Идентификатор (сгенерированная таблица 15)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_15', 'name', 'Наименование (сгенерированная таблица 15)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_15', 'code', 'Код (сгенерированная таблица 15)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_15', 'status', 'Статус (сгенерированная таблица 15)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_15', 'amount', 'Сумма (сгенерированная таблица 15)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_15', 'dt', 'Дата (сгенерированная таблица 15)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_16', 'id', 'Идентификатор (сгенерированная таблица 16)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_16', 'name', 'Наименование (сгенерированная таблица 16)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_16', 'code', 'Код (сгенерированная таблица 16)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_16', 'status', 'Статус (сгенерированная таблица 16)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_16', 'amount', 'Сумма (сгенерированная таблица 16)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_16', 'dt', 'Дата (сгенерированная таблица 16)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_17', 'id', 'Идентификатор (сгенерированная таблица 17)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_17', 'name', 'Наименование (сгенерированная таблица 17)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_17', 'code', 'Код (сгенерированная таблица 17)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_17', 'status', 'Статус (сгенерированная таблица 17)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_17', 'amount', 'Сумма (сгенерированная таблица 17)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_17', 'dt', 'Дата (сгенерированная таблица 17)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_18', 'id', 'Идентификатор (сгенерированная таблица 18)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_18', 'name', 'Наименование (сгенерированная таблица 18)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_18', 'code', 'Код (сгенерированная таблица 18)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_18', 'status', 'Статус (сгенерированная таблица 18)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_18', 'amount', 'Сумма (сгенерированная таблица 18)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_18', 'dt', 'Дата (сгенерированная таблица 18)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_19', 'id', 'Идентификатор (сгенерированная таблица 19)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_19', 'name', 'Наименование (сгенерированная таблица 19)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_19', 'code', 'Код (сгенерированная таблица 19)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_19', 'status', 'Статус (сгенерированная таблица 19)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_19', 'amount', 'Сумма (сгенерированная таблица 19)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_19', 'dt', 'Дата (сгенерированная таблица 19)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_20', 'id', 'Идентификатор (сгенерированная таблица 20)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_20', 'name', 'Наименование (сгенерированная таблица 20)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_20', 'code', 'Код (сгенерированная таблица 20)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_20', 'status', 'Статус (сгенерированная таблица 20)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_20', 'amount', 'Сумма (сгенерированная таблица 20)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_20', 'dt', 'Дата (сгенерированная таблица 20)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_21', 'id', 'Идентификатор (сгенерированная таблица 21)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_21', 'name', 'Наименование (сгенерированная таблица 21)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_21', 'code', 'Код (сгенерированная таблица 21)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_21', 'status', 'Статус (сгенерированная таблица 21)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_21', 'amount', 'Сумма (сгенерированная таблица 21)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_21', 'dt', 'Дата (сгенерированная таблица 21)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_22', 'id', 'Идентификатор (сгенерированная таблица 22)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_22', 'name', 'Наименование (сгенерированная таблица 22)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_22', 'code', 'Код (сгенерированная таблица 22)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_22', 'status', 'Статус (сгенерированная таблица 22)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_22', 'amount', 'Сумма (сгенерированная таблица 22)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_22', 'dt', 'Дата (сгенерированная таблица 22)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_23', 'id', 'Идентификатор (сгенерированная таблица 23)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_23', 'name', 'Наименование (сгенерированная таблица 23)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_23', 'code', 'Код (сгенерированная таблица 23)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_23', 'status', 'Статус (сгенерированная таблица 23)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_23', 'amount', 'Сумма (сгенерированная таблица 23)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_23', 'dt', 'Дата (сгенерированная таблица 23)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_24', 'id', 'Идентификатор (сгенерированная таблица 24)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_24', 'name', 'Наименование (сгенерированная таблица 24)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_24', 'code', 'Код (сгенерированная таблица 24)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_24', 'status', 'Статус (сгенерированная таблица 24)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_24', 'amount', 'Сумма (сгенерированная таблица 24)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_24', 'dt', 'Дата (сгенерированная таблица 24)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_25', 'id', 'Идентификатор (сгенерированная таблица 25)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_25', 'name', 'Наименование (сгенерированная таблица 25)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_25', 'code', 'Код (сгенерированная таблица 25)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_25', 'status', 'Статус (сгенерированная таблица 25)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_25', 'amount', 'Сумма (сгенерированная таблица 25)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_25', 'dt', 'Дата (сгенерированная таблица 25)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_26', 'id', 'Идентификатор (сгенерированная таблица 26)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_26', 'name', 'Наименование (сгенерированная таблица 26)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_26', 'code', 'Код (сгенерированная таблица 26)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_26', 'status', 'Статус (сгенерированная таблица 26)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_26', 'amount', 'Сумма (сгенерированная таблица 26)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_26', 'dt', 'Дата (сгенерированная таблица 26)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_27', 'id', 'Идентификатор (сгенерированная таблица 27)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_27', 'name', 'Наименование (сгенерированная таблица 27)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_27', 'code', 'Код (сгенерированная таблица 27)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_27', 'status', 'Статус (сгенерированная таблица 27)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_27', 'amount', 'Сумма (сгенерированная таблица 27)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_27', 'dt', 'Дата (сгенерированная таблица 27)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_28', 'id', 'Идентификатор (сгенерированная таблица 28)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_28', 'name', 'Наименование (сгенерированная таблица 28)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_28', 'code', 'Код (сгенерированная таблица 28)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_28', 'status', 'Статус (сгенерированная таблица 28)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_28', 'amount', 'Сумма (сгенерированная таблица 28)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_28', 'dt', 'Дата (сгенерированная таблица 28)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_29', 'id', 'Идентификатор (сгенерированная таблица 29)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_29', 'name', 'Наименование (сгенерированная таблица 29)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_29', 'code', 'Код (сгенерированная таблица 29)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_29', 'status', 'Статус (сгенерированная таблица 29)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_29', 'amount', 'Сумма (сгенерированная таблица 29)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_29', 'dt', 'Дата (сгенерированная таблица 29)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_30', 'id', 'Идентификатор (сгенерированная таблица 30)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_30', 'name', 'Наименование (сгенерированная таблица 30)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_30', 'code', 'Код (сгенерированная таблица 30)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_30', 'status', 'Статус (сгенерированная таблица 30)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_30', 'amount', 'Сумма (сгенерированная таблица 30)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_30', 'dt', 'Дата (сгенерированная таблица 30)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_31', 'id', 'Идентификатор (сгенерированная таблица 31)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_31', 'name', 'Наименование (сгенерированная таблица 31)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_31', 'code', 'Код (сгенерированная таблица 31)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_31', 'status', 'Статус (сгенерированная таблица 31)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_31', 'amount', 'Сумма (сгенерированная таблица 31)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_31', 'dt', 'Дата (сгенерированная таблица 31)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_32', 'id', 'Идентификатор (сгенерированная таблица 32)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_32', 'name', 'Наименование (сгенерированная таблица 32)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_32', 'code', 'Код (сгенерированная таблица 32)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_32', 'status', 'Статус (сгенерированная таблица 32)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_32', 'amount', 'Сумма (сгенерированная таблица 32)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_32', 'dt', 'Дата (сгенерированная таблица 32)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_33', 'id', 'Идентификатор (сгенерированная таблица 33)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_33', 'name', 'Наименование (сгенерированная таблица 33)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_33', 'code', 'Код (сгенерированная таблица 33)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_33', 'status', 'Статус (сгенерированная таблица 33)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_33', 'amount', 'Сумма (сгенерированная таблица 33)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_33', 'dt', 'Дата (сгенерированная таблица 33)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_34', 'id', 'Идентификатор (сгенерированная таблица 34)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_34', 'name', 'Наименование (сгенерированная таблица 34)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_34', 'code', 'Код (сгенерированная таблица 34)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_34', 'status', 'Статус (сгенерированная таблица 34)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_34', 'amount', 'Сумма (сгенерированная таблица 34)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_34', 'dt', 'Дата (сгенерированная таблица 34)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_35', 'id', 'Идентификатор (сгенерированная таблица 35)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_35', 'name', 'Наименование (сгенерированная таблица 35)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_35', 'code', 'Код (сгенерированная таблица 35)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_35', 'status', 'Статус (сгенерированная таблица 35)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_35', 'amount', 'Сумма (сгенерированная таблица 35)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_35', 'dt', 'Дата (сгенерированная таблица 35)', 'test_alias_dt', 'date', 'max', true),
             ('test_tbl_gen_36', 'id', 'Идентификатор (сгенерированная таблица 36)', 'test_alias_id', 'bigint', 'count_distinct', true),
             ('test_tbl_gen_36', 'name', 'Наименование (сгенерированная таблица 36)', 'test_alias_name', 'string', 'count', true),
             ('test_tbl_gen_36', 'code', 'Код (сгенерированная таблица 36)', 'test_alias_code', 'string', 'count', true),
             ('test_tbl_gen_36', 'status', 'Статус (сгенерированная таблица 36)', 'test_alias_status', 'string', 'count', true),
             ('test_tbl_gen_36', 'amount', 'Сумма (сгенерированная таблица 36)', 'test_alias_amount', 'decimal', 'sum', true),
             ('test_tbl_gen_36', 'dt', 'Дата (сгенерированная таблица 36)', 'test_alias_dt', 'date', 'max', true)
     ) AS v(table_name, name, description, alias_name, coltype_name, calc_name, show_in_ui)
         JOIN dc.table_cat t ON t."name" = v.table_name
         JOIN dc.alias a ON a."name" = v.alias_name
         JOIN dc.column_type ct ON ct."name" = v.coltype_name
         JOIN dc.calculation_type calc ON calc."name" = v.calc_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
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
-- =========================================================================
