-- =========================================================
-- dc.domain_roles
-- =========================================================

-- name: GetDomainRoleById :one
SELECT *
FROM dc.domain_roles
WHERE id = $1
AND is_deleted = false;

-- name: GetDomainRoles :many
SELECT *
FROM dc.domain_roles
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedDomainRoleById :one
SELECT *
FROM dc.domain_roles
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedDomainRoles :many
SELECT *
FROM dc.domain_roles
WHERE is_deleted = true
ORDER BY id;

-- name: CreateDomainRole :one
INSERT INTO dc.domain_roles
("name", description, created_at, updated_at, is_deleted)
VALUES($1, $2, now(), now(), false)
RETURNING *;

-- name: UpdateDomainRoleById :one
UPDATE dc.domain_roles
SET "name"=$2, description=$3, updated_at=now()
WHERE id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteDomainRoleById :exec
UPDATE dc.domain_roles
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteDomainRoleById :exec
UPDATE dc.domain_roles
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.table_roles
-- =========================================================

-- name: GetTableRoleById :one
SELECT *
FROM dc.table_roles
WHERE id = $1
AND is_deleted = false;

-- name: GetTableRoles :many
SELECT *
FROM dc.table_roles
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedTableRoleById :one
SELECT *
FROM dc.table_roles
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedTableRoles :many
SELECT *
FROM dc.table_roles
WHERE is_deleted = true
ORDER BY id;

-- name: CreateTableRole :one
INSERT INTO dc.table_roles
("name", description, created_at, updated_at, is_deleted)
VALUES($1, $2, now(), now(), false)
RETURNING *;

-- name: UpdateTableRoleById :one
UPDATE dc.table_roles
SET "name"=$2, description=$3, updated_at=now()
WHERE id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteTableRoleById :exec
UPDATE dc.table_roles
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteTableRoleById :exec
UPDATE dc.table_roles
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.user_domain_roles
-- =========================================================

-- name: GetUserDomainRoleById :one
SELECT *
FROM dc.user_domain_roles
WHERE id = $1
AND is_deleted = false;

-- name: GetUserDomainRoles :many
SELECT *
FROM dc.user_domain_roles
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedUserDomainRoleById :one
SELECT *
FROM dc.user_domain_roles
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedUserDomainRoles :many
SELECT *
FROM dc.user_domain_roles
WHERE is_deleted = true
ORDER BY id;

-- name: CreateUserDomainRole :one
INSERT INTO dc.user_domain_roles
(user_id, domain_roles_id, created_at, updated_at, is_deleted, domain_id, updated_by_id)
VALUES($1, $2, now(), now(), false, $3, (SELECT u.id FROM dc."user" u WHERE u.external_id = $4))
RETURNING *;

-- name: UpdateUserDomainRoleById :one
UPDATE dc.user_domain_roles
SET user_id=$2, domain_roles_id=$3, domain_id=$4, updated_by_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $5), updated_at=now()
WHERE dc.user_domain_roles.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteUserDomainRoleById :exec
UPDATE dc.user_domain_roles
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteUserDomainRoleById :exec
UPDATE dc.user_domain_roles
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.user_table_roles
-- =========================================================

-- name: GetUserTableRoleById :one
SELECT *
FROM dc.user_table_roles
WHERE id = $1
AND is_deleted = false;

-- name: GetUserTableRoles :many
SELECT *
FROM dc.user_table_roles
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedUserTableRoleById :one
SELECT *
FROM dc.user_table_roles
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedUserTableRoles :many
SELECT *
FROM dc.user_table_roles
WHERE is_deleted = true
ORDER BY id;

-- name: CreateUserTableRole :one
INSERT INTO dc.user_table_roles
(user_id, table_roles_id, created_at, updated_at, is_deleted, table_id, updated_by_id)
VALUES($1, $2, now(), now(), false, $3, (SELECT u.id FROM dc."user" u WHERE u.external_id = $4))
RETURNING *;

-- name: UpdateUserTableRoleById :one
UPDATE dc.user_table_roles
SET user_id=$2, table_roles_id=$3, table_id=$4, updated_by_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $5), updated_at=now()
WHERE dc.user_table_roles.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteUserTableRoleById :exec
UPDATE dc.user_table_roles
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteUserTableRoleById :exec
UPDATE dc.user_table_roles
SET is_deleted=false, updated_at=now()
WHERE id=$1;
