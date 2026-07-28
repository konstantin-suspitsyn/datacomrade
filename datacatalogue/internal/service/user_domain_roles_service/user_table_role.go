package userdomainrolesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetUserTableRoleById возвращает активную строку dc.user_table_roles по id.
func (s *UserDomainRolesService) GetUserTableRoleById(ctx context.Context, id int64) (user_domain_roles.DcUserTableRole, error) {
	row, err := s.UserDomainRolesRepository.GetUserTableRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcUserTableRole{}, fmt.Errorf("dc.user_table_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcUserTableRole{}, fmt.Errorf("get dc.user_table_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetUserTableRoles возвращает все активные строки dc.user_table_roles.
func (s *UserDomainRolesService) GetUserTableRoles(ctx context.Context) ([]user_domain_roles.DcUserTableRole, error) {
	rows, err := s.UserDomainRolesRepository.GetUserTableRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.user_table_roles: %w", err)
	}

	return rows, nil
}

// GetDeletedUserTableRoleById возвращает мягко удалённую строку dc.user_table_roles по id.
func (s *UserDomainRolesService) GetDeletedUserTableRoleById(ctx context.Context, id int64) (user_domain_roles.DcUserTableRole, error) {
	row, err := s.UserDomainRolesRepository.GetDeletedUserTableRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcUserTableRole{}, fmt.Errorf("deleted dc.user_table_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcUserTableRole{}, fmt.Errorf("get deleted dc.user_table_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedUserTableRoles возвращает все мягко удалённые строки dc.user_table_roles.
func (s *UserDomainRolesService) GetDeletedUserTableRoles(ctx context.Context) ([]user_domain_roles.DcUserTableRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDeletedUserTableRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.user_table_roles: %w", err)
	}

	return rows, nil
}

// CreateUserTableRole вставляет строку dc.user_table_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateUserTableRole(ctx context.Context, params user_domain_roles.CreateUserTableRoleParams) (user_domain_roles.DcUserTableRole, error) {
	row, err := s.UserDomainRolesRepository.CreateUserTableRole(ctx, params)

	if err != nil {
		return user_domain_roles.DcUserTableRole{}, fmt.Errorf("%w: dc.user_table_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateUserTableRoleById обновляет активную строку dc.user_table_roles и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *UserDomainRolesService) UpdateUserTableRoleById(ctx context.Context, params user_domain_roles.UpdateUserTableRoleByIdParams) (user_domain_roles.DcUserTableRole, error) {
	row, err := s.UserDomainRolesRepository.UpdateUserTableRoleById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcUserTableRole{}, fmt.Errorf("dc.user_table_roles id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcUserTableRole{}, fmt.Errorf("%w: dc.user_table_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteUserTableRoleById мягко удаляет строку dc.user_table_roles.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserDomainRolesService) DeleteUserTableRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetUserTableRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserDomainRolesRepository.DeleteUserTableRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.user_table_roles id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteUserTableRoleById восстанавливает мягко удалённую строку dc.user_table_roles.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteUserTableRoleById.
func (s *UserDomainRolesService) UndeleteUserTableRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedUserTableRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.UserDomainRolesRepository.UndeleteUserTableRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.user_table_roles id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
