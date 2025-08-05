package rolesmodel

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
)

func (m RoleModel) InsertResource(ctx context.Context, resource *Resource) error {

	sqlQuery := `INSERT INTO users.resource ("name", description, is_deleted, created_at, updated_at) 
	VALUES($1, $2, false, now(), now())
	RETURNING id, created_at, updated_at;`

	args := []any{resource.Name, resource.Description}

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&resource.Id, &resource.CreatedAt, &resource.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (m RoleModel) GetResourceById(ctx context.Context, id int64) (*Resource, error) {
	resource := Resource{}

	sqlQuery := `SELECT id, "name", description, is_deleted, created_at, updated_at FROM users.resource
	WHERE id = 1
	RETURNING id, "name", description, is_deleted, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, sqlQuery, resource).Scan(
		&resource.Id,
		&resource.Name,
		&resource.Description,
		&resource.Description,
		&resource.IsDeleted,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &resource, nil
}

// Doesn't actually delete. Just mark deleted
func (m RoleModel) DeleteResourceById(ctx context.Context, id int64) error {

	sqlQuery := `UPDATE users.resource
	SET is_deleted = true,
	created_at = now()
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, configs.QueryTimeoutShort)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, sqlQuery, id)

	if err != nil {
		return err
	}

	return nil
}
