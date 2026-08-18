package userdomainrolesservice

import (
	"context"
	"fmt"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	customerrors "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/utils/custom_errors"
)

// GetUserDomainRoleses возвращает страницу строк dc.user_domain_roles и её счётчики.
func (s *UserDomainRolesService) GetUserDomainRoleses(ctx context.Context, params user_domain_roles.GetUserDomainRolesesParams) ([]user_domain_roles.GetUserDomainRolesesRow, user_domain_roles.CountGetUserDomainRolesesRow, error) {
	count, err := s.UserDomainRolesRepository.CountGetUserDomainRoleses(ctx, params.PageLimit)
	if err != nil {
		return nil, user_domain_roles.CountGetUserDomainRolesesRow{}, fmt.Errorf("count dc.user_domain_roles: %w", err)
	}

	if count.TotalItems == 0 {
		return []user_domain_roles.GetUserDomainRolesesRow{}, count, nil
	}

	rows, err := s.UserDomainRolesRepository.GetUserDomainRoleses(ctx, params)
	if err != nil {
		return nil, user_domain_roles.CountGetUserDomainRolesesRow{}, fmt.Errorf("get dc.user_domain_roles page: %w", err)
	}

	return rows, count, nil
}

// CreateUserDomainRoles вставляет строку dc.user_domain_roles и возвращает её целиком.
func (s *UserDomainRolesService) CreateUserDomainRoles(ctx context.Context, params user_domain_roles.CreateUserDomainRolesParams) (user_domain_roles.DcUserDomainRole, error) {
	row, err := s.UserDomainRolesRepository.CreateUserDomainRoles(ctx, params)

	if err != nil {
		return user_domain_roles.DcUserDomainRole{}, fmt.Errorf("%w: dc.user_domain_roles: %w", customerrors.ErrCreate, err)
	}

	return row, nil
}

// UpdateUserDomainRolesById обновляет строку dc.user_domain_roles.
func (s *UserDomainRolesService) UpdateUserDomainRolesById(ctx context.Context, params user_domain_roles.UpdateUserDomainRolesByIdParams) error {
	if err := s.UserDomainRolesRepository.UpdateUserDomainRolesById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.user_domain_roles id = %d: %w", customerrors.ErrUpdate, params.ID, err)
	}

	return nil
}

// DeleteUserDomainRolesById мягко удаляет строку dc.user_domain_roles.
func (s *UserDomainRolesService) DeleteUserDomainRolesById(ctx context.Context, params user_domain_roles.DeleteUserDomainRolesByIdParams) error {
	if err := s.UserDomainRolesRepository.DeleteUserDomainRolesById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.user_domain_roles id = %d: %w", customerrors.ErrDelete, params.ID, err)
	}

	return nil
}

// UndeleteUserDomainRolesById восстанавливает мягко удалённую строку dc.user_domain_roles.
func (s *UserDomainRolesService) UndeleteUserDomainRolesById(ctx context.Context, params user_domain_roles.UndeleteUserDomainRolesByIdParams) error {
	if err := s.UserDomainRolesRepository.UndeleteUserDomainRolesById(ctx, params); err != nil {
		return fmt.Errorf("%w: dc.user_domain_roles id = %d: %w", customerrors.ErrUndelete, params.ID, err)
	}

	return nil
}
