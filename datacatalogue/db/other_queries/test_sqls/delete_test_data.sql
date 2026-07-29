-- =========================================================================
-- Удаление тестовых данных, созданных insert_test_data.sql
-- =========================================================================
-- Принципы, чтобы не зацепить лишнее:
--
--   1. Удаляются ТОЛЬКО записи с именем по маске 'test\_%' и записи,
--      привязанные к таким родителям (колонки, уровни группировок, связки
--      ролей, доступные расчёты). Обратный слеш обязателен: без него `_`
--      в LIKE — это wildcard и маска поймала бы, например, 'testing'.
--   2. Справочники из миграций не трогаются вообще:
--      dc.calculation_type, dc.table_type, dc.column_type, а также роли
--      can_read / can_write / can_grant в dc.domain_roles и dc.table_roles —
--      у них нет префикса test_. Пользователь 'Lomonosov M.' (миграция 37)
--      тоже остаётся.
--   3. Удаление идёт по именам, а не по id: скрипт не зависит от sequence
--      и не может «промахнуться» по чужой записи с подходящим id.
--   4. Порядок — строго от детей к родителям (все FK объявлены как
--      ON DELETE RESTRICT, каскадов нет).
--   5. Родительские сущности удаляются строго по маске test_, без расширения
--      «всё, что лежит внутри тестовой схемы/базы». Если кто-то положил
--      в тестовую схему нетестовую таблицу, скрипт упадёт на FK — это
--      осознанный выбор: лучше явная ошибка, чем тихое удаление чужих данных.
--
-- Скрипт идемпотентен: повторный запуск просто удалит 0 строк.
-- =========================================================================

BEGIN;

-- -------------------------------------------------------------------------
-- 1. Уровни группировок по тестовым колонкам
-- -------------------------------------------------------------------------
DELETE
FROM dc.group_levels gl
WHERE EXISTS (SELECT 1
              FROM dc.column_cat c
                       JOIN dc.table_cat t ON t.id = c.table_id
              WHERE t."name" LIKE 'test\_%'
                AND c.id IN (gl.column_id, gl.parent_column_id));

-- -------------------------------------------------------------------------
-- 2. Связи колонок по тестовым колонкам
-- -------------------------------------------------------------------------
DELETE
FROM dc.has_to_group h
WHERE EXISTS (SELECT 1
              FROM dc.column_cat c
                       JOIN dc.table_cat t ON t.id = c.table_id
              WHERE t."name" LIKE 'test\_%'
                AND c.id IN (h.column_id_a, h.column_id_b));

-- -------------------------------------------------------------------------
-- 3. Доступные расчёты по тестовым колонкам
-- -------------------------------------------------------------------------
DELETE
FROM dc.following_calculation fc
WHERE EXISTS (SELECT 1
              FROM dc.column_cat c
                       JOIN dc.table_cat t ON t.id = c.table_id
              WHERE c.id = fc.column_cat_id
                AND t."name" LIKE 'test\_%');

-- -------------------------------------------------------------------------
-- 4. Доступные расчёты по тестовым базам
-- -------------------------------------------------------------------------
DELETE
FROM dc.database_calculation dcalc
WHERE EXISTS (SELECT 1
              FROM dc.database_cat d
              WHERE d.id = dcalc.database_cat_id
                AND d."name" LIKE 'test\_%');

-- -------------------------------------------------------------------------
-- 5. Связки ролей
-- -------------------------------------------------------------------------
-- insert_test_data.sql не создаёт строк в dc.domain_roles/dc.table_roles —
-- используются только can_read/can_write/can_grant из миграций 39/40.
-- Отдельных связок domain <-> role и table <-> role (бывшие
-- dc.domains_domain_roles / dc.tables_table_roles) в схеме больше нет:
-- домен/таблица хранятся прямо в dc.user_domain_roles.domain_id и
-- dc.user_table_roles.table_id. Поэтому строки удаляются по тестовому
-- пользователю, домену или таблице, а не по роли (роль здесь всегда чужая,
-- боевая). Условие "OR роль LIKE test\_%" оставлено как защитный запасной
-- вариант — оно ловит связки с ролями test_role_*, которые создавала более
-- ранняя версия скрипта, и не мешает новой схеме (сейчас таких ролей не
-- существует, условие просто не сработает).
DELETE
FROM dc.user_domain_roles udr
WHERE EXISTS (SELECT 1
              FROM dc."user" u
              WHERE u.id = udr.user_id
                AND u."name" LIKE 'test\_%')
   OR EXISTS (SELECT 1
              FROM dc.domain_cat d
              WHERE d.id = udr.domain_id
                AND d.domain_name LIKE 'test\_%')
   OR EXISTS (SELECT 1
              FROM dc.domain_roles r
              WHERE r.id = udr.domain_roles_id
                AND r."name" LIKE 'test\_%');

DELETE
FROM dc.user_table_roles utr
WHERE EXISTS (SELECT 1
              FROM dc."user" u
              WHERE u.id = utr.user_id
                AND u."name" LIKE 'test\_%')
   OR EXISTS (SELECT 1
              FROM dc.table_cat t
              WHERE t.id = utr.table_id
                AND t."name" LIKE 'test\_%')
   OR EXISTS (SELECT 1
              FROM dc.table_roles r
              WHERE r.id = utr.table_roles_id
                AND r."name" LIKE 'test\_%');

-- -------------------------------------------------------------------------
-- 6. Колонки тестовых таблиц
-- -------------------------------------------------------------------------
-- Имена колонок специально без префикса (order_id, amount, ...), поэтому
-- единственный безопасный критерий — принадлежность тестовой таблице.
DELETE
FROM dc.column_cat c
WHERE EXISTS (SELECT 1
              FROM dc.table_cat t
              WHERE t.id = c.table_id
                AND t."name" LIKE 'test\_%');

-- -------------------------------------------------------------------------
-- 7. Таблицы
-- -------------------------------------------------------------------------
DELETE
FROM dc.table_cat
WHERE "name" LIKE 'test\_%';

-- -------------------------------------------------------------------------
-- 8. Схемы (двойное условие: тестовое имя И тестовая база)
-- -------------------------------------------------------------------------
DELETE
FROM dc.schema_cat s
WHERE s."name" LIKE 'test\_%'
  AND EXISTS (SELECT 1
              FROM dc.database_cat d
              WHERE d.id = s.database_id
                AND d."name" LIKE 'test\_%');

-- -------------------------------------------------------------------------
-- 9. Базы, домены, хосты, типы БД, алиасы
-- -------------------------------------------------------------------------
-- dc.column_type здесь нет: скрипт не вставляет в него собственных строк,
-- он только читает реальные типы, засеянные миграцией 000000041.
DELETE
FROM dc.database_cat
WHERE "name" LIKE 'test\_%';

DELETE
FROM dc.domain_cat
WHERE domain_name LIKE 'test\_%';

DELETE
FROM dc.host
WHERE "name" LIKE 'test\_%';

DELETE
FROM dc.database_type
WHERE "name" LIKE 'test\_%';

DELETE
FROM dc.alias
WHERE "name" LIKE 'test\_%';

-- -------------------------------------------------------------------------
-- 10. Legacy-очистка: роли test_role_* от прежней версии скрипта
-- -------------------------------------------------------------------------
-- Текущая версия insert_test_data.sql в dc.domain_roles/dc.table_roles ничего
-- не создаёт (см. секцию 5). Эти два DELETE — не часть активной схемы,
-- а разовая уборка за более ранней версией скрипта, которая заводила
-- test_role_domain_sales_read и подобные строки. can_read/can_write/can_grant
-- из миграций 39/40 не совпадают с маской 'test\_%' и не задеваются.
-- На новых базах, где legacy-строк не было, оба DELETE удаляют 0 строк.
DELETE
FROM dc.domain_roles
WHERE "name" LIKE 'test\_%';

DELETE
FROM dc.table_roles
WHERE "name" LIKE 'test\_%';

-- -------------------------------------------------------------------------
-- 11. Тестовые пользователи ('Lomonosov M.' из миграции 37 не трогаем)
-- -------------------------------------------------------------------------
DELETE
FROM dc."user"
WHERE "name" LIKE 'test\_%';

COMMIT;

-- =========================================================================
-- Предварительный просмотр: что именно будет удалено (безопасно, только SELECT).
-- Запускать ДО удаления, если нужно убедиться, что под маску не попало лишнее.
--
-- SELECT 'user'          AS entity, "name"      AS obj FROM dc."user"        WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'host',          "name"          FROM dc.host          WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'database_type', "name"          FROM dc.database_type WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'database_cat',  "name"          FROM dc.database_cat  WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'schema_cat',    "name"          FROM dc.schema_cat    WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'domain_cat',    domain_name     FROM dc.domain_cat    WHERE domain_name LIKE 'test\_%'
-- UNION ALL SELECT 'table_cat',     "name"          FROM dc.table_cat     WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'alias',         "name"          FROM dc.alias         WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'domain_roles',  "name"          FROM dc.domain_roles  WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'table_roles',   "name"          FROM dc.table_roles   WHERE "name"      LIKE 'test\_%'
-- UNION ALL SELECT 'column_cat',    t."name" || '.' || c."name"
--     FROM dc.column_cat c JOIN dc.table_cat t ON t.id = c.table_id
--     WHERE t."name" LIKE 'test\_%'
-- ORDER BY 1, 2;
-- =========================================================================
