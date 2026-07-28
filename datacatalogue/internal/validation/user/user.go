package user

import (
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userWritableFields проверяет поля, общие для вставки и обновления dc.user.
func userWritableFields(
	v *validator.Validator,
	name string,
	incomingUserId int64,
) {
	v.StringVarchar("name", name, userNameMaxLen)
	v.Int64ID("incoming_user_id", incomingUserId)
}

// ValidateCreateUser проверяет запрос на вставку строки dc.user.
func ValidateCreateUser(req *userv1.CreateUserRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	userWritableFields(
		v,
		req.GetName(),
		req.GetIncomingUserId(),
	)

	return v.Err()
}

// ValidateUpdateUserById проверяет запрос на обновление строки dc.user.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateUserById(req *userv1.UpdateUserByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	userWritableFields(
		v,
		req.GetName(),
		req.GetIncomingUserId(),
	)

	return v.Err()
}
