-- name: GetRoleAccessById :one
SELECT * FROM  users.role_access
WHERE id = $1 LIMIT 1;

-- name: ListUserAccess :many
SELECT * FROM users.role_access
ORDER BY id;

-- name: CreateRoleAccess :one
INSERT INTO users.role_access (
	role_id, resource_id, action_id, resource_type_id, is_deleted
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteAuthor :exec
UPDATE users.role_access
	SET is_deleted = false,
	updated_at = now()
WHERE id = $1;

-- name: GetUserAccessById :one
SELECT * FROM  users.user_access
WHERE id = $1 LIMIT 1;

-- name: GetJWTShortRolesByUserId :many
SELECT user_role.id, user_role.user_id, user_role.role_id, role.role_name_short 
FROM users.user_role
inner join users."role"
	on "role".id = user_role.role_id
	and "role".is_deleted = false
	and "role".jwt_export = true
where user_id = $1
and user_role.is_deleted = false;
