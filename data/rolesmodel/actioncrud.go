package rolesmodel

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
)

func (m RoleModel) InsertAction(ctx context.Context, action *Action) error {

	sqlQuery := `INSERT INTO users."action" ("name", description, is_deleted, created_at, updated_at) VALUES($1, $2, false, now(), now())
	RETURNING id, created_at, updated_at;`

	args := []any{action.Name, action.Description}

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&action.Id,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (m RoleModel) GetActionById(ctx context.Context, id int64) (*Action, error) {
	action := Action{}

	sqlQuery := `SELECT id, "name", description, is_deleted, created_at, updated_at FROM users."action"
	WHERE id = $1
	RETURNING id, "name", description, is_deleted, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, sqlQuery, action).Scan(
		&action.Id,
		&action.Name,
		&action.Description,
		&action.IsDeleted,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &action, nil
}

// Doesn't actually delete. Just mark deleted
func (m RoleModel) DeleteActionById(ctx context.Context, id int64) error {

	sqlQuery := `UPDATE users.action
	SET is_deleted = true,
	updated_at = now()
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, sqlQuery, id)

	if err != nil {
		return err
	}

	return nil
}
