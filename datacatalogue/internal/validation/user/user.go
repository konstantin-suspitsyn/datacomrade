package user

import (
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// userWritableFields проверяет поля, общие для вставки и обновления dc.user.
func userWritableFields(
	v *validator.Validator,
	name string,
	externalId string,
) {
	v.StringVarchar("name", name, userNameMaxLen)
	v.StringUUID("external_id", externalId)
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
		req.GetExternalId(),
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
		req.GetExternalId(),
	)

	return v.Err()
}

// ValidateGetUserByExternalId проверяет запрос на выборку строки dc.user по external_id.
func ValidateGetUserByExternalId(req *userv1.GetUserByExternalIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("external_id", req.GetExternalId())

	return v.Err()
}
