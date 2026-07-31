package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// hasToGroupWritableFields проверяет поля, общие для вставки и обновления dc.has_to_group.
func hasToGroupWritableFields(
	v *validator.Validator,
	columnIdA int64,
	columnIdB int64,
	description string,
	userExternalId string,
) {
	v.Int64ID("column_id_a", columnIdA)
	v.Int64ID("column_id_b", columnIdB)
	v.StringVarchar("description", description, hasToGroupDescriptionMaxLen)
	v.StringUUID("user_id", userExternalId)
}

// ValidateCreateHasToGroup проверяет запрос на вставку строки dc.has_to_group.
func ValidateCreateHasToGroup(req *tablesv1.CreateHasToGroupRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	hasToGroupWritableFields(
		v,
		req.GetColumnIdA(),
		req.GetColumnIdB(),
		req.GetDescription(),
		req.GetUserExternalId(),
	)

	return v.Err()
}

// ValidateUpdateHasToGroupById проверяет запрос на обновление строки dc.has_to_group.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateHasToGroupById(req *tablesv1.UpdateHasToGroupByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	hasToGroupWritableFields(
		v,
		req.GetColumnIdA(),
		req.GetColumnIdB(),
		req.GetDescription(),
		req.GetUserExternalId(),
	)

	return v.Err()
}
