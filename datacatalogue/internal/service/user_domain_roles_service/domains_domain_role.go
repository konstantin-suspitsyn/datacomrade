package userdomainrolesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetDomainsDomainRoleById возвращает активную строку dc.domains_domain_roles по id.
func (s *UserDomainRolesService) GetDomainsDomainRoleById(ctx context.Context, id int64) (user_domain_roles.DcDomainsDomainRole, error) {
	row, err := s.UserDomainRolesRepository.GetDomainsDomainRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcDomainsDomainRole{}, fmt.Errorf("dc.domains_domain_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcDomainsDomainRole{}, fmt.Errorf("get dc.domains_domain_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDomainsDomainRoles возвращает все активные строки dc.domains_domain_roles.
func (s *UserDomainRolesService) GetDomainsDomainRoles(ctx context.Context) ([]user_domain_roles.DcDomainsDomainRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDomainsDomainRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.domains_domain_roles: %w", err)
	}

	return rows, nil
}

// GetDeletedDomainsDomainRoleById возвращает мягко удалённую строку dc.domains_domain_roles по id.
func (s *UserDomainRolesService) GetDeletedDomainsDomainRoleById(ctx context.Context, id int64) (user_domain_roles.DcDomainsDomainRole, error) {
	row, err := s.UserDomainRolesRepository.GetDeletedDomainsDomainRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcDomainsDomainRole{}, fmt.Errorf("deleted dc.domains_domain_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcDomainsDomainRole{}, fmt.Errorf("get deleted dc.domains_domain_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedDomainsDomainRoles возвращает все мягко удалённые строки dc.domains_domain_roles.
func (s *UserDomainRolesService) GetDeletedDomainsDomainRoles(ctx context.Context) ([]user_domain_roles.DcDomainsDomainRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDeletedDomainsDomainRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.domains_domain_roles: %w", err)
	}

	return rows, nil
}

// CreateDomainsDomainRole вставляет строку dc.domains_domain_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateDomainsDomainRole(ctx context.Context, params user_domain_roles.CreateDomainsDomainRoleParams) (user_domain_roles.DcDomainsDomainRole, error) {
	row, err := s.UserDomainRolesRepository.CreateDomainsDomainRole(ctx, params)

	if err != nil {
		return user_domain_roles.DcDomainsDomainRole{}, fmt.Errorf("%w: dc.domains_domain_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateDomainsDomainRoleById обновляет активную строку dc.domains_domain_roles и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *UserDomainRolesService) UpdateDomainsDomainRoleById(ctx context.Context, params user_domain_roles.UpdateDomainsDomainRoleByIdParams) (user_domain_roles.DcDomainsDomainRole, error) {
	row, err := s.UserDomainRolesRepository.UpdateDomainsDomainRoleById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcDomainsDomainRole{}, fmt.Errorf("dc.domains_domain_roles id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcDomainsDomainRole{}, fmt.Errorf("%w: dc.domains_domain_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteDomainsDomainRoleById мягко удаляет строку dc.domains_domain_roles.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserDomainRolesService) DeleteDomainsDomainRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDomainsDomainRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserDomainRolesRepository.DeleteDomainsDomainRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.domains_domain_roles id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteDomainsDomainRoleById восстанавливает мягко удалённую строку dc.domains_domain_roles.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteDomainsDomainRoleById.
func (s *UserDomainRolesService) UndeleteDomainsDomainRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedDomainsDomainRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.UserDomainRolesRepository.UndeleteDomainsDomainRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.domains_domain_roles id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
