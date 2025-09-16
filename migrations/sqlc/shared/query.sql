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
	updated_at = now(), 
	user_id = $2
WHERE id = $1;

-- name: CountActiveRows :one
SELECT count(*) from shared."domain"
where is_deleted = false;

-- name: GetDomainsWithPager :many
SELECT id, "name", description, user_id, created_at, updated_at
FROM shared."domain"
where is_deleted = false
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateOne :one
UPDATE shared."domain"
	SET 
	"name" = $1,
	description = $2,
	user_id = $3,
	updated_at = now()
WHERE id = $4
RETURNING *;


