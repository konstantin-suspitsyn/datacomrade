package userdomainrolesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	userdomainrolesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/userdomainroles"
	userdomainrolesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/userdomainroles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// GetUserDomainRoleses отдаёт страницу строк dc.user_domain_roles.
func (u *UserDomainRolesApiV1) GetUserDomainRoleses(ctx context.Context, req *userdomainrolesv1.GetUserDomainRolesesRequest) (*userdomainrolesv1.GetUserDomainRolesesResponse, error) {
	if err := userdomainrolesvalidation.ValidateGetUserDomainRoleses(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := userdomainrolesconv.ToGetUserDomainRolesesParams(req)

	rows, count, err := u.services.UserDomainRolesService.GetUserDomainRoleses(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserDomainRolesesResponse{
		Data: userdomainrolesconv.GetUserDomainRolesesRowsToProto(rows),
		Pagination: &userdomainrolesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateUserDomainRoles вставляет строку dc.user_domain_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateUserDomainRoles(ctx context.Context, req *userdomainrolesv1.CreateUserDomainRolesRequest) (*userdomainrolesv1.CreateUserDomainRolesResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateUserDomainRoles(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateUserDomainRoles(ctx, userdomainrolesconv.ToCreateUserDomainRolesParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateUserDomainRolesResponse{UserDomainRoles: userdomainrolesconv.UserDomainRolesToProto(row)}, nil
}

// UpdateUserDomainRolesById обновляет строку dc.user_domain_roles.
func (u *UserDomainRolesApiV1) UpdateUserDomainRolesById(ctx context.Context, req *userdomainrolesv1.UpdateUserDomainRolesByIdRequest) (*userdomainrolesv1.UpdateUserDomainRolesByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateUserDomainRolesById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UpdateUserDomainRolesById(ctx, userdomainrolesconv.ToUpdateUserDomainRolesByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateUserDomainRolesByIdResponse{}, nil
}

// DeleteUserDomainRolesById мягко удаляет строку dc.user_domain_roles.
func (u *UserDomainRolesApiV1) DeleteUserDomainRolesById(ctx context.Context, req *userdomainrolesv1.DeleteUserDomainRolesByIdRequest) (*userdomainrolesv1.DeleteUserDomainRolesByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateDeleteUserDomainRolesById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteUserDomainRolesById(ctx, userdomainrolesconv.ToDeleteUserDomainRolesByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteUserDomainRolesByIdResponse{}, nil
}

// UndeleteUserDomainRolesById восстанавливает мягко удалённую строку dc.user_domain_roles.
func (u *UserDomainRolesApiV1) UndeleteUserDomainRolesById(ctx context.Context, req *userdomainrolesv1.UndeleteUserDomainRolesByIdRequest) (*userdomainrolesv1.UndeleteUserDomainRolesByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUndeleteUserDomainRolesById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteUserDomainRolesById(ctx, userdomainrolesconv.ToUndeleteUserDomainRolesByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteUserDomainRolesByIdResponse{}, nil
}
