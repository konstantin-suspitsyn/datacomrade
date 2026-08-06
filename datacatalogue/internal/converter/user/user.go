package user

import (
	"github.com/google/uuid"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_model"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
)

// UserToProto переводит строку dc.user в сущность gRPC.
func UserToProto(row user_model.DcUser) *userv1.User {
	return &userv1.User{
		Id:         row.ID,
		Name:       row.Name,
		CreatedAt:  converter.TimeToProto(row.CreatedAt),
		UpdatedAt:  converter.TimeToProto(row.UpdatedAt),
		IsDeleted:  row.IsDeleted,
		ExternalId: converter.UUIDToProto(row.ExternalID),
	}
}

// UsersToProto переводит список строк dc.user в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func UsersToProto(rows []user_model.DcUser) []*userv1.User {
	items := make([]*userv1.User, 0, len(rows))

	for _, row := range rows {
		items = append(items, UserToProto(row))
	}

	return items
}

// ToCreateUserParams собирает параметры вставки dc.user из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateUserParams(req *userv1.CreateUserRequest) user_model.CreateUserParams {
	return user_model.CreateUserParams{
		Name:       req.GetName(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToUpdateUserByIdParams собирает параметры обновления dc.user из запроса gRPC.
func ToUpdateUserByIdParams(req *userv1.UpdateUserByIdRequest) user_model.UpdateUserByIdParams {
	return user_model.UpdateUserByIdParams{
		ID:         req.GetId(),
		Name:       req.GetName(),
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
	}
}

// ToGetUserByExternalIdArg достаёт из запроса gRPC значение external_id для выборки dc.user.
func ToGetUserByExternalIdArg(req *userv1.GetUserByExternalIdRequest) uuid.UUID {
	return converter.ProtoToUUID(req.GetExternalId())
}
