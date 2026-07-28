package tables

import (
	"math"

	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// groupLevelWritableFields проверяет поля, общие для вставки и обновления dc.group_levels.
func groupLevelWritableFields(
	v *validator.Validator,
	columnId int64,
	parentColumnId int64,
	level int32,
	description string,
	userId int64,
) {
	v.Int64ID("column_id", columnId)
	v.Int64ID("parent_column_id", parentColumnId)
	v.Int32Between("level", level, math.MinInt16, math.MaxInt16)
	v.StringVarchar("description", description, groupLevelDescriptionMaxLen)
	v.Int64ID("user_id", userId)
}

// ValidateCreateGroupLevel проверяет запрос на вставку строки dc.group_levels.
func ValidateCreateGroupLevel(req *tablesv1.CreateGroupLevelRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	groupLevelWritableFields(
		v,
		req.GetColumnId(),
		req.GetParentColumnId(),
		req.GetLevel(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateGroupLevelById проверяет запрос на обновление строки dc.group_levels.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateGroupLevelById(req *tablesv1.UpdateGroupLevelByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	groupLevelWritableFields(
		v,
		req.GetColumnId(),
		req.GetParentColumnId(),
		req.GetLevel(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}
