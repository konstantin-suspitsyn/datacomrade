package authlogic

import (
	authlogicv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/auth_logic/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateGetTableIdsByExternalUserIdAndRoles проверяет запрос на выборку auth_logic.GetTableIdsByExternalUserIdAndRoles.
func ValidateGetTableIdsByExternalUserIdAndRoles(req *authlogicv1.GetTableIdsByExternalUserIdAndRolesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("external_id", req.GetExternalId())
	v.StringVarchar("name", req.GetName(), nameMaxLen)

	return v.Err()
}

// ValidateGetTableIdsByUserIdAndRoles проверяет запрос на выборку auth_logic.GetTableIdsByUserIdAndRoles.
func ValidateGetTableIdsByUserIdAndRoles(req *authlogicv1.GetTableIdsByUserIdAndRolesRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("user_id", req.GetUserId())
	v.StringVarchar("name", req.GetName(), nameMaxLen)

	return v.Err()
}
