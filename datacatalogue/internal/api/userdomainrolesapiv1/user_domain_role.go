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

// GetUserDomainRoleById отдаёт активную строку dc.user_domain_roles по id.
func (u *UserDomainRolesApiV1) GetUserDomainRoleById(ctx context.Context, req *userdomainrolesv1.GetUserDomainRoleByIdRequest) (*userdomainrolesv1.GetUserDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetUserDomainRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserDomainRoleByIdResponse{UserDomainRole: userdomainrolesconv.UserDomainRoleToProto(row)}, nil
}

// GetUserDomainRoles отдаёт все активные строки dc.user_domain_roles.
func (u *UserDomainRolesApiV1) GetUserDomainRoles(ctx context.Context, req *userdomainrolesv1.GetUserDomainRolesRequest) (*userdomainrolesv1.GetUserDomainRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetUserDomainRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetUserDomainRolesResponse{UserDomainRoles: userdomainrolesconv.UserDomainRolesToProto(rows)}, nil
}

// GetDeletedUserDomainRoleById отдаёт мягко удалённую строку dc.user_domain_roles по id.
func (u *UserDomainRolesApiV1) GetDeletedUserDomainRoleById(ctx context.Context, req *userdomainrolesv1.GetDeletedUserDomainRoleByIdRequest) (*userdomainrolesv1.GetDeletedUserDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDeletedUserDomainRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedUserDomainRoleByIdResponse{UserDomainRole: userdomainrolesconv.UserDomainRoleToProto(row)}, nil
}

// GetDeletedUserDomainRoles отдаёт все мягко удалённые строки dc.user_domain_roles.
func (u *UserDomainRolesApiV1) GetDeletedUserDomainRoles(ctx context.Context, req *userdomainrolesv1.GetDeletedUserDomainRolesRequest) (*userdomainrolesv1.GetDeletedUserDomainRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDeletedUserDomainRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedUserDomainRolesResponse{UserDomainRoles: userdomainrolesconv.UserDomainRolesToProto(rows)}, nil
}

// CreateUserDomainRole вставляет строку dc.user_domain_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateUserDomainRole(ctx context.Context, req *userdomainrolesv1.CreateUserDomainRoleRequest) (*userdomainrolesv1.CreateUserDomainRoleResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateUserDomainRole(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateUserDomainRole(ctx, userdomainrolesconv.ToCreateUserDomainRoleParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateUserDomainRoleResponse{UserDomainRole: userdomainrolesconv.UserDomainRoleToProto(row)}, nil
}

// UpdateUserDomainRoleById обновляет активную строку dc.user_domain_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) UpdateUserDomainRoleById(ctx context.Context, req *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) (*userdomainrolesv1.UpdateUserDomainRoleByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateUserDomainRoleById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.UpdateUserDomainRoleById(ctx, userdomainrolesconv.ToUpdateUserDomainRoleByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateUserDomainRoleByIdResponse{UserDomainRole: userdomainrolesconv.UserDomainRoleToProto(row)}, nil
}

// DeleteUserDomainRoleById мягко удаляет строку dc.user_domain_roles.
func (u *UserDomainRolesApiV1) DeleteUserDomainRoleById(ctx context.Context, req *userdomainrolesv1.DeleteUserDomainRoleByIdRequest) (*userdomainrolesv1.DeleteUserDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteUserDomainRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteUserDomainRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteUserDomainRoleById восстанавливает мягко удалённую строку dc.user_domain_roles.
func (u *UserDomainRolesApiV1) UndeleteUserDomainRoleById(ctx context.Context, req *userdomainrolesv1.UndeleteUserDomainRoleByIdRequest) (*userdomainrolesv1.UndeleteUserDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteUserDomainRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteUserDomainRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}
