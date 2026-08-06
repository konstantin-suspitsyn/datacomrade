package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateUser проверяет запрос на вставку строки dc.user.
func ValidateCreateUser(req *tablesv1.CreateUserRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), userNameMaxLen)
	v.StringUUID("external_id", req.GetExternalId())

	return v.Err()
}

// ValidateGetUserByExternalId проверяет запрос на выборку строки dc.user по external_id.
func ValidateGetUserByExternalId(req *tablesv1.GetUserByExternalIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("external_id", req.GetExternalId())

	return v.Err()
}

// ValidateGetUsers проверяет запрос страницы dc.user.
func ValidateGetUsers(req *tablesv1.GetUsersRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int32Between("page_limit", req.GetPageLimit(), 0, validation.MaxPageSize)
	v.Int32Min("page", req.GetPage(), 0)
	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}
