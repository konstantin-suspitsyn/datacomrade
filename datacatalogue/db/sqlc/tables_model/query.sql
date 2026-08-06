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
