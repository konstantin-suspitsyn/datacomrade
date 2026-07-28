package userdomainrolesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetDomainRoleById возвращает активную строку dc.domain_roles по id.
func (s *UserDomainRolesService) GetDomainRoleById(ctx context.Context, id int64) (user_domain_roles.DcDomainRole, error) {
	row, err := s.UserDomainRolesRepository.GetDomainRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcDomainRole{}, fmt.Errorf("dc.domain_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcDomainRole{}, fmt.Errorf("get dc.domain_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDomainRoles возвращает все активные строки dc.domain_roles.
func (s *UserDomainRolesService) GetDomainRoles(ctx context.Context) ([]user_domain_roles.DcDomainRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDomainRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.domain_roles: %w", err)
	}

	return rows, nil
}

// GetDeletedDomainRoleById возвращает мягко удалённую строку dc.domain_roles по id.
func (s *UserDomainRolesService) GetDeletedDomainRoleById(ctx context.Context, id int64) (user_domain_roles.DcDomainRole, error) {
	row, err := s.UserDomainRolesRepository.GetDeletedDomainRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcDomainRole{}, fmt.Errorf("deleted dc.domain_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcDomainRole{}, fmt.Errorf("get deleted dc.domain_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedDomainRoles возвращает все мягко удалённые строки dc.domain_roles.
func (s *UserDomainRolesService) GetDeletedDomainRoles(ctx context.Context) ([]user_domain_roles.DcDomainRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDeletedDomainRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.domain_roles: %w", err)
	}

	return rows, nil
}

// CreateDomainRole вставляет строку dc.domain_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateDomainRole(ctx context.Context, params user_domain_roles.CreateDomainRoleParams) (user_domain_roles.DcDomainRole, error) {
	row, err := s.UserDomainRolesRepository.CreateDomainRole(ctx, params)

	if err != nil {
		return user_domain_roles.DcDomainRole{}, fmt.Errorf("%w: dc.domain_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDomainRoleById обновляет активную строку dc.domain_roles и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *UserDomainRolesService) UpdateDomainRoleById(ctx context.Context, params user_domain_roles.UpdateDomainRoleByIdParams) (user_domain_roles.DcDomainRole, error) {
	row, err := s.UserDomainRolesRepository.UpdateDomainRoleById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcDomainRole{}, fmt.Errorf("dc.domain_roles id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcDomainRole{}, fmt.Errorf("%w: dc.domain_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteDomainRoleById мягко удаляет строку dc.domain_roles.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserDomainRolesService) DeleteDomainRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDomainRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserDomainRolesRepository.DeleteDomainRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.domain_roles id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteDomainRoleById восстанавливает мягко удалённую строку dc.domain_roles.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteDomainRoleById.
func (s *UserDomainRolesService) UndeleteDomainRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedDomainRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.UserDomainRolesRepository.UndeleteDomainRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.domain_roles id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
