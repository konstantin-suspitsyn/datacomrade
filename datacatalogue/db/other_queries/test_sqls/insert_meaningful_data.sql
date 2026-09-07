-- =========================================================================
-- Осмысленные тестовые данные для dc.*
-- =========================================================================
-- Отдельный от insert_test_data.sql набор: те же соглашения, но данные
-- подобраны как связная предметная область, а не для проверки объёмных
-- инвариантов.
--
-- Соглашения (важны для delete_test_data.sql):
--
--   1. Все именованные сущности создаются с префиксом `test_`.
--      Удаление идёт по LIKE 'test\_%'.
--   2. Сущности без собственного имени (dc.column_cat, dc.group_levels,
--      dc.has_to_group, dc.tables_joins, связки ролей, dc.*_calculation)
--      удаляются через связь с родителем, помеченным префиксом `test_`.
--   3. Скрипт НИЧЕГО не вставляет в справочники из миграций
--      (dc.calculation_type, dc.table_type, dc.column_type,
--      dc.domain_roles, dc.table_roles) и не трогает пользователя
--      'Lomonosov M.' — только читает их по имени.
--   4. Никаких захардкоженных id: все FK разрешаются подзапросами
--      по уникальным именам.
--   5. Скрипт идемпотентен: ON CONFLICT DO NOTHING там, где есть
--      unique-констрейнт, NOT EXISTS — там, где его нет.
--
-- Владелец всех тестовых записей — dc."user".name = 'test_user_admin'.
-- =========================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- 1. Пользователи
-- -------------------------------------------------------------------------
-- test_user_admin   — владелец всех тестовых объектов набора (user_id во
--                     всех секциях ниже), выдаёт роли (updated_by_id).
-- test_user_analyst — читатель: доменный can_read.
-- test_user_steward — куратор данных: доменный can_write + can_grant.
--
-- external_id обязателен и уникален; диапазон 'bbbbbbbb-...' выбран так,
-- чтобы не пересекаться с 'aaaaaaaa-...' из insert_test_data.sql.
INSERT INTO dc."user" ("name", external_id, created_at, updated_at, is_deleted)
VALUES ('test_user_admin', 'bbbbbbbb-0000-4000-8000-000000000001', now(), now(), false),
       ('test_user_analyst', 'bbbbbbbb-0000-4000-8000-000000000002', now(), now(), false),
       ('test_user_steward', 'bbbbbbbb-0000-4000-8000-000000000003', now(), now(), false)
ON CONFLICT ON CONSTRAINT user_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 2. Хосты
-- -------------------------------------------------------------------------
-- Один хост — ClickHouse-кластер аналитического хранилища.
-- host_env / port_env / username_env / password_env уникальны в пределах
-- таблицы; префикс TEST_ANALYTICS_CH_ выбран так, чтобы не пересекаться
-- с TEST_CH_* из insert_test_data.sql.
INSERT INTO dc.host ("name", description, host_env, port_env, username_env, password_env,
                     is_deleted, created_at, updated_at, user_id)
SELECT v.name, v.description, v.host_env, v.port_env, v.username_env, v.password_env,
       false, now(), now(), u.id
FROM (VALUES ('test_host_analytics_ch', 'Кластер ClickHouse аналитического хранилища (площадка Москва-М9)',
              'TEST_ANALYTICS_CH_HOST', 'TEST_ANALYTICS_CH_PORT',
              'TEST_ANALYTICS_CH_USER', 'TEST_ANALYTICS_CH_PASSWORD')
     ) AS v(name, description, host_env, port_env, username_env, password_env)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT host_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 3. Типы БД
-- -------------------------------------------------------------------------
-- Тип под ClickHouse; db_version — актуальная LTS на момент набора.
INSERT INTO dc.database_type ("name", db_version, is_deleted, created_at, updated_at, user_id)
SELECT v.name, v.db_version, false, now(), now(), u.id
FROM (VALUES ('test_dbtype_clickhouse', '24.8')
     ) AS v(name, db_version)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT database_type_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 4. Базы данных
-- -------------------------------------------------------------------------
-- Одна база — аналитическое хранилище на ClickHouse-хосте из секции 2.
INSERT INTO dc.database_cat ("name", host_id, database_type_id, description,
                             is_deleted, created_at, updated_at, user_id)
SELECT v.name, h.id, dt.id, v.description, false, now(), now(), u.id
FROM (VALUES ('test_db_analytics', 'test_host_analytics_ch', 'test_dbtype_clickhouse',
              'Аналитическое хранилище: витрины продаж и справочники')
     ) AS v(name, host_name, dbtype_name, description)
         JOIN dc.host h ON h."name" = v.host_name
         JOIN dc.database_type dt ON dt."name" = v.dbtype_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT database_cat_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 5. Схемы (в dc.schema_cat нет unique-констрейнта -> защита через NOT EXISTS)
-- -------------------------------------------------------------------------
-- Две схемы одной базы:
--   test_schema_dims — справочники / измерения (dimension);
--   test_schema_fact — фактовые таблицы.
-- Имена под маску test_ (иначе delete_test_data.sql их не уберёт);
-- физические имена в целевой БД — 'dims' и 'fact'.
INSERT INTO dc.schema_cat (database_id, "name", is_deleted, created_at, updated_at, user_id)
SELECT d.id, v.name, false, now(), now(), u.id
FROM (VALUES ('test_schema_dims', 'test_db_analytics'),
             ('test_schema_fact', 'test_db_analytics')
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
-- Тестовый бизнес-домен «продажи». Фиксированный домен 'no_domain'
-- (миграция 000000048) здесь не трогаем — он уже есть в базе, таблицы
-- общих словарей будут ссылаться на него напрямую.
INSERT INTO dc.domain_cat (domain_name, is_deleted, created_at, updated_at, user_id)
SELECT v.domain_name, false, now(), now(), u.id
FROM (VALUES ('test_domain_sales')
     ) AS v(domain_name)
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT domain_cat_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 7. Доменные роли пользователей (dc.user_domain_roles)
-- -------------------------------------------------------------------------
-- Роли can_read / can_write / can_grant берутся из миграции 000000039,
-- новых строк в dc.domain_roles не создаётся. Домен хранится прямо в
-- строке связки (domain_id). updated_by_id — кто выдал роль: всегда
-- test_user_admin. В dc.user_domain_roles нет unique-констрейнта ->
-- защита через NOT EXISTS по (user_id, domain_roles_id, domain_id).
--
--   test_user_analyst -> can_read  на test_domain_sales
--   test_user_steward -> can_write на test_domain_sales
--   test_user_steward -> can_grant на test_domain_sales
INSERT INTO dc.user_domain_roles (user_id, domain_roles_id, domain_id,
                                  created_at, updated_at, is_deleted, updated_by_id)
SELECT u.id, r.id, dom.id, now(), now(), false, admin.id
FROM (VALUES ('test_user_analyst', 'can_read', 'test_domain_sales'),
             ('test_user_steward', 'can_write', 'test_domain_sales'),
             ('test_user_steward', 'can_grant', 'test_domain_sales')
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
-- 8. Алиасы (бизнес-имена колонок)
-- -------------------------------------------------------------------------
-- У алиаса есть собственный column_type_id: он задаёт тип для каждой
-- колонки, которая на него ссылается (правило ColumnCatGuard
-- "alias.column_type_id = column_cat.column_type_id"). Поэтому у каждого
-- test_alias_* ровно один тип, и шаблоны колонок ниже берут (alias, тип)
-- парой. Типы читаются из dc.column_type (миграция 000000041).
INSERT INTO dc.alias ("name", description, column_type_id, created_at, updated_at, is_deleted, user_id)
SELECT v.name, v.description, ct.id, now(), now(), false, u.id
FROM (VALUES ('test_alias_date', 'Дата', 'date'),
             ('test_alias_year', 'Год', 'int'),
             ('test_alias_quarter', 'Квартал', 'int'),
             ('test_alias_month', 'Месяц', 'int'),
             ('test_alias_month_name', 'Название месяца', 'string'),
             ('test_alias_week', 'Неделя года', 'int'),
             ('test_alias_week_of_month', 'Неделя месяца', 'int'),
             ('test_alias_year_week', 'Год-неделя', 'int'),
             ('test_alias_year_month', 'Год-месяц', 'int'),
             ('test_alias_day', 'День месяца', 'int'),
             ('test_alias_day_of_week', 'День недели', 'int'),
             ('test_alias_day_name', 'Название дня недели', 'string'),
             ('test_alias_is_weekend', 'Выходной', 'bool'),
             ('test_alias_is_holiday', 'Праздник', 'bool'),
             ('test_alias_holiday_name', 'Название праздника', 'string')
     ) AS v(name, description, coltype_name)
         JOIN dc.column_type ct ON ct."name" = v.coltype_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT alias_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 9. Таблица-словарь: календарь
-- -------------------------------------------------------------------------
-- Одна строка на календарную дату. Лежит в схеме измерений test_schema_dims
-- и в фиксированном домене no_domain (миграция 000000048) — как общий
-- словарь, доступный всем. Тип dimension, is_get_dict = true: в ClickHouse
-- к нему обращаются через dictGet, а не join.
INSERT INTO dc.table_cat ("name", description, schema_id, table_type_id, domain_id,
                          is_deleted, is_get_dict, created_at, updated_at, user_id)
SELECT v.name, v.description, s.id, tt.id, dom.id, false, true, now(), now(), u.id
FROM (VALUES ('test_dim_calendar', 'Календарь: одна строка на дату, атрибуты дня для группировок и фильтров',
              'test_schema_dims', 'dimension', 'no_domain')
     ) AS v(name, description, schema_name, table_type_name, domain_name)
         JOIN dc.schema_cat s ON s."name" = v.schema_name
         JOIN dc.table_type tt ON tt."name" = v.table_type_name
         JOIN dc.domain_cat dom ON dom.domain_name = v.domain_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT table_cat_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 10. Колонки календаря
-- -------------------------------------------------------------------------
-- date_id — суррогатный ключ вида YYYYMMDD, скрыт и без алиаса
-- (show_in_ui = false, alias_id IS NULL — разрешено миграцией 000000046).
-- Остальные колонки видимы, у каждой alias из секции 8; column_type колонки
-- совпадает с типом алиаса.
--
-- calculation_type_id IS NULL у всех колонок календаря: это словарь-измерение,
-- агрегата у его атрибутов нет (nullable с миграции 000000049).
INSERT INTO dc.column_cat (table_id, "name", alias_id, column_type_id, description,
                           calculation_type_id, is_deleted, show_in_ui,
                           created_at, updated_at, user_id)
SELECT t.id, v.name, a.id, ct.id, v.description, NULL, false, v.show_in_ui, now(), now(), u.id
FROM (VALUES ('date_id', 'Суррогатный ключ даты (YYYYMMDD)', NULL, 'bigint', false),
             ('calendar_date', 'Календарная дата', 'test_alias_date', 'date', true),
             ('year_num', 'Год', 'test_alias_year', 'int', true),
             ('quarter_num', 'Квартал (1-4)', 'test_alias_quarter', 'int', true),
             ('month_num', 'Месяц (1-12)', 'test_alias_month', 'int', true),
             ('month_name', 'Название месяца', 'test_alias_month_name', 'string', true),
             ('week_of_year', 'Номер ISO-недели в году (1-53)', 'test_alias_week', 'int', true),
             ('week_of_month', 'Номер недели в месяце (1-5)', 'test_alias_week_of_month', 'int', true),
             ('year_week', 'Год-неделя YYYYWW (ISO)', 'test_alias_year_week', 'int', true),
             ('year_month', 'Год-месяц YYYYMM', 'test_alias_year_month', 'int', true),
             ('day_of_month', 'День месяца (1-31)', 'test_alias_day', 'int', true),
             ('day_of_week', 'День недели (1-7, понедельник = 1)', 'test_alias_day_of_week', 'int', true),
             ('day_name', 'Название дня недели', 'test_alias_day_name', 'string', true),
             ('is_weekend', 'Признак выходного дня', 'test_alias_is_weekend', 'bool', true),
             ('is_holiday', 'Признак праздничного дня', 'test_alias_is_holiday', 'bool', true),
             ('holiday_name', 'Название праздника (если есть)', 'test_alias_holiday_name', 'string', true)
     ) AS v(name, description, alias_name, coltype_name, show_in_ui)
         JOIN dc.table_cat t ON t."name" = 'test_dim_calendar'
         LEFT JOIN dc.alias a ON a."name" = v.alias_name
         JOIN dc.column_type ct ON ct."name" = v.coltype_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT column_cat_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 11. Алиасы витрины продаж
-- -------------------------------------------------------------------------
-- test_alias_date (секция 8) переиспользуется для sale_date, здесь только
-- новые. Тип алиаса = тип колонки, которая на него сошлётся (ColumnCatGuard).
INSERT INTO dc.alias ("name", description, column_type_id, created_at, updated_at, is_deleted, user_id)
SELECT v.name, v.description, ct.id, now(), now(), false, u.id
FROM (VALUES ('test_alias_product', 'Товар', 'bigint'),
             ('test_alias_customer', 'Клиент', 'bigint'),
             ('test_alias_channel', 'Канал продаж', 'string'),
             ('test_alias_payment_method', 'Способ оплаты', 'string'),
             ('test_alias_quantity', 'Количество', 'int'),
             ('test_alias_amount', 'Сумма', 'decimal'),
             ('test_alias_discount', 'Скидка', 'decimal'),
             ('test_alias_net_amount', 'Сумма со скидкой', 'decimal')
     ) AS v(name, description, coltype_name)
         JOIN dc.column_type ct ON ct."name" = v.coltype_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT alias_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 12. Таблица fact: продажи по дням
-- -------------------------------------------------------------------------
-- Схема test_schema_fact, тип fact, домен test_domain_sales. is_get_dict = false:
-- это фактовая таблица, к ней идёт обычный join, а не dictGet.
-- Гранулярность строки: день × товар × клиент × канал × способ оплаты.
INSERT INTO dc.table_cat ("name", description, schema_id, table_type_id, domain_id,
                          is_deleted, is_get_dict, created_at, updated_at, user_id)
SELECT v.name, v.description, s.id, tt.id, dom.id, false, false, now(), now(), u.id
FROM (VALUES ('test_fact_sales', 'Продажи по дням: агрегаты выручки и количества в разрезе товара, клиента и канала',
              'test_schema_fact', 'fact', 'test_domain_sales')
     ) AS v(name, description, schema_name, table_type_name, domain_name)
         JOIN dc.schema_cat s ON s."name" = v.schema_name
         JOIN dc.table_type tt ON tt."name" = v.table_type_name
         JOIN dc.domain_cat dom ON dom.domain_name = v.domain_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT table_cat_name_unique DO NOTHING;

-- -------------------------------------------------------------------------
-- 13. Колонки витрины продаж
-- -------------------------------------------------------------------------
-- sale_id, date_id — скрытые технические ключи без алиаса. date_id ссылается
-- на test_dim_calendar.date_id (связь оформляется в секции has_to_group ниже).
-- Меры (quantity, *_amount) — calc = sum. Атрибуты-измерения (sale_date,
-- sales_channel, payment_method) — calc IS NULL: по ним только группируют
-- (nullable с миграции 000000049).
INSERT INTO dc.column_cat (table_id, "name", alias_id, column_type_id, description,
                           calculation_type_id, is_deleted, show_in_ui,
                           created_at, updated_at, user_id)
SELECT t.id, v.name, a.id, ct.id, v.description, calc.id, false, v.show_in_ui, now(), now(), u.id
FROM (VALUES ('sale_id', 'Суррогатный ключ продажи', NULL, 'bigint', 'count_distinct', false),
             ('date_id', 'Ключ даты (YYYYMMDD), FK на test_dim_calendar', NULL, 'bigint', 'count_distinct', false),
             ('sale_date', 'Дата продажи', 'test_alias_date', 'date', NULL, true),
             ('product_id', 'Идентификатор товара', 'test_alias_product', 'bigint', 'count_distinct', true),
             ('customer_id', 'Идентификатор клиента', 'test_alias_customer', 'bigint', 'count_distinct', true),
             ('sales_channel', 'Канал продаж (retail, online, partner)', 'test_alias_channel', 'string', NULL, true),
             ('payment_method', 'Способ оплаты (card, cash, transfer)', 'test_alias_payment_method', 'string', NULL, true),
             ('quantity', 'Количество проданных единиц', 'test_alias_quantity', 'int', 'sum', true),
             ('gross_amount', 'Выручка до скидки', 'test_alias_amount', 'decimal', 'sum', true),
             ('discount_amount', 'Сумма скидки', 'test_alias_discount', 'decimal', 'sum', true),
             ('net_amount', 'Выручка со скидкой', 'test_alias_net_amount', 'decimal', 'sum', true)
     ) AS v(name, description, alias_name, coltype_name, calc_name, show_in_ui)
         JOIN dc.table_cat t ON t."name" = 'test_fact_sales'
         LEFT JOIN dc.alias a ON a."name" = v.alias_name
         JOIN dc.column_type ct ON ct."name" = v.coltype_name
         LEFT JOIN dc.calculation_type calc ON calc."name" = v.calc_name
         CROSS JOIN (SELECT id FROM dc."user" WHERE "name" = 'test_user_admin') u
ON CONFLICT ON CONSTRAINT column_cat_name_unique DO NOTHING;

COMMIT;
