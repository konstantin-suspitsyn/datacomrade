-- Файл сгенерирован программой SG Buddy https://github.com/konstantin-suspitsyn/sg_buddy
-- Правки будут затёрты при следующей генерации: правьте настройки, а не этот файл.

-- =========================================================
-- dc.alias
-- =========================================================

-- name: GetOlapAliases :many
SELECT
    alias.id,
    alias.name,
    alias.description,
    alias.is_deleted,
    alias.column_type_id,
    column_type.name AS column_type_name,
    column_type.description AS column_type_description
FROM dc.alias
INNER JOIN dc.column_type column_type ON column_type.id = alias.column_type_id
WHERE alias.is_deleted = false
  AND column_type.is_deleted = false
  AND (alias.id in (
  select column_cat.alias_id from dc.column_cat 
  where 
  column_cat.table_id in (select get_table_ids_by_external_user_and_role from dc.get_table_ids_by_external_user_and_role(sqlc.arg(external_id), 'can_read'))
  and column_cat.is_deleted = false))
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN alias.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN alias.id END DESC;
