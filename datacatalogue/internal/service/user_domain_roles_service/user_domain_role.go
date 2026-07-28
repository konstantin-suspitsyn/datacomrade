package userdomainrolesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetUserDomainRoleById возвращает активную строку dc.user_domain_roles по id.
func (s *UserDomainRolesService) GetUserDomainRoleById(ctx context.Context, id int64) (user_domain_roles.DcUserDomainRole, error) {
	row, err := s.UserDomainRolesRepository.GetUserDomainRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("dc.user_domain_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("get dc.user_domain_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetUserDomainRoles возвращает все активные строки dc.user_domain_roles.
func (s *UserDomainRolesService) GetUserDomainRoles(ctx context.Context) ([]user_domain_roles.DcUserDomainRole, error) {
	rows, err := s.UserDomainRolesRepository.GetUserDomainRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.user_domain_roles: %w", err)
	}

	return rows, nil
}

// GetDeletedUserDomainRoleById возвращает мягко удалённую строку dc.user_domain_roles по id.
func (s *UserDomainRolesService) GetDeletedUserDomainRoleById(ctx context.Context, id int64) (user_domain_roles.DcUserDomainRole, error) {
	row, err := s.UserDomainRolesRepository.GetDeletedUserDomainRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("deleted dc.user_domain_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("get deleted dc.user_domain_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedUserDomainRoles возвращает все мягко удалённые строки dc.user_domain_roles.
func (s *UserDomainRolesService) GetDeletedUserDomainRoles(ctx context.Context) ([]user_domain_roles.DcUserDomainRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDeletedUserDomainRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.user_domain_roles: %w", err)
	}

	return rows, nil
}

// CreateUserDomainRole вставляет строку dc.user_domain_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateUserDomainRole(ctx context.Context, params user_domain_roles.CreateUserDomainRoleParams) (user_domain_roles.DcUserDomainRole, error) {
	row, err := s.UserDomainRolesRepository.CreateUserDomainRole(ctx, params)

	if err != nil {
		return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("%w: dc.user_domain_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateUserDomainRoleById обновляет активную строку dc.user_domain_roles и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *UserDomainRolesService) UpdateUserDomainRoleById(ctx context.Context, params user_domain_roles.UpdateUserDomainRoleByIdParams) (user_domain_roles.DcUserDomainRole, error) {
	row, err := s.UserDomainRolesRepository.UpdateUserDomainRoleById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("dc.user_domain_roles id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("%w: dc.user_domain_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteUserDomainRoleById мягко удаляет строку dc.user_domain_roles.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserDomainRolesService) DeleteUserDomainRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetUserDomainRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserDomainRolesRepository.DeleteUserDomainRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.user_domain_roles id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteUserDomainRoleById восстанавливает мягко удалённую строку dc.user_domain_roles.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteUserDomainRoleById.
func (s *UserDomainRolesService) UndeleteUserDomainRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedUserDomainRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.UserDomainRolesRepository.UndeleteUserDomainRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.user_domain_roles id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
