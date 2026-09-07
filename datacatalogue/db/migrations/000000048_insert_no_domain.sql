-- +goose up
-- Фиксированный домен-«заглушка». Существует всегда, наравне со справочными
-- строками can_read/can_write/can_grant (миграции 39/40) и пользователем
-- 'Lomonosov M.' (миграция 37).
--
-- Назначение: таблицы, не относящиеся ни к одной бизнес-области (общие
-- словари, календари, географические справочники), кладутся в этот домен.
-- Доступ к нему открыт всем пользователям в обход dc.user_domain_roles —
-- см. documentation/.../01_architecture/Разграничение доступа.md.
--
-- Владелец строки — 'Lomonosov M.', как у прочих засеянных миграциями данных.
INSERT INTO dc.domain_cat (domain_name, is_deleted, created_at, updated_at, user_id)
SELECT 'no_domain', false, now(), now(), u.id
FROM (SELECT id FROM dc."user" WHERE "name" = 'Lomonosov M.') u;

-- +goose down
DELETE FROM dc.domain_cat
WHERE domain_name = 'no_domain';
