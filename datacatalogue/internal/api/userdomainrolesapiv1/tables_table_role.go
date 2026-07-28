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

// GetTablesTableRoleById отдаёт активную строку dc.tables_table_roles по id.
func (u *UserDomainRolesApiV1) GetTablesTableRoleById(ctx context.Context, req *userdomainrolesv1.GetTablesTableRoleByIdRequest) (*userdomainrolesv1.GetTablesTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetTablesTableRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetTablesTableRoleByIdResponse{TablesTableRole: userdomainrolesconv.TablesTableRoleToProto(row)}, nil
}

// GetTablesTableRoles отдаёт все активные строки dc.tables_table_roles.
func (u *UserDomainRolesApiV1) GetTablesTableRoles(ctx context.Context, req *userdomainrolesv1.GetTablesTableRolesRequest) (*userdomainrolesv1.GetTablesTableRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetTablesTableRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetTablesTableRolesResponse{TablesTableRoles: userdomainrolesconv.TablesTableRolesToProto(rows)}, nil
}

// GetDeletedTablesTableRoleById отдаёт мягко удалённую строку dc.tables_table_roles по id.
func (u *UserDomainRolesApiV1) GetDeletedTablesTableRoleById(ctx context.Context, req *userdomainrolesv1.GetDeletedTablesTableRoleByIdRequest) (*userdomainrolesv1.GetDeletedTablesTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.GetDeletedTablesTableRoleById(ctx, req.GetId())
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedTablesTableRoleByIdResponse{TablesTableRole: userdomainrolesconv.TablesTableRoleToProto(row)}, nil
}

// GetDeletedTablesTableRoles отдаёт все мягко удалённые строки dc.tables_table_roles.
func (u *UserDomainRolesApiV1) GetDeletedTablesTableRoles(ctx context.Context, req *userdomainrolesv1.GetDeletedTablesTableRolesRequest) (*userdomainrolesv1.GetDeletedTablesTableRolesResponse, error) {
	rows, err := u.services.UserDomainRolesService.GetDeletedTablesTableRoles(ctx)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.GetDeletedTablesTableRolesResponse{TablesTableRoles: userdomainrolesconv.TablesTableRolesToProto(rows)}, nil
}

// CreateTablesTableRole вставляет строку dc.tables_table_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) CreateTablesTableRole(ctx context.Context, req *userdomainrolesv1.CreateTablesTableRoleRequest) (*userdomainrolesv1.CreateTablesTableRoleResponse, error) {
	if err := userdomainrolesvalidation.ValidateCreateTablesTableRole(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.CreateTablesTableRole(ctx, userdomainrolesconv.ToCreateTablesTableRoleParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.CreateTablesTableRoleResponse{TablesTableRole: userdomainrolesconv.TablesTableRoleToProto(row)}, nil
}

// UpdateTablesTableRoleById обновляет активную строку dc.tables_table_roles и отдаёт её целиком.
func (u *UserDomainRolesApiV1) UpdateTablesTableRoleById(ctx context.Context, req *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) (*userdomainrolesv1.UpdateTablesTableRoleByIdResponse, error) {
	if err := userdomainrolesvalidation.ValidateUpdateTablesTableRoleById(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := u.services.UserDomainRolesService.UpdateTablesTableRoleById(ctx, userdomainrolesconv.ToUpdateTablesTableRoleByIdParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UpdateTablesTableRoleByIdResponse{TablesTableRole: userdomainrolesconv.TablesTableRoleToProto(row)}, nil
}

// DeleteTablesTableRoleById мягко удаляет строку dc.tables_table_roles.
func (u *UserDomainRolesApiV1) DeleteTablesTableRoleById(ctx context.Context, req *userdomainrolesv1.DeleteTablesTableRoleByIdRequest) (*userdomainrolesv1.DeleteTablesTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.DeleteTablesTableRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.DeleteTablesTableRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}

// UndeleteTablesTableRoleById восстанавливает мягко удалённую строку dc.tables_table_roles.
func (u *UserDomainRolesApiV1) UndeleteTablesTableRoleById(ctx context.Context, req *userdomainrolesv1.UndeleteTablesTableRoleByIdRequest) (*userdomainrolesv1.UndeleteTablesTableRoleByIdResponse, error) {
	if err := validation.ValidateID(req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	if err := u.services.UserDomainRolesService.UndeleteTablesTableRoleById(ctx, req.GetId()); err != nil {
		return nil, apierror.Wrap(err)
	}

	return &userdomainrolesv1.UndeleteTablesTableRoleByIdResponse{Empty: &emptypb.Empty{}}, nil
}
