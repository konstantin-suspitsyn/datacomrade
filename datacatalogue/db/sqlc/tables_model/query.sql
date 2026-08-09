-- Файл сгенерирован программой SG Buddy https://github.com/konstantin-suspitsyn/sg_buddy
-- Правки будут затёрты при следующей генерации: правьте настройки, а не этот файл.

-- =========================================================
-- dc.alias
-- =========================================================

-- name: CreateAlias :one
INSERT INTO dc.alias (
    name,
    description,
    created_at,
    updated_at,
    is_deleted,
    user_id
) VALUES (
    @name,
    @description,
    now(),
    now(),
    false,
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetAliasById :one
SELECT
    alias.id,
    alias.name,
    alias.description,
    alias.created_at,
    alias.updated_at,
    alias.is_deleted,
    alias.user_id
FROM dc.alias
WHERE alias.id = @id
  AND alias.is_deleted = false
LIMIT 1;

-- name: GetAliasesDeleted :many
SELECT
    alias.id,
    alias.name,
    alias.description,
    alias.created_at,
    alias.updated_at,
    alias.is_deleted,
    alias.user_id
FROM dc.alias
WHERE alias.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN alias.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN alias.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetAliasesDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.alias
WHERE alias.is_deleted = true;

-- name: GetAliases :many
SELECT
    alias.id,
    alias.name,
    alias.description,
    alias.created_at,
    alias.updated_at,
    alias.is_deleted,
    alias.user_id
FROM dc.alias
WHERE alias.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN alias.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN alias.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetAliases :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.alias
WHERE alias.is_deleted = false;

-- name: UpdateAliasById :exec
UPDATE dc.alias
SET name = @name,
    description = @description,
    updated_at = now(),
    is_deleted = @is_deleted,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE alias.id = @id;

-- name: DeleteAliasById :exec
UPDATE dc.alias
SET updated_at = now(),
    is_deleted = true,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE alias.id = @id;

-- name: UndeleteAliasById :exec
UPDATE dc.alias
SET updated_at = now(),
    is_deleted = false,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE alias.id = @id;

-- =========================================================
-- dc."user"
-- =========================================================

-- name: CreateUser :one
INSERT INTO dc."user" (
    name,
    created_at,
    updated_at,
    is_deleted,
    external_id
) VALUES (
    @name,
    now(),
    now(),
    false,
    @external_id
)
RETURNING *;

-- name: GetUserByExternalId :one
SELECT
    "user".id,
    "user".name,
    "user".created_at,
    "user".updated_at,
    "user".is_deleted,
    "user".external_id
FROM dc."user"
WHERE "user".is_deleted = false
  AND "user".external_id = @external_id
LIMIT 1;

-- name: GetUsers :many
SELECT
    "user".id,
    "user".name,
    "user".created_at,
    "user".updated_at,
    "user".is_deleted,
    "user".external_id
FROM dc."user"
WHERE "user".is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN "user".id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN "user".id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetUsers :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc."user"
WHERE "user".is_deleted = false;

-- =========================================================
-- dc.host
-- =========================================================

-- name: CreateHost :one
INSERT INTO dc.host (
    name,
    description,
    host_env,
    port_env,
    username_env,
    password_env,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @name,
    @description,
    @host_env,
    @port_env,
    @username_env,
    @password_env,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetHostById :one
SELECT
    host.id,
    host.name,
    host.description,
    host.host_env,
    host.port_env,
    host.username_env,
    host.password_env,
    host.is_deleted,
    host.created_at,
    host.updated_at,
    host.user_id
FROM dc.host
WHERE host.id = @id
LIMIT 1;

-- name: GetHosts :many
SELECT
    host.id,
    host.name,
    host.description,
    host.host_env,
    host.port_env,
    host.username_env,
    host.password_env,
    host.is_deleted,
    host.created_at,
    host.updated_at,
    host.user_id
FROM dc.host
WHERE host.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN host.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN host.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetHosts :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.host
WHERE host.is_deleted = false;

-- name: GetHostsSearchName :many
SELECT
    host.id,
    host.name,
    host.description,
    host.host_env,
    host.port_env,
    host.username_env,
    host.password_env,
    host.is_deleted,
    host.created_at,
    host.updated_at,
    host.user_id
FROM dc.host
WHERE host.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN host.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN host.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetHostsSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.host
WHERE host.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: GetHostDeleted :many
SELECT
    host.id,
    host.name,
    host.description,
    host.host_env,
    host.port_env,
    host.username_env,
    host.password_env,
    host.is_deleted,
    host.created_at,
    host.updated_at,
    host.user_id
FROM dc.host
WHERE host.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN host.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN host.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetHostDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.host
WHERE host.is_deleted = true;

-- name: UpdateHostById :exec
UPDATE dc.host
SET id = @id,
    name = @name,
    description = @description,
    host_env = @host_env,
    port_env = @port_env,
    username_env = @username_env,
    password_env = @password_env,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE host.id = @id;

-- name: DeleteHostById :exec
UPDATE dc.host
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE host.id = @id;

-- name: UndeleteHostById :exec
UPDATE dc.host
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE host.id = @id;

-- =========================================================
-- dc.database_type
-- =========================================================

-- name: CreateDatabaseType :one
INSERT INTO dc.database_type (
    name,
    db_version,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @name,
    @db_version,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetDatabaseTypeById :one
SELECT
    database_type.id,
    database_type.name,
    database_type.db_version,
    database_type.is_deleted,
    database_type.created_at,
    database_type.updated_at,
    database_type.user_id
FROM dc.database_type
WHERE database_type.id = @id
  AND database_type.is_deleted = false
LIMIT 1;

-- name: GetDatabaseTypesSearchName :many
SELECT
    database_type.id,
    database_type.name,
    database_type.db_version,
    database_type.is_deleted,
    database_type.created_at,
    database_type.updated_at,
    database_type.user_id
FROM dc.database_type
WHERE database_type.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseTypesSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_type
WHERE database_type.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: GetDatabaseTypes :many
SELECT
    database_type.id,
    database_type.name,
    database_type.db_version,
    database_type.is_deleted,
    database_type.created_at,
    database_type.updated_at,
    database_type.user_id
FROM dc.database_type
WHERE database_type.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseTypes :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_type
WHERE database_type.is_deleted = false;

-- name: GetDatabaseTypesDeleted :many
SELECT
    database_type.id,
    database_type.name,
    database_type.db_version,
    database_type.is_deleted,
    database_type.created_at,
    database_type.updated_at,
    database_type.user_id
FROM dc.database_type
WHERE database_type.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseTypesDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_type
WHERE database_type.is_deleted = true;

-- name: UpdateDatabaseTypeById :exec
UPDATE dc.database_type
SET name = @name,
    db_version = @db_version,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_type.id = @id;

-- name: DeleteDatabaseTypeById :exec
UPDATE dc.database_type
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_type.id = @id;

-- name: UndeleteDatabaseTypeById :exec
UPDATE dc.database_type
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_type.id = @id;

-- =========================================================
-- dc.database_cat
-- =========================================================

-- name: CreateDatabaseCat :one
INSERT INTO dc.database_cat (
    name,
    host_id,
    database_type_id,
    description,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @name,
    @host_id,
    @database_type_id,
    @description,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetDatabaseCatById :one
SELECT
    database_cat.id,
    database_cat.name,
    database_cat.host_id,
    database_cat.database_type_id,
    database_cat.description,
    database_cat.is_deleted,
    database_cat.created_at,
    database_cat.updated_at,
    database_cat.user_id
FROM dc.database_cat
WHERE database_cat.id = @id
LIMIT 1;

-- name: GetDatabaseCats :many
SELECT
    database_cat.id,
    database_cat.name,
    database_cat.host_id,
    database_cat.database_type_id,
    database_cat.description,
    database_cat.is_deleted,
    database_cat.created_at,
    database_cat.updated_at,
    database_cat.user_id
FROM dc.database_cat
WHERE database_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCats :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_cat
WHERE database_cat.is_deleted = false;

-- name: GetDatabaseCatsDeleted :many
SELECT
    database_cat.id,
    database_cat.name,
    database_cat.host_id,
    database_cat.database_type_id,
    database_cat.description,
    database_cat.is_deleted,
    database_cat.created_at,
    database_cat.updated_at,
    database_cat.user_id
FROM dc.database_cat
WHERE database_cat.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCatsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_cat
WHERE database_cat.is_deleted = true;

-- name: GetDatabaseCatsByHostIdDatabaseTypeId :many
SELECT
    database_cat.id,
    database_cat.name,
    database_cat.host_id,
    database_cat.database_type_id,
    database_cat.description,
    database_cat.is_deleted,
    database_cat.created_at,
    database_cat.updated_at,
    database_cat.user_id
FROM dc.database_cat
WHERE (sqlc.narg('host_id')::bigint IS NULL OR database_cat.host_id = sqlc.narg('host_id'))
  AND (sqlc.narg('database_type_id')::bigint IS NULL OR database_cat.database_type_id = sqlc.narg('database_type_id'))
  AND database_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' AND @order_by::text = 'name' THEN database_cat.name END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' AND @order_by::text = 'name' THEN database_cat.name END DESC,
    CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCatsByHostIdDatabaseTypeId :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_cat
WHERE (sqlc.narg('host_id')::bigint IS NULL OR database_cat.host_id = sqlc.narg('host_id'))
  AND (sqlc.narg('database_type_id')::bigint IS NULL OR database_cat.database_type_id = sqlc.narg('database_type_id'))
  AND database_cat.is_deleted = false;

-- name: GetDatabaseCatsSearchName :many
SELECT
    database_cat.id,
    database_cat.name,
    database_cat.host_id,
    database_cat.database_type_id,
    database_cat.description,
    database_cat.is_deleted,
    database_cat.created_at,
    database_cat.updated_at,
    database_cat.user_id
FROM dc.database_cat
WHERE database_cat.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCatsSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_cat
WHERE database_cat.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: UpdateDatabaseCatById :exec
UPDATE dc.database_cat
SET name = @name,
    host_id = @host_id,
    database_type_id = @database_type_id,
    description = @description,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_cat.id = @id;

-- name: DeleteDatabaseCatById :exec
UPDATE dc.database_cat
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_cat.id = @id;

-- name: UndeleteDatabaseCatById :exec
UPDATE dc.database_cat
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_cat.id = @id;

-- =========================================================
-- dc.column_cat
-- =========================================================

-- name: CreateColumnCat :one
INSERT INTO dc.column_cat (
    table_id,
    name,
    alias_id,
    column_type_id,
    description,
    calculation_type_id,
    is_deleted,
    show_in_ui,
    created_at,
    updated_at,
    user_id
) VALUES (
    @table_id,
    @name,
    @alias_id,
    @column_type_id,
    @description,
    @calculation_type_id,
    false,
    @show_in_ui,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetColumnCatById :one
SELECT
    column_cat.id,
    column_cat.table_id,
    column_cat.name,
    column_cat.alias_id,
    column_cat.column_type_id,
    column_cat.description,
    column_cat.calculation_type_id,
    column_cat.is_deleted,
    column_cat.show_in_ui,
    column_cat.created_at,
    column_cat.updated_at,
    column_cat.user_id
FROM dc.column_cat
WHERE column_cat.id = @id
LIMIT 1;

-- name: GetColumnCats :many
SELECT
    column_cat.id,
    column_cat.table_id,
    column_cat.name,
    column_cat.alias_id,
    column_cat.column_type_id,
    column_cat.description,
    column_cat.calculation_type_id,
    column_cat.is_deleted,
    column_cat.show_in_ui,
    column_cat.created_at,
    column_cat.updated_at,
    column_cat.user_id
FROM dc.column_cat
WHERE column_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetColumnCats :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.column_cat
WHERE column_cat.is_deleted = false;

-- name: GetColumnCatsDeleted :many
SELECT
    column_cat.id,
    column_cat.table_id,
    column_cat.name,
    column_cat.alias_id,
    column_cat.column_type_id,
    column_cat.description,
    column_cat.calculation_type_id,
    column_cat.is_deleted,
    column_cat.show_in_ui,
    column_cat.created_at,
    column_cat.updated_at,
    column_cat.user_id
FROM dc.column_cat
WHERE column_cat.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetColumnCatsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.column_cat
WHERE column_cat.is_deleted = true;

-- name: GetColumnCatsOptionalFilters :many
SELECT
    column_cat.id,
    column_cat.table_id,
    column_cat.name,
    column_cat.alias_id,
    column_cat.column_type_id,
    column_cat.description,
    column_cat.calculation_type_id,
    column_cat.is_deleted,
    column_cat.show_in_ui,
    column_cat.created_at,
    column_cat.updated_at,
    column_cat.user_id
FROM dc.column_cat
WHERE (sqlc.narg('table_id')::bigint IS NULL OR column_cat.table_id = sqlc.narg('table_id'))
  AND (sqlc.narg('alias_id')::bigint IS NULL OR column_cat.alias_id = sqlc.narg('alias_id'))
  AND (sqlc.narg('column_type_id')::bigint IS NULL OR column_cat.column_type_id = sqlc.narg('column_type_id'))
  AND (sqlc.narg('calculation_type_id')::bigint IS NULL OR column_cat.calculation_type_id = sqlc.narg('calculation_type_id'))
  AND column_cat.is_deleted = false
  AND (sqlc.narg('show_in_ui')::boolean IS NULL OR column_cat.show_in_ui = sqlc.narg('show_in_ui'))
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetColumnCatsOptionalFilters :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.column_cat
WHERE (sqlc.narg('table_id')::bigint IS NULL OR column_cat.table_id = sqlc.narg('table_id'))
  AND (sqlc.narg('alias_id')::bigint IS NULL OR column_cat.alias_id = sqlc.narg('alias_id'))
  AND (sqlc.narg('column_type_id')::bigint IS NULL OR column_cat.column_type_id = sqlc.narg('column_type_id'))
  AND (sqlc.narg('calculation_type_id')::bigint IS NULL OR column_cat.calculation_type_id = sqlc.narg('calculation_type_id'))
  AND column_cat.is_deleted = false
  AND (sqlc.narg('show_in_ui')::boolean IS NULL OR column_cat.show_in_ui = sqlc.narg('show_in_ui'));

-- name: GetColumnCatsSearchName :many
SELECT
    column_cat.id,
    column_cat.table_id,
    column_cat.name,
    column_cat.alias_id,
    column_cat.column_type_id,
    column_cat.description,
    column_cat.calculation_type_id,
    column_cat.is_deleted,
    column_cat.show_in_ui,
    column_cat.created_at,
    column_cat.updated_at,
    column_cat.user_id
FROM dc.column_cat
WHERE (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetColumnCatsSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.column_cat
WHERE (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: UpdateColumnCatById :exec
UPDATE dc.column_cat
SET table_id = @table_id,
    name = @name,
    alias_id = @alias_id,
    column_type_id = @column_type_id,
    description = @description,
    calculation_type_id = @calculation_type_id,
    is_deleted = false,
    show_in_ui = @show_in_ui,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE column_cat.id = @id;

-- name: DeleteColumnCatById :exec
UPDATE dc.column_cat
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE column_cat.id = @id;

-- name: UndeleteColumnCatById :exec
UPDATE dc.column_cat
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE column_cat.id = @id;

-- =========================================================
-- dc.column_type
-- =========================================================

-- name: CreateColumnType :one
INSERT INTO dc.column_type (
    name,
    description,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @name,
    @description,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetColumnTypeById :one
SELECT
    column_type.id,
    column_type.name,
    column_type.description,
    column_type.is_deleted,
    column_type.created_at,
    column_type.updated_at,
    column_type.user_id
FROM dc.column_type
WHERE column_type.id = @id
LIMIT 1;

-- name: GetColumnTypes :many
SELECT
    column_type.id,
    column_type.name,
    column_type.description,
    column_type.is_deleted,
    column_type.created_at,
    column_type.updated_at,
    column_type.user_id
FROM dc.column_type
WHERE column_type.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetColumnTypes :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.column_type
WHERE column_type.is_deleted = false;

-- name: GetColumnTypesSearchName :many
SELECT
    column_type.id,
    column_type.name,
    column_type.description,
    column_type.is_deleted,
    column_type.created_at,
    column_type.updated_at,
    column_type.user_id
FROM dc.column_type
WHERE column_type.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetColumnTypesSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.column_type
WHERE column_type.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: GetColumnTypesDeleted :many
SELECT
    column_type.id,
    column_type.name,
    column_type.description,
    column_type.is_deleted,
    column_type.created_at,
    column_type.updated_at,
    column_type.user_id
FROM dc.column_type
WHERE column_type.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN column_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN column_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetColumnTypesDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.column_type
WHERE column_type.is_deleted = true;

-- name: UpdateColumnTypeById :exec
UPDATE dc.column_type
SET name = @name,
    description = @description,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE column_type.id = @id;

-- name: DeleteColumnTypeById :exec
UPDATE dc.column_type
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE column_type.id = @id;

-- name: UndeleteColumnTypeById :exec
UPDATE dc.column_type
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE column_type.id = @id;

-- =========================================================
-- dc.database_calculation
-- =========================================================

-- name: CreateDatabaseCalculation :one
INSERT INTO dc.database_calculation (
    database_cat_id,
    calculation_type_id,
    created_at,
    updated_at,
    is_deleted,
    user_id
) VALUES (
    @database_cat_id,
    @calculation_type_id,
    now(),
    now(),
    false,
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetDatabaseCalculationsSearchName :many
SELECT
    database_calculation.id,
    database_calculation.database_cat_id,
    database_calculation.calculation_type_id,
    database_calculation.created_at,
    database_calculation.updated_at,
    database_calculation.is_deleted,
    database_calculation.user_id
FROM dc.database_calculation
WHERE database_calculation.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_calculation.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_calculation.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCalculationsSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_calculation
WHERE database_calculation.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: GetDatabaseCalculationById :one
SELECT
    database_calculation.id,
    database_calculation.database_cat_id,
    database_calculation.calculation_type_id,
    database_calculation.created_at,
    database_calculation.updated_at,
    database_calculation.is_deleted,
    database_calculation.user_id
FROM dc.database_calculation
WHERE database_calculation.id = @id
LIMIT 1;

-- name: GetDatabaseCalculationsFiltered :many
SELECT
    database_calculation.id,
    database_calculation.database_cat_id,
    database_calculation.calculation_type_id,
    database_calculation.created_at,
    database_calculation.updated_at,
    database_calculation.is_deleted,
    database_calculation.user_id
FROM dc.database_calculation
WHERE (sqlc.narg('database_cat_id')::bigint IS NULL OR database_calculation.database_cat_id = sqlc.narg('database_cat_id'))
  AND (sqlc.narg('calculation_type_id')::bigint IS NULL OR database_calculation.calculation_type_id = sqlc.narg('calculation_type_id'))
  AND database_calculation.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_calculation.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_calculation.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCalculationsFiltered :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_calculation
WHERE (sqlc.narg('database_cat_id')::bigint IS NULL OR database_calculation.database_cat_id = sqlc.narg('database_cat_id'))
  AND (sqlc.narg('calculation_type_id')::bigint IS NULL OR database_calculation.calculation_type_id = sqlc.narg('calculation_type_id'))
  AND database_calculation.is_deleted = false;

-- name: GetDatabaseCalculationsDeleted :many
SELECT
    database_calculation.id,
    database_calculation.database_cat_id,
    database_calculation.calculation_type_id,
    database_calculation.created_at,
    database_calculation.updated_at,
    database_calculation.is_deleted,
    database_calculation.user_id
FROM dc.database_calculation
WHERE database_calculation.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_calculation.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_calculation.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCalculationsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_calculation
WHERE database_calculation.is_deleted = true;

-- name: GetDatabaseCalculations :many
SELECT
    database_calculation.id,
    database_calculation.database_cat_id,
    database_calculation.calculation_type_id,
    database_calculation.created_at,
    database_calculation.updated_at,
    database_calculation.is_deleted,
    database_calculation.user_id
FROM dc.database_calculation
WHERE database_calculation.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN database_calculation.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN database_calculation.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDatabaseCalculations :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.database_calculation
WHERE database_calculation.is_deleted = false;

-- name: UpdateDatabaseCalculationById :exec
UPDATE dc.database_calculation
SET database_cat_id = @database_cat_id,
    calculation_type_id = @calculation_type_id,
    updated_at = now(),
    is_deleted = false,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_calculation.id = @id;

-- name: DeleteDatabaseCalculationById :exec
UPDATE dc.database_calculation
SET updated_at = now(),
    is_deleted = true,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_calculation.id = @id;

-- name: UndeleteDatabaseCalculationById :exec
UPDATE dc.database_calculation
SET updated_at = now(),
    is_deleted = false,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE database_calculation.id = @id;

-- =========================================================
-- dc.domain_cat
-- =========================================================

-- name: CreateDomainCat :one
INSERT INTO dc.domain_cat (
    domain_name,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @domain_name,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetDomainCatsSearchName :many
SELECT
    domain_cat.id,
    domain_cat.domain_name,
    domain_cat.is_deleted,
    domain_cat.created_at,
    domain_cat.updated_at,
    domain_cat.user_id
FROM dc.domain_cat
WHERE domain_cat.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN domain_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN domain_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDomainCatsSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.domain_cat
WHERE domain_cat.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: GetDomainCatById :one
SELECT
    domain_cat.id,
    domain_cat.domain_name,
    domain_cat.is_deleted,
    domain_cat.created_at,
    domain_cat.updated_at,
    domain_cat.user_id
FROM dc.domain_cat
WHERE domain_cat.id = @id
LIMIT 1;

-- name: GetDomainCats :many
SELECT
    domain_cat.id,
    domain_cat.domain_name,
    domain_cat.is_deleted,
    domain_cat.created_at,
    domain_cat.updated_at,
    domain_cat.user_id
FROM dc.domain_cat
WHERE domain_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN domain_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN domain_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDomainCats :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.domain_cat
WHERE domain_cat.is_deleted = false;

-- name: GetDomainCatsDeleted :many
SELECT
    domain_cat.id,
    domain_cat.domain_name,
    domain_cat.is_deleted,
    domain_cat.created_at,
    domain_cat.updated_at,
    domain_cat.user_id
FROM dc.domain_cat
WHERE domain_cat.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN domain_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN domain_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetDomainCatsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.domain_cat
WHERE domain_cat.is_deleted = true;

-- name: UpdateDomainCatById :exec
UPDATE dc.domain_cat
SET domain_name = @domain_name,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE domain_cat.id = @id;

-- name: DeleteDomainCatById :exec
UPDATE dc.domain_cat
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE domain_cat.id = @id;

-- name: UndeleteDomainCatById :exec
UPDATE dc.domain_cat
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE domain_cat.id = @id;

-- =========================================================
-- dc.following_calculation
-- =========================================================

-- name: CreateFollowingCalculation :one
INSERT INTO dc.following_calculation (
    column_cat_id,
    calculation_type_id,
    created_at,
    updated_at,
    is_deleted,
    user_id
) VALUES (
    @column_cat_id,
    @calculation_type_id,
    now(),
    now(),
    false,
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetFollowingCalculationById :one
SELECT
    following_calculation.id,
    following_calculation.column_cat_id,
    following_calculation.calculation_type_id,
    following_calculation.created_at,
    following_calculation.updated_at,
    following_calculation.is_deleted,
    following_calculation.user_id
FROM dc.following_calculation
WHERE following_calculation.id = @id
LIMIT 1;

-- name: GetFollowingCalculations :many
SELECT
    following_calculation.id,
    following_calculation.column_cat_id,
    following_calculation.calculation_type_id,
    following_calculation.created_at,
    following_calculation.updated_at,
    following_calculation.is_deleted,
    following_calculation.user_id
FROM dc.following_calculation
WHERE following_calculation.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN following_calculation.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN following_calculation.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetFollowingCalculations :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.following_calculation
WHERE following_calculation.is_deleted = false;

-- name: GetFollowingCalculationsDeleted :many
SELECT
    following_calculation.id,
    following_calculation.column_cat_id,
    following_calculation.calculation_type_id,
    following_calculation.created_at,
    following_calculation.updated_at,
    following_calculation.is_deleted,
    following_calculation.user_id
FROM dc.following_calculation
WHERE following_calculation.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN following_calculation.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN following_calculation.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetFollowingCalculationsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.following_calculation
WHERE following_calculation.is_deleted = true;

-- name: GetFollowingCalculationsFilters :many
SELECT
    following_calculation.id,
    following_calculation.column_cat_id,
    following_calculation.calculation_type_id,
    following_calculation.created_at,
    following_calculation.updated_at,
    following_calculation.is_deleted,
    following_calculation.user_id
FROM dc.following_calculation
WHERE (sqlc.narg('column_cat_id')::bigint IS NULL OR following_calculation.column_cat_id = sqlc.narg('column_cat_id'))
  AND (sqlc.narg('calculation_type_id')::bigint IS NULL OR following_calculation.calculation_type_id = sqlc.narg('calculation_type_id'))
  AND following_calculation.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN following_calculation.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN following_calculation.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetFollowingCalculationsFilters :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.following_calculation
WHERE (sqlc.narg('column_cat_id')::bigint IS NULL OR following_calculation.column_cat_id = sqlc.narg('column_cat_id'))
  AND (sqlc.narg('calculation_type_id')::bigint IS NULL OR following_calculation.calculation_type_id = sqlc.narg('calculation_type_id'))
  AND following_calculation.is_deleted = false;

-- name: UpdateFollowingCalculationById :exec
UPDATE dc.following_calculation
SET column_cat_id = @column_cat_id,
    calculation_type_id = @calculation_type_id,
    updated_at = now(),
    is_deleted = false,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE following_calculation.id = @id;

-- name: DeleteFollowingCalculationById :exec
UPDATE dc.following_calculation
SET updated_at = now(),
    is_deleted = true,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE following_calculation.id = @id;

-- name: UndeleteFollowingCalculationById :exec
UPDATE dc.following_calculation
SET updated_at = now(),
    is_deleted = false,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE following_calculation.id = @id;

-- =========================================================
-- dc.group_levels
-- =========================================================

-- name: CreateGroupLevels :one
INSERT INTO dc.group_levels (
    column_id,
    parent_column_id,
    level,
    description,
    created_at,
    updated_at,
    is_deleted,
    user_id
) VALUES (
    @column_id,
    @parent_column_id,
    @level,
    @description,
    now(),
    now(),
    false,
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetGroupLevelsSearchName :many
SELECT
    group_levels.id,
    group_levels.column_id,
    group_levels.parent_column_id,
    group_levels.level,
    group_levels.description,
    group_levels.created_at,
    group_levels.updated_at,
    group_levels.is_deleted,
    group_levels.user_id
FROM dc.group_levels
WHERE group_levels.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN group_levels.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN group_levels.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetGroupLevelsSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.group_levels
WHERE group_levels.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: GetGroupLevels :many
SELECT
    group_levels.id,
    group_levels.column_id,
    group_levels.parent_column_id,
    group_levels.level,
    group_levels.description,
    group_levels.created_at,
    group_levels.updated_at,
    group_levels.is_deleted,
    group_levels.user_id
FROM dc.group_levels
WHERE group_levels.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN group_levels.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN group_levels.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetGroupLevels :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.group_levels
WHERE group_levels.is_deleted = false;

-- name: GetGroupLevelsById :one
SELECT
    group_levels.id,
    group_levels.column_id,
    group_levels.parent_column_id,
    group_levels.level,
    group_levels.description,
    group_levels.created_at,
    group_levels.updated_at,
    group_levels.is_deleted,
    group_levels.user_id
FROM dc.group_levels
WHERE group_levels.id = @id
LIMIT 1;

-- name: GetGroupLevelsFiltered :many
SELECT
    group_levels.id,
    group_levels.column_id,
    group_levels.parent_column_id,
    group_levels.level,
    group_levels.description,
    group_levels.created_at,
    group_levels.updated_at,
    group_levels.is_deleted,
    group_levels.user_id
FROM dc.group_levels
WHERE (sqlc.narg('column_id')::bigint IS NULL OR group_levels.column_id = sqlc.narg('column_id'))
  AND (sqlc.narg('parent_column_id')::bigint IS NULL OR group_levels.parent_column_id = sqlc.narg('parent_column_id'))
  AND group_levels.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN group_levels.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN group_levels.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetGroupLevelsFiltered :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.group_levels
WHERE (sqlc.narg('column_id')::bigint IS NULL OR group_levels.column_id = sqlc.narg('column_id'))
  AND (sqlc.narg('parent_column_id')::bigint IS NULL OR group_levels.parent_column_id = sqlc.narg('parent_column_id'))
  AND group_levels.is_deleted = false;

-- name: GetGroupLevelsDeleted :many
SELECT
    group_levels.id,
    group_levels.column_id,
    group_levels.parent_column_id,
    group_levels.level,
    group_levels.description,
    group_levels.created_at,
    group_levels.updated_at,
    group_levels.is_deleted,
    group_levels.user_id
FROM dc.group_levels
WHERE group_levels.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN group_levels.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN group_levels.id END DESC;

-- name: UpdateGroupLevelsById :exec
UPDATE dc.group_levels
SET column_id = @column_id,
    parent_column_id = @parent_column_id,
    level = @level,
    description = @description,
    updated_at = now(),
    is_deleted = false,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE group_levels.id = @id;

-- name: DeleteGroupLevelsById :exec
UPDATE dc.group_levels
SET updated_at = now(),
    is_deleted = true,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE group_levels.id = @id;

-- name: UndeleteGroupLevelsById :exec
UPDATE dc.group_levels
SET updated_at = now(),
    is_deleted = false,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE group_levels.id = @id;

-- =========================================================
-- dc.has_to_group
-- =========================================================

-- name: CreateHasToGroup :one
INSERT INTO dc.has_to_group (
    column_id_a,
    column_id_b,
    description,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @column_id_a,
    @column_id_b,
    @description,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetHasToGroupById :one
SELECT
    has_to_group.id,
    has_to_group.column_id_a,
    has_to_group.column_id_b,
    has_to_group.description,
    has_to_group.is_deleted,
    has_to_group.created_at,
    has_to_group.updated_at,
    has_to_group.user_id
FROM dc.has_to_group
WHERE has_to_group.id = @id
LIMIT 1;

-- name: GetHasToGroupsSearchName :many
SELECT
    has_to_group.id,
    has_to_group.column_id_a,
    has_to_group.column_id_b,
    has_to_group.description,
    has_to_group.is_deleted,
    has_to_group.created_at,
    has_to_group.updated_at,
    has_to_group.user_id
FROM dc.has_to_group
WHERE has_to_group.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%')
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN has_to_group.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN has_to_group.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetHasToGroupsSearchName :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.has_to_group
WHERE has_to_group.is_deleted = false
  AND (lower(name) LIKE '%' || lower(sqlc.arg(search_name)) || '%');

-- name: GetHasToGroups :many
SELECT
    has_to_group.id,
    has_to_group.column_id_a,
    has_to_group.column_id_b,
    has_to_group.description,
    has_to_group.is_deleted,
    has_to_group.created_at,
    has_to_group.updated_at,
    has_to_group.user_id
FROM dc.has_to_group
WHERE has_to_group.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN has_to_group.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN has_to_group.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetHasToGroups :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.has_to_group
WHERE has_to_group.is_deleted = false;

-- name: GetHasToGroupsDeleted :many
SELECT
    has_to_group.id,
    has_to_group.column_id_a,
    has_to_group.column_id_b,
    has_to_group.description,
    has_to_group.is_deleted,
    has_to_group.created_at,
    has_to_group.updated_at,
    has_to_group.user_id
FROM dc.has_to_group
WHERE has_to_group.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN has_to_group.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN has_to_group.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetHasToGroupsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.has_to_group
WHERE has_to_group.is_deleted = true;

-- name: UpdateHasToGroupById :exec
UPDATE dc.has_to_group
SET column_id_a = @column_id_a,
    column_id_b = @column_id_b,
    description = @description,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE has_to_group.id = @id;

-- name: DeleteHasToGroupById :exec
UPDATE dc.has_to_group
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE has_to_group.id = @id;

-- name: UndeleteHasToGroupById :exec
UPDATE dc.has_to_group
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE has_to_group.id = @id;

-- =========================================================
-- dc.schema_cat
-- =========================================================

-- name: CreateSchemaCat :one
INSERT INTO dc.schema_cat (
    database_id,
    name,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @database_id,
    @name,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetSchemaCatById :one
SELECT
    schema_cat.id,
    schema_cat.database_id,
    schema_cat.name,
    schema_cat.is_deleted,
    schema_cat.created_at,
    schema_cat.updated_at,
    schema_cat.user_id
FROM dc.schema_cat
WHERE schema_cat.id = @id
LIMIT 1;

-- name: GetSchemaCats :many
SELECT
    schema_cat.id,
    schema_cat.database_id,
    schema_cat.name,
    schema_cat.is_deleted,
    schema_cat.created_at,
    schema_cat.updated_at,
    schema_cat.user_id
FROM dc.schema_cat
WHERE schema_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN schema_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN schema_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetSchemaCats :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.schema_cat
WHERE schema_cat.is_deleted = false;

-- name: GetSchemaCatsDeleted :many
SELECT
    schema_cat.id,
    schema_cat.database_id,
    schema_cat.name,
    schema_cat.is_deleted,
    schema_cat.created_at,
    schema_cat.updated_at,
    schema_cat.user_id
FROM dc.schema_cat
WHERE schema_cat.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN schema_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN schema_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetSchemaCatsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.schema_cat
WHERE schema_cat.is_deleted = true;

-- name: UpdateSchemaCatById :exec
UPDATE dc.schema_cat
SET database_id = @database_id,
    name = @name,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE schema_cat.id = @id;

-- name: DeleteSchemaCatById :exec
UPDATE dc.schema_cat
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE schema_cat.id = @id;

-- name: UndeleteSchemaCatById :exec
UPDATE dc.schema_cat
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE schema_cat.id = @id;

-- =========================================================
-- dc.table_cat
-- =========================================================

-- name: CreateTableCat :one
INSERT INTO dc.table_cat (
    name,
    description,
    schema_id,
    table_type_id,
    domain_id,
    is_deleted,
    created_at,
    updated_at,
    is_get_dict,
    user_id
) VALUES (
    @name,
    @description,
    @schema_id,
    @table_type_id,
    @domain_id,
    false,
    now(),
    now(),
    false,
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetTableCatById :one
SELECT
    table_cat.id,
    table_cat.name,
    table_cat.description,
    table_cat.schema_id,
    table_cat.table_type_id,
    table_cat.domain_id,
    table_cat.is_deleted,
    table_cat.created_at,
    table_cat.updated_at,
    table_cat.is_get_dict,
    table_cat.user_id
FROM dc.table_cat
WHERE table_cat.id = @id
LIMIT 1;

-- name: GetTableCats :many
SELECT
    table_cat.id,
    table_cat.name,
    table_cat.description,
    table_cat.schema_id,
    table_cat.table_type_id,
    table_cat.domain_id,
    table_cat.is_deleted,
    table_cat.created_at,
    table_cat.updated_at,
    table_cat.is_get_dict,
    table_cat.user_id
FROM dc.table_cat
WHERE table_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN table_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN table_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetTableCats :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.table_cat
WHERE table_cat.is_deleted = false;

-- name: GetTableCatsDeleted :many
SELECT
    table_cat.id,
    table_cat.name,
    table_cat.description,
    table_cat.schema_id,
    table_cat.table_type_id,
    table_cat.domain_id,
    table_cat.is_deleted,
    table_cat.created_at,
    table_cat.updated_at,
    table_cat.is_get_dict,
    table_cat.user_id
FROM dc.table_cat
WHERE table_cat.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN table_cat.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN table_cat.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetTableCatsDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.table_cat
WHERE table_cat.is_deleted = true;

-- name: UpdateTableCatById :exec
UPDATE dc.table_cat
SET name = @name,
    description = @description,
    schema_id = @schema_id,
    table_type_id = @table_type_id,
    domain_id = @domain_id,
    is_deleted = false,
    updated_at = now(),
    is_get_dict = @is_get_dict,
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE table_cat.id = @id;

-- name: DeleteTableCatById :exec
UPDATE dc.table_cat
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE table_cat.id = @id;

-- name: UndeleteTableCatById :exec
UPDATE dc.table_cat
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE table_cat.id = @id;

-- =========================================================
-- dc.table_type
-- =========================================================

-- name: CreateTableType :one
INSERT INTO dc.table_type (
    name,
    description,
    is_deleted,
    created_at,
    updated_at,
    user_id
) VALUES (
    @name,
    @description,
    false,
    now(),
    now(),
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetTableTypeById :one
SELECT
    table_type.id,
    table_type.name,
    table_type.description,
    table_type.is_deleted,
    table_type.created_at,
    table_type.updated_at,
    table_type.user_id
FROM dc.table_type
WHERE table_type.id = @id
LIMIT 1;

-- name: GetTableTypes :many
SELECT
    table_type.id,
    table_type.name,
    table_type.description,
    table_type.is_deleted,
    table_type.created_at,
    table_type.updated_at,
    table_type.user_id
FROM dc.table_type
WHERE table_type.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN table_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN table_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetTableTypes :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.table_type
WHERE table_type.is_deleted = false;

-- name: GetTableTypesDeleted :many
SELECT
    table_type.id,
    table_type.name,
    table_type.description,
    table_type.is_deleted,
    table_type.created_at,
    table_type.updated_at,
    table_type.user_id
FROM dc.table_type
WHERE table_type.is_deleted = true
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN table_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN table_type.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetTableTypesDeleted :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.table_type
WHERE table_type.is_deleted = true;

-- name: UpdateTableTypeById :exec
UPDATE dc.table_type
SET name = @name,
    description = @description,
    is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE table_type.id = @id;

-- name: DeleteTableTypeById :exec
UPDATE dc.table_type
SET is_deleted = true,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE table_type.id = @id;

-- name: UndeleteTableTypeById :exec
UPDATE dc.table_type
SET is_deleted = false,
    updated_at = now(),
    user_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE table_type.id = @id;

-- =========================================================
-- dc.calculation_type
-- =========================================================

-- name: GetCalculationTypeById :one
SELECT
    calculation_type.id,
    calculation_type.name,
    calculation_type.description,
    calculation_type.created_at,
    calculation_type.updated_at,
    calculation_type.is_deleted
FROM dc.calculation_type
WHERE calculation_type.id = @id
LIMIT 1;

-- name: GetCalculationTypes :many
SELECT
    calculation_type.id,
    calculation_type.name,
    calculation_type.description,
    calculation_type.created_at,
    calculation_type.updated_at,
    calculation_type.is_deleted
FROM dc.calculation_type
WHERE calculation_type.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN calculation_type.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN calculation_type.id END DESC;
