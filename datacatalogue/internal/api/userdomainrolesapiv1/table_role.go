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

// GetTableRoleById отдаёт активную строку dc.table_roles по id.
func (u *UserDomainRolesApiV1) GetTableRoleById(ctx context.Context, req *userdomainrolesv1.GetTableRoleByIdRequest) (*userdomainrolesv1.GetTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetTableRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetTableRoleByIdResponse{TableRole: userdomainrolesconv.TableRoleToProto(row)}, nil
}

// GetDeletedTableRoleById отдаёт мягко удалённую строку dc.table_roles по id.
func (u *UserDomainRolesApiV1) GetDeletedTableRoleById(ctx context.Context, req *userdomainrolesv1.GetDeletedTableRoleByIdRequest) (*userdomainrolesv1.GetDeletedTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDeletedTableRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedTableRoleByIdResponse{TableRole: userdomainrolesconv.TableRoleToProto(row)}, nil
}

// GetTableRoles отдаёт строки dc.table_roles.
func (u *UserDomainRolesApiV1) GetTableRoles(ctx context.Context, req *userdomainrolesv1.GetTableRolesRequest) (*userdomainrolesv1.GetTableRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetTableRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetTableRolesResponse{TableRoles: userdomainrolesconv.TableRolesToProto(rows)}, nil
}

// GetDeletedTableRoles отдаёт строки dc.table_roles.
func (u *UserDomainRolesApiV1) GetDeletedTableRoles(ctx context.Context, req *userdomainrolesv1.GetDeletedTableRolesRequest) (*userdomainrolesv1.GetDeletedTableRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDeletedTableRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedTableRolesResponse{TableRoles: userdomainrolesconv.TableRolesToProto(rows)}, nil
}

// CreateTableRole вставляет строку dc.table_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateTableRole(ctx context.Context, req *userdomainrolesv1.CreateTableRoleRequest) (*userdomainrolesv1.CreateTableRoleResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateTableRole(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateTableRole(ctx, userdomainrolesconv.ToCreateTableRoleParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateTableRoleResponse{TableRole: userdomainrolesconv.TableRoleToProto(row)}, nil
}

// UpdateTableRoleById обновляет строку dc.table_roles.
func (u *UserDomainRolesApiV1) UpdateTableRoleById(ctx context.Context, req *userdomainrolesv1.UpdateTableRoleByIdRequest) (*userdomainrolesv1.UpdateTableRoleByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateTableRoleById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.UpdateTableRoleById(ctx, userdomainrolesconv.ToUpdateTableRoleByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateTableRoleByIdResponse{TableRole: userdomainrolesconv.TableRoleToProto(row)}, nil
}

// DeleteTableRoleById мягко удаляет строку dc.table_roles.
func (u *UserDomainRolesApiV1) DeleteTableRoleById(ctx context.Context, req *userdomainrolesv1.DeleteTableRoleByIdRequest) (*userdomainrolesv1.DeleteTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteTableRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteTableRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteTableRoleById восстанавливает мягко удалённую строку dc.table_roles.
func (u *UserDomainRolesApiV1) UndeleteTableRoleById(ctx context.Context, req *userdomainrolesv1.UndeleteTableRoleByIdRequest) (*userdomainrolesv1.UndeleteTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteTableRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteTableRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}
