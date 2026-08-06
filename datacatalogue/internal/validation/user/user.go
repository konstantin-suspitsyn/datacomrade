package user

import (
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateUser проверяет запрос на вставку строки dc.user.
func ValidateCreateUser(req *userv1.CreateUserRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), userNameMaxLen)
	v.StringUUID("external_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateUserById проверяет запрос на обновление строки dc.user.
func ValidateUpdateUserById(req *userv1.UpdateUserByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())
	v.StringVarchar("name", req.GetName(), userNameMaxLen)
	v.StringUUID("external_id", req.GetExternalId())

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
