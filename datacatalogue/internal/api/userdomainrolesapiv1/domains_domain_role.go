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

// GetDomainsDomainRoleById отдаёт активную строку dc.domains_domain_roles по id.
func (u *UserDomainRolesApiV1) GetDomainsDomainRoleById(ctx context.Context, req *userdomainrolesv1.GetDomainsDomainRoleByIdRequest) (*userdomainrolesv1.GetDomainsDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDomainsDomainRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDomainsDomainRoleByIdResponse{DomainsDomainRole: userdomainrolesconv.DomainsDomainRoleToProto(row)}, nil
}

// GetDomainsDomainRoles отдаёт все активные строки dc.domains_domain_roles.
func (u *UserDomainRolesApiV1) GetDomainsDomainRoles(ctx context.Context, req *userdomainrolesv1.GetDomainsDomainRolesRequest) (*userdomainrolesv1.GetDomainsDomainRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDomainsDomainRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDomainsDomainRolesResponse{DomainsDomainRoles: userdomainrolesconv.DomainsDomainRolesToProto(rows)}, nil
}

// GetDeletedDomainsDomainRoleById отдаёт мягко удалённую строку dc.domains_domain_roles по id.
func (u *UserDomainRolesApiV1) GetDeletedDomainsDomainRoleById(ctx context.Context, req *userdomainrolesv1.GetDeletedDomainsDomainRoleByIdRequest) (*userdomainrolesv1.GetDeletedDomainsDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDeletedDomainsDomainRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedDomainsDomainRoleByIdResponse{DomainsDomainRole: userdomainrolesconv.DomainsDomainRoleToProto(row)}, nil
}

// GetDeletedDomainsDomainRoles отдаёт все мягко удалённые строки dc.domains_domain_roles.
func (u *UserDomainRolesApiV1) GetDeletedDomainsDomainRoles(ctx context.Context, req *userdomainrolesv1.GetDeletedDomainsDomainRolesRequest) (*userdomainrolesv1.GetDeletedDomainsDomainRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDeletedDomainsDomainRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedDomainsDomainRolesResponse{DomainsDomainRoles: userdomainrolesconv.DomainsDomainRolesToProto(rows)}, nil
}

// CreateDomainsDomainRole вставляет строку dc.domains_domain_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateDomainsDomainRole(ctx context.Context, req *userdomainrolesv1.CreateDomainsDomainRoleRequest) (*userdomainrolesv1.CreateDomainsDomainRoleResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateDomainsDomainRole(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateDomainsDomainRole(ctx, userdomainrolesconv.ToCreateDomainsDomainRoleParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateDomainsDomainRoleResponse{DomainsDomainRole: userdomainrolesconv.DomainsDomainRoleToProto(row)}, nil
}

// UpdateDomainsDomainRoleById обновляет активную строку dc.domains_domain_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) UpdateDomainsDomainRoleById(ctx context.Context, req *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) (*userdomainrolesv1.UpdateDomainsDomainRoleByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateDomainsDomainRoleById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.UpdateDomainsDomainRoleById(ctx, userdomainrolesconv.ToUpdateDomainsDomainRoleByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateDomainsDomainRoleByIdResponse{DomainsDomainRole: userdomainrolesconv.DomainsDomainRoleToProto(row)}, nil
}

// DeleteDomainsDomainRoleById мягко удаляет строку dc.domains_domain_roles.
func (u *UserDomainRolesApiV1) DeleteDomainsDomainRoleById(ctx context.Context, req *userdomainrolesv1.DeleteDomainsDomainRoleByIdRequest) (*userdomainrolesv1.DeleteDomainsDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteDomainsDomainRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteDomainsDomainRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteDomainsDomainRoleById восстанавливает мягко удалённую строку dc.domains_domain_roles.
func (u *UserDomainRolesApiV1) UndeleteDomainsDomainRoleById(ctx context.Context, req *userdomainrolesv1.UndeleteDomainsDomainRoleByIdRequest) (*userdomainrolesv1.UndeleteDomainsDomainRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteDomainsDomainRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteDomainsDomainRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}
