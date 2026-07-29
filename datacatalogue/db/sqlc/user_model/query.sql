-- =========================================================
-- dc.user
-- =========================================================

-- name: GetUserById :one
SELECT *
FROM dc."user"
WHERE id = $1
AND is_deleted = false;

-- name: GetUserByExternalId :one
SELECT *
FROM dc."user"
WHERE external_id = $1
AND is_deleted = false;

-- name: GetUsers :many
SELECT *
FROM dc."user"
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedUserById :one
SELECT *
FROM dc."user"
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedUsers :many
SELECT *
FROM dc."user"
WHERE is_deleted = true
ORDER BY id;

-- name: CreateUser :one
INSERT INTO dc."user"
("name", external_id, is_deleted, created_at, updated_at)
VALUES($1, $2, false, now(), now())
RETURNING *;

-- name: UpdateUserById :one
UPDATE dc."user"
SET "name"=$2, external_id=$3, updated_at=now()
WHERE id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteUserById :exec
UPDATE dc."user"
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteUserById :exec
UPDATE dc."user"
SET is_deleted=false, updated_at=now()
WHERE id=$1;
