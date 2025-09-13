-- name: GetDomains :many
SELECT id, "name", description, user_id, created_at, updated_at
FROM shared."domain"
where is_deleted = false
ORDER BY id;


-- name: GetDomainById :one
SELECT id, "name", description, user_id, created_at, updated_at
FROM shared."domain"
where is_deleted = false
AND id = $1;

-- name: CreateDomain :one
INSERT INTO shared."domain"
("name", description, user_id, is_deleted, created_at, updated_at)
VALUES($1, $2, $3, false, now(), now())
RETURNING *;

-- name: DeleteDomain :exec
UPDATE shared."domain"
	SET is_deleted = true,
	updated_at = now()
WHERE id = $1;

-- name: CountActiveRows :one
SELECT count(*) from shared."domain"
where is_deleted = false;
