package userdomainrolesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetTablesTableRoleById возвращает активную строку dc.tables_table_roles по id.
func (s *UserDomainRolesService) GetTablesTableRoleById(ctx context.Context, id int64) (user_domain_roles.DcTablesTableRole, error) {
	row, err := s.UserDomainRolesRepository.GetTablesTableRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcTablesTableRole{}, fmt.Errorf("dc.tables_table_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcTablesTableRole{}, fmt.Errorf("get dc.tables_table_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetTablesTableRoles возвращает все активные строки dc.tables_table_roles.
func (s *UserDomainRolesService) GetTablesTableRoles(ctx context.Context) ([]user_domain_roles.DcTablesTableRole, error) {
	rows, err := s.UserDomainRolesRepository.GetTablesTableRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get dc.tables_table_roles: %w", err)
	}

	return rows, nil
}

// GetDeletedTablesTableRoleById возвращает мягко удалённую строку dc.tables_table_roles по id.
func (s *UserDomainRolesService) GetDeletedTablesTableRoleById(ctx context.Context, id int64) (user_domain_roles.DcTablesTableRole, error) {
	row, err := s.UserDomainRolesRepository.GetDeletedTablesTableRoleById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcTablesTableRole{}, fmt.Errorf("deleted dc.tables_table_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcTablesTableRole{}, fmt.Errorf("get deleted dc.tables_table_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetDeletedTablesTableRoles возвращает все мягко удалённые строки dc.tables_table_roles.
func (s *UserDomainRolesService) GetDeletedTablesTableRoles(ctx context.Context) ([]user_domain_roles.DcTablesTableRole, error) {
	rows, err := s.UserDomainRolesRepository.GetDeletedTablesTableRoles(ctx)

	if err != nil {
		return nil, fmt.Errorf("get deleted dc.tables_table_roles: %w", err)
	}

	return rows, nil
}

// CreateTablesTableRole вставляет строку dc.tables_table_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateTablesTableRole(ctx context.Context, params user_domain_roles.CreateTablesTableRoleParams) (user_domain_roles.DcTablesTableRole, error) {
	row, err := s.UserDomainRolesRepository.CreateTablesTableRole(ctx, params)

	if err != nil {
		return user_domain_roles.DcTablesTableRole{}, fmt.Errorf("%w: dc.tables_table_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateTablesTableRoleById обновляет активную строку dc.tables_table_roles и возвращает её целиком.
//
// Запрос фильтрует по is_deleted = false, поэтому попытка обновить удалённую
// или несуществующую запись даёт sql.ErrNoRows — переводим его в ErrNotFound,
// чтобы api-слой ответил NotFound, а не Internal.
func (s *UserDomainRolesService) UpdateTablesTableRoleById(ctx context.Context, params user_domain_roles.UpdateTablesTableRoleByIdParams) (user_domain_roles.DcTablesTableRole, error) {
	row, err := s.UserDomainRolesRepository.UpdateTablesTableRoleById(ctx, params)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.DcTablesTableRole{}, fmt.Errorf("dc.tables_table_roles id = %d: %w", params.ID, customerrors.ErrNotFound)
		}

		return user_domain_roles.DcTablesTableRole{}, fmt.Errorf("%w: dc.tables_table_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return row, nil
}

// DeleteTablesTableRoleById мягко удаляет строку dc.tables_table_roles.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserDomainRolesService) DeleteTablesTableRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetTablesTableRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserDomainRolesRepository.DeleteTablesTableRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.tables_table_roles id = %d: %w", customerrors.ErrDelete, id, err)
	}

	return nil
}

// UndeleteTablesTableRoleById восстанавливает мягко удалённую строку dc.tables_table_roles.
// Существование удалённой записи проверяется заранее по той же причине,
// что и в DeleteTablesTableRoleById.
func (s *UserDomainRolesService) UndeleteTablesTableRoleById(ctx context.Context, id int64) error {
	if _, err := s.GetDeletedTablesTableRoleById(ctx, id); err != nil {
		return errors.Join(customerrors.ErrUndelete, err)
	}

	if err := s.UserDomainRolesRepository.UndeleteTablesTableRoleById(ctx, id); err != nil {
		return fmt.Errorf("%w: dc.tables_table_roles id = %d: %w", customerrors.ErrUndelete, id, err)
	}

	return nil
}
