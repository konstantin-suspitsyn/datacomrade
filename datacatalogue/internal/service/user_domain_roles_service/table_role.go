package userdomainrolesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetTableRoleById возвращает активную строку dc.table_roles по id.
func (s *UserDomainRolesService) GetTableRoleById(ctx context.Context, id int64) (user_domain_roles.DcTableRole, error) {
	row, err := s.UserDomainRolesRepository.GetTableRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcTableRole{}, fmt.Errorf("dc.table_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcTableRole{}, fmt.Errorf("get dc.table_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedTableRoleById возвращает мягко удалённую строку dc.table_roles по id.
func (s *UserDomainRolesService) GetDeletedTableRoleById(ctx context.Context, id int64) (user_domain_roles.DcTableRole, error) {
	row, err := s.UserDomainRolesRepository.GetDeletedTableRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcTableRole{}, fmt.Errorf("deleted dc.table_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcTableRole{}, fmt.Errorf("get deleted dc.table_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetTableRoles возвращает строки dc.table_roles.
func (s *UserDomainRolesService) GetTableRoles(ctx context.Context) ([]user_domain_roles.DcTableRole, error) {
	rows, err := s.UserDomainRolesRepository.GetTableRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("GetTableRoles: %w", err)
	}

	return rows, nil
}

// GetDeletedTableRoles возвращает строки dc.table_roles.
func (s *UserDomainRolesService) GetDeletedTableRoles(ctx context.Context) ([]user_domain_roles.DcTableRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDeletedTableRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("GetDeletedTableRoles: %w", err)
	}

	return rows, nil
}

// CreateTableRole вставляет строку dc.table_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateTableRole(ctx context.Context, params user_domain_roles.CreateTableRoleParams) (user_domain_roles.DcTableRole, error) {
	row, err := s.UserDomainRolesRepository.CreateTableRole(ctx, params)

	if err != nil {
		return user_domain_roles.DcTableRole{}, fmt.Errorf("%w: dc.table_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateTableRoleById обновляет активную строку dc.table_roles и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *UserDomainRolesService) UpdateTableRoleById(ctx context.Context, params user_domain_roles.UpdateTableRoleByIdParams) (user_domain_roles.DcTableRole, error) {
	row, err := s.UserDomainRolesRepository.UpdateTableRoleById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcTableRole{}, fmt.Errorf("dc.table_roles id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcTableRole{}, fmt.Errorf("%w: dc.table_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteTableRoleById мягко удаляет строку dc.table_roles.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserDomainRolesService) DeleteTableRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetTableRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserDomainRolesRepository.DeleteTableRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.table_roles id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteTableRoleById восстанавливает мягко удалённую строку dc.table_roles.
//
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteTableRoleById.
func (s *UserDomainRolesService) UndeleteTableRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedTableRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.UserDomainRolesRepository.UndeleteTableRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.table_roles id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
