-- Файл сгенерирован программой SG Buddy https://github.com/konstantin-suspitsyn/sg_buddy
-- Правки будут затёрты при следующей генерации: правьте настройки, а не этот файл.

-- =========================================================
-- dc.user_domain_roles
-- =========================================================

-- name: CreateUserDomainRoles :one
INSERT INTO dc.user_domain_roles (
    user_id,
    domain_roles_id,
    created_at,
    updated_at,
    is_deleted,
    domain_id,
    updated_by_id
) VALUES (
    @user_id,
    @domain_roles_id,
    now(),
    now(),
    false,
    @domain_id,
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetUserDomainRoleses :many
SELECT
    user_domain_roles.id,
    user_domain_roles.user_id,
    user_domain_roles.domain_roles_id,
    user_domain_roles.created_at,
    user_domain_roles.updated_at,
    user_domain_roles.is_deleted,
    user_domain_roles.domain_id,
    user_domain_roles.updated_by_id,
    domain_roles.name AS domain_roles_name,
    domain_roles.description AS domain_roles_description,
    dc_user.name AS dc_user_name
FROM dc.user_domain_roles
INNER JOIN dc.domain_roles domain_roles ON user_domain_roles.domain_roles_id = domain_roles.id
INNER JOIN dc."user" dc_user ON user_domain_roles.user_id = dc_user.id
WHERE user_domain_roles.is_deleted = false
  AND domain_roles.is_deleted = false
  AND dc_user.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN user_domain_roles.user_id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN user_domain_roles.user_id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetUserDomainRoleses :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.user_domain_roles
INNER JOIN dc.domain_roles domain_roles ON user_domain_roles.domain_roles_id = domain_roles.id
INNER JOIN dc."user" dc_user ON user_domain_roles.user_id = dc_user.id
WHERE user_domain_roles.is_deleted = false
  AND domain_roles.is_deleted = false
  AND dc_user.is_deleted = false;

-- name: UpdateUserDomainRolesById :exec
UPDATE dc.user_domain_roles
SET user_id = @user_id,
    domain_roles_id = @domain_roles_id,
    updated_at = now(),
    is_deleted = false,
    domain_id = @domain_id,
    updated_by_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE user_domain_roles.id = @id;

-- name: DeleteUserDomainRolesById :exec
UPDATE dc.user_domain_roles
SET is_deleted = true,
    updated_by_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE user_domain_roles.id = @id;

-- name: UndeleteUserDomainRolesById :exec
UPDATE dc.user_domain_roles
SET updated_at = now(),
    is_deleted = false,
    updated_by_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE user_domain_roles.id = @id;

-- =========================================================
-- dc.user_table_roles
-- =========================================================

-- name: CreateUserTableRoles :one
INSERT INTO dc.user_table_roles (
    user_id,
    table_roles_id,
    created_at,
    updated_at,
    is_deleted,
    table_id,
    updated_by_id
) VALUES (
    @user_id,
    @table_roles_id,
    now(),
    now(),
    false,
    @table_id,
    (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
)
RETURNING *;

-- name: GetUserTableRoleses :many
SELECT
    user_table_roles.id,
    user_table_roles.user_id,
    user_table_roles.created_at,
    user_table_roles.updated_at,
    user_table_roles.table_id,
    dc_user.id AS dc_user_id,
    dc_user.name AS dc_user_name,
    dc_user.external_id AS dc_user_external_id,
    table_roles.id AS table_roles_id_2,
    table_roles.name AS table_roles_name,
    table_roles.description AS table_roles_description,
    table_cat.id AS table_cat_id,
    table_cat.name AS table_cat_name
FROM dc.user_table_roles
INNER JOIN dc."user" dc_user ON user_table_roles.user_id = dc_user.id
INNER JOIN dc.table_roles table_roles ON user_table_roles.table_roles_id = table_roles.id
INNER JOIN dc.table_cat table_cat ON user_table_roles.table_id = table_cat.id
WHERE user_table_roles.is_deleted = false
  AND dc_user.is_deleted = false
  AND table_roles.is_deleted = false
  AND table_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN user_table_roles.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN user_table_roles.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetUserTableRoleses :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.user_table_roles
INNER JOIN dc."user" dc_user ON user_table_roles.user_id = dc_user.id
INNER JOIN dc.table_roles table_roles ON user_table_roles.table_roles_id = table_roles.id
INNER JOIN dc.table_cat table_cat ON user_table_roles.table_id = table_cat.id
WHERE user_table_roles.is_deleted = false
  AND dc_user.is_deleted = false
  AND table_roles.is_deleted = false
  AND table_cat.is_deleted = false;

-- name: GetUserTableRolesesFiltered :many
SELECT
    user_table_roles.id,
    user_table_roles.user_id,
    user_table_roles.created_at,
    user_table_roles.updated_at,
    user_table_roles.table_id,
    dc_user.id AS dc_user_id,
    dc_user.name AS dc_user_name,
    dc_user.external_id AS dc_user_external_id,
    table_roles.name AS table_roles_name,
    table_roles.description AS table_roles_description,
    table_cat.id AS table_cat_id,
    table_cat.name AS table_cat_name
FROM dc.user_table_roles
INNER JOIN dc."user" dc_user ON user_table_roles.user_id = dc_user.id
INNER JOIN dc.table_roles table_roles ON user_table_roles.table_roles_id = table_roles.id
INNER JOIN dc.table_cat table_cat ON user_table_roles.table_id = table_cat.id
WHERE user_table_roles.is_deleted = false
  AND (sqlc.narg('dc_user_name')::varchar(512) IS NULL OR dc_user.name = sqlc.narg('dc_user_name'))
  AND dc_user.is_deleted = false
  AND (sqlc.narg('dc_user_external_id')::uuid IS NULL OR dc_user.external_id = sqlc.narg('dc_user_external_id'))
  AND (sqlc.narg('table_roles_name')::varchar(128) IS NULL OR table_roles.name = sqlc.narg('table_roles_name'))
  AND table_roles.is_deleted = false
  AND (sqlc.narg('table_cat_id')::bigint IS NULL OR table_cat.id = sqlc.narg('table_cat_id'))
  AND (sqlc.narg('table_cat_name')::varchar(128) IS NULL OR table_cat.name = sqlc.narg('table_cat_name'))
  AND table_cat.is_deleted = false
ORDER BY CASE WHEN sqlc.arg('order')::text <> 'DESC' THEN user_table_roles.id END ASC,
    CASE WHEN sqlc.arg('order')::text = 'DESC' THEN user_table_roles.id END DESC
LIMIT @page_limit::int OFFSET (sqlc.arg('page')::int-1)*sqlc.arg('page_limit')::int;

-- name: CountGetUserTableRolesesFiltered :one
SELECT
    count(*)                                                        AS total_items,
    ceil(count(*)::numeric / GREATEST(@page_limit::int, 1))::bigint AS total_pages
FROM dc.user_table_roles
INNER JOIN dc."user" dc_user ON user_table_roles.user_id = dc_user.id
INNER JOIN dc.table_roles table_roles ON user_table_roles.table_roles_id = table_roles.id
INNER JOIN dc.table_cat table_cat ON user_table_roles.table_id = table_cat.id
WHERE user_table_roles.is_deleted = false
  AND (sqlc.narg('dc_user_name')::varchar(512) IS NULL OR dc_user.name = sqlc.narg('dc_user_name'))
  AND dc_user.is_deleted = false
  AND (sqlc.narg('dc_user_external_id')::uuid IS NULL OR dc_user.external_id = sqlc.narg('dc_user_external_id'))
  AND (sqlc.narg('table_roles_name')::varchar(128) IS NULL OR table_roles.name = sqlc.narg('table_roles_name'))
  AND table_roles.is_deleted = false
  AND (sqlc.narg('table_cat_id')::bigint IS NULL OR table_cat.id = sqlc.narg('table_cat_id'))
  AND (sqlc.narg('table_cat_name')::varchar(128) IS NULL OR table_cat.name = sqlc.narg('table_cat_name'))
  AND table_cat.is_deleted = false;

-- name: GetUserTableRolesById :one
SELECT
    user_table_roles.id,
    user_table_roles.user_id,
    user_table_roles.table_roles_id,
    user_table_roles.created_at,
    user_table_roles.updated_at,
    user_table_roles.table_id,
    dc_user.id AS dc_user_id,
    dc_user.name AS dc_user_name,
    dc_user.external_id AS dc_user_external_id,
    table_roles.name AS table_roles_name,
    table_roles.description AS table_roles_description,
    table_cat.id AS table_cat_id,
    table_cat.name AS table_cat_name
FROM dc.user_table_roles
INNER JOIN dc."user" dc_user ON user_table_roles.user_id = dc_user.id
INNER JOIN dc.table_roles table_roles ON user_table_roles.table_roles_id = table_roles.id
INNER JOIN dc.table_cat table_cat ON user_table_roles.table_id = table_cat.id
WHERE user_table_roles.id = @id
  AND user_table_roles.is_deleted = false
  AND dc_user.is_deleted = false
  AND table_roles.is_deleted = false
  AND table_cat.is_deleted = false
LIMIT 1;

-- name: UpdateUserTableRolesById :exec
UPDATE dc.user_table_roles
SET user_id = @user_id,
    table_roles_id = @table_roles_id,
    updated_at = now(),
    is_deleted = false,
    table_id = @table_id,
    updated_by_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE user_table_roles.id = @id;

-- name: DeleteUserTableRolesById :exec
UPDATE dc.user_table_roles
SET updated_at = now(),
    is_deleted = true,
    updated_by_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE user_table_roles.id = @id;

-- name: UndeleteUserTableRolesById :exec
UPDATE dc.user_table_roles
SET updated_at = now(),
    is_deleted = false,
    updated_by_id = (SELECT u.id FROM dc."user" u WHERE u.external_id = @external_id)
WHERE user_table_roles.id = @id;
