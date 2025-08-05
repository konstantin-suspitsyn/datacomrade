package rolesmodel

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
)

func (m RoleModel) InsertRole(ctx context.Context, role *Role) error {

	sqlQuery := `INSERT INTO users."role" (role_name_long, role_name_short, description, jwt_export, is_deleted, created_at, updated_at) VALUES($1, $2, $3, $4, $5, now(), now())
	RETURNING id, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	args := []any{role.RoleNameLong, role.RoleNameShort, role.Description, role.JwtExport, role.IsDeleted}

	err := m.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&role.Id, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (m RoleModel) GetRoleById(ctx context.Context, id int64) (*Role, error) {

	role := Role{}

	sqlQuery := `SELECT id, role_name_long, role_name_short, description, jwt_export, is_deleted, created_at, updated_at FROM users."role"
	WHERE id = $1
	RETURNING id, role_name_long, role_name_short, description, jwt_export, is_deleted, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, sqlQuery, role).Scan(
		&role.Id,
		&role.RoleNameLong,
		&role.RoleNameShort,
		&role.Description,
		&role.JwtExport,
		&role.IsDeleted,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// Doesn't actually delete. Just mark deleted
func (m RoleModel) DeleteRoleById(ctx context.Context, id int64) error {

	sqlQuery := `
	UPDATE users.role
	SET is_deleted = true,
	created_at = now()
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	return nil
}
