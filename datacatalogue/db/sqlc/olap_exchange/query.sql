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

-- =========================================================
-- dc.column_cat
-- =========================================================

-- name: GetTablesMeta :many
SELECT
    column_cat.id,
    column_cat.table_id,
    column_cat.name,
    column_cat.alias_id,
    column_cat.column_type_id,
    column_cat.description,
    column_cat.calculation_type_id,
    table_cat.id AS table_cat_id,
    table_cat.name AS table_cat_name,
    domain_cat.id AS domain_cat_id,
    domain_cat.domain_name AS domain_cat_domain_name,
    table_type.id AS table_type_id,
    table_type.name AS table_type_name,
    schema_cat.id AS schema_cat_id,
    schema_cat.name AS schema_cat_name,
    database_cat.id AS database_cat_id,
    database_cat.name AS database_cat_name,
    database_type.id AS database_type_id,
    database_type.name AS database_type_name,
    database_type.db_version AS database_type_db_version
FROM dc.column_cat
INNER JOIN dc.table_cat table_cat ON table_cat.id = column_cat.table_id
INNER JOIN dc.domain_cat domain_cat ON domain_cat.id = table_cat.domain_id
INNER JOIN dc.table_type table_type ON table_type.id = table_cat.table_type_id
INNER JOIN dc.schema_cat schema_cat ON schema_cat.id = table_cat.schema_id
INNER JOIN dc.database_cat database_cat ON database_cat.id = schema_cat.database_id
INNER JOIN dc.database_type database_type ON database_type.id = database_cat.database_type_id
WHERE column_cat.is_deleted = false
  AND table_cat.is_deleted = false
  AND domain_cat.is_deleted = false
  AND table_type.is_deleted = false
  AND schema_cat.is_deleted = false
  AND database_cat.is_deleted = false
  AND database_type.is_deleted = false
  AND (column_cat.alias_id in (sqlc.arg(alias_ids_list)::string) and (column_cat.table_id in (select get_table_ids_by_external_user_and_role from dc.get_table_ids_by_external_user_and_role(sqlc.arg(external_id), 'can_read'))))
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_cat.id END DESC;
