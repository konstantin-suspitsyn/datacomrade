package userdomainrolesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	userdomainrolesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/userdomainroles"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/userdomainroles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// GetUserTableRolesById отдаёт активную строку dc.user_table_roles по id.
func (u *UserDomainRolesApiV1) GetUserTableRolesById(ctx context.Context, req *userdomainrolesv1.GetUserTableRolesByIdRequest) (*userdomainrolesv1.GetUserTableRolesByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetUserTableRolesById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserTableRolesByIdResponse{UserTableRoles: userdomainrolesconv.GetUserTableRolesByIdRowToProto(row)}, nil
}

// GetUserTableRoleses отдаёт страницу строк dc.user_table_roles.
func (u *UserDomainRolesApiV1) GetUserTableRoleses(ctx context.Context, req *userdomainrolesv1.GetUserTableRolesesRequest) (*userdomainrolesv1.GetUserTableRolesesResponse, error) {
	if err := userdomainrolesvalidation.ValidateGetUserTableRoleses(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := userdomainrolesconv.ToGetUserTableRolesesParams(req)

	rows, count, err := u.services.UserDomainRolesService.GetUserTableRoleses(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserTableRolesesResponse{
		Data: userdomainrolesconv.GetUserTableRolesesRowsToProto(rows),
		Pagination: &userdomainrolesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetUserTableRolesesFiltered отдаёт страницу строк dc.user_table_roles.
func (u *UserDomainRolesApiV1) GetUserTableRolesesFiltered(ctx context.Context, req *userdomainrolesv1.GetUserTableRolesesFilteredRequest) (*userdomainrolesv1.GetUserTableRolesesFilteredResponse, error) {
	if err := userdomainrolesvalidation.ValidateGetUserTableRolesesFiltered(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := userdomainrolesconv.ToGetUserTableRolesesFilteredParams(req)

	rows, count, err := u.services.UserDomainRolesService.GetUserTableRolesesFiltered(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserTableRolesesFilteredResponse{
		Data: userdomainrolesconv.GetUserTableRolesesFilteredRowsToProto(rows),
		Pagination: &userdomainrolesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// CreateUserTableRoles вставляет строку dc.user_table_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateUserTableRoles(ctx context.Context, req *userdomainrolesv1.CreateUserTableRolesRequest) (*userdomainrolesv1.CreateUserTableRolesResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateUserTableRoles(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateUserTableRoles(ctx, userdomainrolesconv.ToCreateUserTableRolesParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateUserTableRolesResponse{UserTableRoles: userdomainrolesconv.UserTableRolesToProto(row)}, nil
}

// UpdateUserTableRolesById обновляет строку dc.user_table_roles.
func (u *UserDomainRolesApiV1) UpdateUserTableRolesById(ctx context.Context, req *userdomainrolesv1.UpdateUserTableRolesByIdRequest) (*userdomainrolesv1.UpdateUserTableRolesByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateUserTableRolesById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UpdateUserTableRolesById(ctx, userdomainrolesconv.ToUpdateUserTableRolesByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateUserTableRolesByIdResponse{}, nil
}

// DeleteUserTableRolesById мягко удаляет строку dc.user_table_roles.
func (u *UserDomainRolesApiV1) DeleteUserTableRolesById(ctx context.Context, req *userdomainrolesv1.DeleteUserTableRolesByIdRequest) (*userdomainrolesv1.DeleteUserTableRolesByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateDeleteUserTableRolesById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteUserTableRolesById(ctx, userdomainrolesconv.ToDeleteUserTableRolesByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteUserTableRolesByIdResponse{}, nil
}

// UndeleteUserTableRolesById восстанавливает мягко удалённую строку dc.user_table_roles.
func (u *UserDomainRolesApiV1) UndeleteUserTableRolesById(ctx context.Context, req *userdomainrolesv1.UndeleteUserTableRolesByIdRequest) (*userdomainrolesv1.UndeleteUserTableRolesByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUndeleteUserTableRolesById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteUserTableRolesById(ctx, userdomainrolesconv.ToUndeleteUserTableRolesByIdParams(req)); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteUserTableRolesByIdResponse{}, nil
}
