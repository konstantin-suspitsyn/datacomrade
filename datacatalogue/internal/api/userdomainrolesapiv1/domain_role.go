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

// GetDomainRoleById отдаёт активную строку dc.domain_roles по id.
func (u *UserDomainRolesApiV1) GetDomainRoleById(ctx context.Context, req *userdomainrolesv1.GetDomainRoleByIdRequest) (*userdomainrolesv1.GetDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDomainRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDomainRoleByIdResponse{DomainRole: userdomainrolesconv.DomainRoleToProto(row)}, nil
}

// GetDeletedDomainRoleById отдаёт мягко удалённую строку dc.domain_roles по id.
func (u *UserDomainRolesApiV1) GetDeletedDomainRoleById(ctx context.Context, req *userdomainrolesv1.GetDeletedDomainRoleByIdRequest) (*userdomainrolesv1.GetDeletedDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDeletedDomainRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedDomainRoleByIdResponse{DomainRole: userdomainrolesconv.DomainRoleToProto(row)}, nil
}

// GetDomainRoles отдаёт строки dc.domain_roles.
func (u *UserDomainRolesApiV1) GetDomainRoles(ctx context.Context, req *userdomainrolesv1.GetDomainRolesRequest) (*userdomainrolesv1.GetDomainRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDomainRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDomainRolesResponse{DomainRoles: userdomainrolesconv.DomainRolesToProto(rows)}, nil
}

// GetDeletedDomainRoles отдаёт строки dc.domain_roles.
func (u *UserDomainRolesApiV1) GetDeletedDomainRoles(ctx context.Context, req *userdomainrolesv1.GetDeletedDomainRolesRequest) (*userdomainrolesv1.GetDeletedDomainRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDeletedDomainRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedDomainRolesResponse{DomainRoles: userdomainrolesconv.DomainRolesToProto(rows)}, nil
}

// CreateDomainRole вставляет строку dc.domain_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateDomainRole(ctx context.Context, req *userdomainrolesv1.CreateDomainRoleRequest) (*userdomainrolesv1.CreateDomainRoleResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateDomainRole(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateDomainRole(ctx, userdomainrolesconv.ToCreateDomainRoleParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateDomainRoleResponse{DomainRole: userdomainrolesconv.DomainRoleToProto(row)}, nil
}

// UpdateDomainRoleById обновляет строку dc.domain_roles.
func (u *UserDomainRolesApiV1) UpdateDomainRoleById(ctx context.Context, req *userdomainrolesv1.UpdateDomainRoleByIdRequest) (*userdomainrolesv1.UpdateDomainRoleByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateDomainRoleById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.UpdateDomainRoleById(ctx, userdomainrolesconv.ToUpdateDomainRoleByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateDomainRoleByIdResponse{DomainRole: userdomainrolesconv.DomainRoleToProto(row)}, nil
}

// DeleteDomainRoleById мягко удаляет строку dc.domain_roles.
func (u *UserDomainRolesApiV1) DeleteDomainRoleById(ctx context.Context, req *userdomainrolesv1.DeleteDomainRoleByIdRequest) (*userdomainrolesv1.DeleteDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteDomainRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteDomainRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteDomainRoleById восстанавливает мягко удалённую строку dc.domain_roles.
func (u *UserDomainRolesApiV1) UndeleteDomainRoleById(ctx context.Context, req *userdomainrolesv1.UndeleteDomainRoleByIdRequest) (*userdomainrolesv1.UndeleteDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteDomainRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteDomainRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}
