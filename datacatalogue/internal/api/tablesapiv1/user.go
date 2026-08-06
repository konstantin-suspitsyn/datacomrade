package tablesapiv1

import (
	"context"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/api/apierror"
	tablesconv "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter/tables"
	tablesvalidation "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation/tables"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// GetUsers отдаёт страницу строк dc.user.
func (t *TablesApiV1) GetUsers(ctx context.Context, req *tablesv1.GetUsersRequest) (*tablesv1.GetUsersResponse, error) {
	if err := tablesvalidation.ValidateGetUsers(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	params := tablesconv.ToGetUsersParams(req)

	rows, count, err := t.services.TablesService.GetUsers(ctx, params)
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetUsersResponse{
		Data: tablesconv.UsersToProto(rows),
		Pagination: &tablesv1.Pagination{
			Page:       params.Page,
			PerPage:    params.PageLimit,
			TotalItems: count.TotalItems,
			TotalPages: count.TotalPages,
		},
	}, nil
}

// GetUserByExternalId отдаёт активную строку dc.user по уникальной колонке external_id.
func (t *TablesApiV1) GetUserByExternalId(ctx context.Context, req *tablesv1.GetUserByExternalIdRequest) (*tablesv1.GetUserByExternalIdResponse, error) {
	if err := tablesvalidation.ValidateGetUserByExternalId(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.GetUserByExternalId(ctx, tablesconv.ToGetUserByExternalIdArg(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.GetUserByExternalIdResponse{User: tablesconv.UserToProto(row)}, nil
}

// CreateUser вставляет строку dc.user и отдаёт её целиком.
func (t *TablesApiV1) CreateUser(ctx context.Context, req *tablesv1.CreateUserRequest) (*tablesv1.CreateUserResponse, error) {
	if err := tablesvalidation.ValidateCreateUser(req); err != nil {
		return nil, apierror.Wrap(err)
	}

	row, err := t.services.TablesService.CreateUser(ctx, tablesconv.ToCreateUserParams(req))
	if err != nil {
		return nil, apierror.Wrap(err)
	}

	return &tablesv1.CreateUserResponse{User: tablesconv.UserToProto(row)}, nil
}
