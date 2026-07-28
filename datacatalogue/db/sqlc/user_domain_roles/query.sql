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
-- dc.domains_domain_roles
-- =========================================================

-- name: GetDomainsDomainRoleById :one
SELECT *
FROM dc.domains_domain_roles
WHERE id = $1
AND is_deleted = false;

-- name: GetDomainsDomainRoles :many
SELECT *
FROM dc.domains_domain_roles
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedDomainsDomainRoleById :one
SELECT *
FROM dc.domains_domain_roles
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedDomainsDomainRoles :many
SELECT *
FROM dc.domains_domain_roles
WHERE is_deleted = true
ORDER BY id;

-- name: CreateDomainsDomainRole :one
INSERT INTO dc.domains_domain_roles
(domain_cat_id, domain_roles_id, created_at, updated_at, is_deleted)
VALUES($1, $2, now(), now(), false)
RETURNING *;

-- name: UpdateDomainsDomainRoleById :one
UPDATE dc.domains_domain_roles
SET domain_cat_id=$2, domain_roles_id=$3, updated_at=now()
WHERE id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteDomainsDomainRoleById :exec
UPDATE dc.domains_domain_roles
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteDomainsDomainRoleById :exec
UPDATE dc.domains_domain_roles
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
-- dc.tables_table_roles
-- =========================================================

-- name: GetTablesTableRoleById :one
SELECT *
FROM dc.tables_table_roles
WHERE id = $1
AND is_deleted = false;

-- name: GetTablesTableRoles :many
SELECT *
FROM dc.tables_table_roles
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedTablesTableRoleById :one
SELECT *
FROM dc.tables_table_roles
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedTablesTableRoles :many
SELECT *
FROM dc.tables_table_roles
WHERE is_deleted = true
ORDER BY id;

-- name: CreateTablesTableRole :one
INSERT INTO dc.tables_table_roles
(table_cat_id, table_roles_id, created_at, updated_at, is_deleted)
VALUES($1, $2, now(), now(), false)
RETURNING *;

-- name: UpdateTablesTableRoleById :one
UPDATE dc.tables_table_roles
SET table_cat_id=$2, table_roles_id=$3, updated_at=now()
WHERE id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteTablesTableRoleById :exec
UPDATE dc.tables_table_roles
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteTablesTableRoleById :exec
UPDATE dc.tables_table_roles
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
(user_id, domain_roles_id, created_at, updated_at, is_deleted)
VALUES($1, $2, now(), now(), false)
RETURNING *;

-- name: UpdateUserDomainRoleById :one
UPDATE dc.user_domain_roles
SET user_id=$2, domain_roles_id=$3, updated_at=now()
WHERE id=$1
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
(user_id, table_roles_id, created_at, updated_at, is_deleted)
VALUES($1, $2, now(), now(), false)
RETURNING *;

-- name: UpdateUserTableRoleById :one
UPDATE dc.user_table_roles
SET user_id=$2, table_roles_id=$3, updated_at=now()
WHERE id=$1
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
