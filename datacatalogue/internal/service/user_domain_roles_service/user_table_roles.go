package userdomainrolesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetUserTableRolesById возвращает активную строку dc.user_table_roles по id.
func (s *UserDomainRolesService) GetUserTableRolesById(ctx context.Context, id int64) (user_domain_roles.GetUserTableRolesByIdRow, error) {
	row, err := s.UserDomainRolesRepository.GetUserTableRolesById(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user_domain_roles.GetUserTableRolesByIdRow{}, fmt.Errorf("dc.user_table_roles id = %d: %w", id, customerrors.ErrNotFound)
		}

		return user_domain_roles.GetUserTableRolesByIdRow{}, fmt.Errorf("get dc.user_table_roles id = %d: %w", id, err)
	}

	return row, nil
}

// GetUserTableRoleses возвращает страницу строк dc.user_table_roles и её счётчики.
func (s *UserDomainRolesService) GetUserTableRoleses(ctx context.Context, params user_domain_roles.GetUserTableRolesesParams) ([]user_domain_roles.GetUserTableRolesesRow, user_domain_roles.CountGetUserTableRolesesRow, error) {
	count, err := s.UserDomainRolesRepository.CountGetUserTableRoleses(ctx, params.PageLimit)
	if err != nil {
		return nil, user_domain_roles.CountGetUserTableRolesesRow{}, fmt.Errorf("count dc.user_table_roles: %w", err)
	}

	if count.TotalItems == 0 {
		return []user_domain_roles.GetUserTableRolesesRow{}, count, nil
	}

	rows, err := s.UserDomainRolesRepository.GetUserTableRoleses(ctx, params)
	if err != nil {
		return nil, user_domain_roles.CountGetUserTableRolesesRow{}, fmt.Errorf("get dc.user_table_roles page: %w", err)
	}

	return rows, count, nil
}

// GetUserTableRolesesFiltered возвращает страницу строк dc.user_table_roles и её счётчики.
func (s *UserDomainRolesService) GetUserTableRolesesFiltered(ctx context.Context, params user_domain_roles.GetUserTableRolesesFilteredParams) ([]user_domain_roles.GetUserTableRolesesFilteredRow, user_domain_roles.CountGetUserTableRolesesFilteredRow, error) {
	countParams := user_domain_roles.CountGetUserTableRolesesFilteredParams{
		PageLimit:        params.PageLimit,
		DcUserName:       params.DcUserName,
		DcUserExternalID: params.DcUserExternalID,
		TableRolesName:   params.TableRolesName,
		TableCatID:       params.TableCatID,
		TableCatName:     params.TableCatName,
	}
	count, err := s.UserDomainRolesRepository.CountGetUserTableRolesesFiltered(ctx, countParams)
	if err != nil {
		return nil, user_domain_roles.CountGetUserTableRolesesFilteredRow{}, fmt.Errorf("count dc.user_table_roles: %w", err)
	}

	if count.TotalItems == 0 {
		return []user_domain_roles.GetUserTableRolesesFilteredRow{}, count, nil
	}

	rows, err := s.UserDomainRolesRepository.GetUserTableRolesesFiltered(ctx, params)
	if err != nil {
		return nil, user_domain_roles.CountGetUserTableRolesesFilteredRow{}, fmt.Errorf("get dc.user_table_roles page: %w", err)
	}

	return rows, count, nil
}

// CreateUserTableRoles вставляет строку dc.user_table_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateUserTableRoles(ctx context.Context, params user_domain_roles.CreateUserTableRolesParams) (user_domain_roles.DcUserTableRole, error) {
	row, err := s.UserDomainRolesRepository.CreateUserTableRoles(ctx, params)

	if err != nil {
		return user_domain_roles.DcUserTableRole{}, fmt.Errorf("%w: dc.user_table_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateUserTableRolesById обновляет строку dc.user_table_roles.
//
// Запрос — :exec и не сообщает число затронутых строк, поэтому
// существование активной записи проверяется заранее.
func (s *UserDomainRolesService) UpdateUserTableRolesById(ctx context.Context, params user_domain_roles.UpdateUserTableRolesByIdParams) error {
	if _, err := s.GetUserTableRolesById(ctx, params.ID); err != nil {
		return fmt.Errorf("%w: dc.user_table_roles: %w", customerrors.ErrUpdate, err)
	}

	if err := s.UserDomainRolesRepository.UpdateUserTableRolesById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.user_table_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteUserTableRolesById мягко удаляет строку dc.user_table_roles.
//
// Сам UPDATE не фильтрует по is_deleted и не сообщает, была ли затронута
// строка, поэтому существование активной записи проверяем заранее —
// иначе удаление несуществующего id молча возвращало бы успех.
func (s *UserDomainRolesService) DeleteUserTableRolesById(ctx context.Context, params user_domain_roles.DeleteUserTableRolesByIdParams) error {
	if _, err := s.GetUserTableRolesById(ctx, params.ID); err != nil {
		return errors.Join(customerrors.ErrDelete, err)
	}

	if err := s.UserDomainRolesRepository.DeleteUserTableRolesById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.user_table_roles id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteUserTableRolesById восстанавливает мягко удалённую строку dc.user_table_roles.
func (s *UserDomainRolesService) UndeleteUserTableRolesById(ctx context.Context, params user_domain_roles.UndeleteUserTableRolesByIdParams) error {
	if err := s.UserDomainRolesRepository.UndeleteUserTableRolesById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.user_table_roles id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
