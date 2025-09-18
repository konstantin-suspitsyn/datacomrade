package roles

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/konstantin-suspitsyn/datacomrade/data/paginationmodel"
	"github.com/konstantin-suspitsyn/datacomrade/data/rolesmodel"
	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/urlparams"
)

func (rs *RoleService) getDataWithPager(r *http.Request) (*rolesmodel.RolesWithPager, error) {
	ctx := r.Context()

	urlPagingParams, err := urlparams.GetPager(r)
	if err != nil {
		return nil, err
	}

	paginator, err := rs.generatePaginator(ctx, urlPagingParams)

	if err != nil {
		return nil, err
	}

	argsForPager := rolesmodel.GetRolesWithPagerParams{
		Limit:  paginator.GetLimit(),
		Offset: paginator.GetOffset(),
	}

	rows, err := rs.Models.GetRolesWithPager(ctx, argsForPager)

	if err != nil {
		return nil, err
	}

	dto := rolesmodel.RolesWithPager{
		Data:   rows,
		Paging: paginator,
	}

	return &dto, nil
}

func (rs *RoleService) generatePaginator(ctx context.Context, pager *urlparams.Pager) (*paginationmodel.Pagination, error) {

	totalItems, err := rs.Models.CountActiveRows(ctx)

	if err != nil {
		return nil, err
	}

	paginator, err := paginationmodel.New(totalItems, pager.PageSize, pager.CurrentPage, configs.DOMAIN_LINK)

	if err != nil {
		return nil, err
	}

	return paginator, nil
}

func (rd *RoleService) CreateRole(ctx context.Context, shortName string, longName string, description string, jwtExport bool, userId int64) (rolesmodel.UsersRole, error) {

	rolesDescription := sql.NullString{
		String: description,
		Valid:  true,
	}

	createRoleParams := rolesmodel.CreateRoleParams{
		RoleNameLong:  longName,
		RoleNameShort: shortName,
		Description:   rolesDescription,
		JwtExport:     jwtExport,
		UserID:        userId,
	}

	return rd.Models.CreateRole(ctx, createRoleParams)

}

func (rs *RoleService) DeleteRole(ctx context.Context, roleId int64, userId int64) error {

	deleteParams := rolesmodel.DeleteRoleParams{
		UserID: userId,
		ID:     roleId,
	}
	return rs.Models.DeleteRole(ctx, deleteParams)
}

func (rs *RoleService) UpdateRole(ctx context.Context, roleId int64, roleNameShort string, roleNameLong string, jwtExport bool, description string, userId int64) (rolesmodel.UsersRole, error) {
	rolesDescription := sql.NullString{
		String: description,
		Valid:  true,
	}

	roleParams := rolesmodel.UpdateRoleParams{
		Description:   rolesDescription,
		ID:            roleId,
		RoleNameLong:  roleNameLong,
		UserID:        userId,
		RoleNameShort: roleNameShort,
		JwtExport:     jwtExport,
	}
	userRole, err := rs.Models.UpdateRole(ctx, roleParams)

	return userRole, err
}
