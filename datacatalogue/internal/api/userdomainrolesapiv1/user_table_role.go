package userdomainrolesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	userdomainrolesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/userdomainroles"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	userdomainrolesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/userdomainroles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetUserTableRoleById отдаёт активную строку dc.user_table_roles по id.
func (u *UserDomainRolesApiV1) GetUserTableRoleById(ctx context.Context, req *userdomainrolesv1.GetUserTableRoleByIdRequest) (*userdomainrolesv1.GetUserTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetUserTableRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserTableRoleByIdResponse{UserTableRole: userdomainrolesconv.UserTableRoleToProto(row)}, nil
}

// GetUserTableRoles отдаёт все активные строки dc.user_table_roles.
func (u *UserDomainRolesApiV1) GetUserTableRoles(ctx context.Context, req *userdomainrolesv1.GetUserTableRolesRequest) (*userdomainrolesv1.GetUserTableRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetUserTableRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserTableRolesResponse{UserTableRoles: userdomainrolesconv.UserTableRolesToProto(rows)}, nil
}

// GetDeletedUserTableRoleById отдаёт мягко удалённую строку dc.user_table_roles по id.
func (u *UserDomainRolesApiV1) GetDeletedUserTableRoleById(ctx context.Context, req *userdomainrolesv1.GetDeletedUserTableRoleByIdRequest) (*userdomainrolesv1.GetDeletedUserTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDeletedUserTableRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedUserTableRoleByIdResponse{UserTableRole: userdomainrolesconv.UserTableRoleToProto(row)}, nil
}

// GetDeletedUserTableRoles отдаёт все мягко удалённые строки dc.user_table_roles.
func (u *UserDomainRolesApiV1) GetDeletedUserTableRoles(ctx context.Context, req *userdomainrolesv1.GetDeletedUserTableRolesRequest) (*userdomainrolesv1.GetDeletedUserTableRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDeletedUserTableRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedUserTableRolesResponse{UserTableRoles: userdomainrolesconv.UserTableRolesToProto(rows)}, nil
}

// CreateUserTableRole вставляет строку dc.user_table_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateUserTableRole(ctx context.Context, req *userdomainrolesv1.CreateUserTableRoleRequest) (*userdomainrolesv1.CreateUserTableRoleResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateUserTableRole(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateUserTableRole(ctx, userdomainrolesconv.ToCreateUserTableRoleParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateUserTableRoleResponse{UserTableRole: userdomainrolesconv.UserTableRoleToProto(row)}, nil
}

// UpdateUserTableRoleById обновляет активную строку dc.user_table_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) UpdateUserTableRoleById(ctx context.Context, req *userdomainrolesv1.UpdateUserTableRoleByIdRequest) (*userdomainrolesv1.UpdateUserTableRoleByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateUserTableRoleById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.UpdateUserTableRoleById(ctx, userdomainrolesconv.ToUpdateUserTableRoleByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateUserTableRoleByIdResponse{UserTableRole: userdomainrolesconv.UserTableRoleToProto(row)}, nil
}

// DeleteUserTableRoleById мягко удаляет строку dc.user_table_roles.
func (u *UserDomainRolesApiV1) DeleteUserTableRoleById(ctx context.Context, req *userdomainrolesv1.DeleteUserTableRoleByIdRequest) (*userdomainrolesv1.DeleteUserTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteUserTableRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteUserTableRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteUserTableRoleById восстанавливает мягко удалённую строку dc.user_table_roles.
func (u *UserDomainRolesApiV1) UndeleteUserTableRoleById(ctx context.Context, req *userdomainrolesv1.UndeleteUserTableRoleByIdRequest) (*userdomainrolesv1.UndeleteUserTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteUserTableRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteUserTableRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}
