package tables

import (
	"github.com/google/uuid"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// UserToProto переводит строку dc.user в сущность gRPC.
func UserToProto(row tables_model.DcUser) *tablesv1.User {
	return &tablesv1.User{
		Id:         converter.Int64ToProto(row.ID),
		Name:       row.Name,
		CreatedAt:  converter.TimeToProto(row.CreatedAt),
		UpdatedAt:  converter.TimeToProto(row.UpdatedAt),
		IsDeleted:  row.IsDeleted,
		ExternalId: converter.UUIDToProto(row.ExternalID),
	}
}

// UsersToProto переводит список строк dc.user в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func UsersToProto(rows []tables_model.DcUser) []*tablesv1.User {
	items := make([]*tablesv1.User, 0, len(rows))

	for _, row := range rows {
		items = append(items, UserToProto(row))
	}

	return items
}

// ToCreateUserParams собирает параметры вставки dc.user из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateUserParams(req *tablesv1.CreateUserRequest) tables_model.CreateUserParams {
	return tables_model.CreateUserParams{
		Name:       req.GetName(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToGetUserByExternalIdArg достаёт из запроса gRPC значение external_id для выборки dc.user.
func ToGetUserByExternalIdArg(req *tablesv1.GetUserByExternalIdRequest) uuid.UUID {
	return converter.ProtoToUUID(req.GetExternalId())
}

// ToGetUsersParams собирает параметры страницы dc.user из запроса gRPC.
func ToGetUsersParams(req *tablesv1.GetUsersRequest) tables_model.GetUsersParams {
	limit := req.GetPageLimit()
	if limit == 0 {
		limit = validation.DefaultPageSize
	}

	page := req.GetPage()
	if page == 0 {
		page = 1
	}

	return tables_model.GetUsersParams{
		Order:     req.GetOrder(),
		Page:      page,
		PageLimit: limit,
	}
}
