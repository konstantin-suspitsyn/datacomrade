package tables

import (
	"math"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateGroupLevels проверяет запрос на вставку строки dc.group_levels.
func ValidateCreateGroupLevels(req *tablesv1.CreateGroupLevelsRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("column_id", req.GetColumnId())
	v.Int64ID("parent_column_id", req.GetParentColumnId())
	v.Int32Between("level", req.GetLevel(), math.MinInt16, math.MaxInt16)
	v.StringVarchar("description", req.GetDescription(), groupLevelsDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateGroupLevelsById проверяет запрос на обновление строки dc.group_levels.
func ValidateUpdateGroupLevelsById(req *tablesv1.UpdateGroupLevelsByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("column_id", req.GetColumnId())
	v.Int64ID("parent_column_id", req.GetParentColumnId())
	v.Int32Between("level", req.GetLevel(), math.MinInt16, math.MaxInt16)
	v.StringVarchar("description", req.GetDescription(), groupLevelsDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteGroupLevelsById проверяет запрос на мягкое удаление строки dc.group_levels.
func ValidateDeleteGroupLevelsById(req *tablesv1.DeleteGroupLevelsByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteGroupLevelsById проверяет запрос на обратное удаление строки dc.group_levels.
func ValidateUndeleteGroupLevelsById(req *tablesv1.UndeleteGroupLevelsByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetGroupLevelsDeleted проверяет запрос списка dc.group_levels.
func ValidateGetGroupLevelsDeleted(req *tablesv1.GetGroupLevelsDeletedRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringIn("order", req.GetOrder(), "", "ASC", "DESC")

	return v.Err()
}

// ValidateGetGroupLevelsSearchName проверяет запрос страницы dc.group_levels.
func ValidateGetGroupLevelsSearchName(req *tablesv1.GetGroupLevelsSearchNameRequest) error {
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

// ValidateGetGroupLevels проверяет запрос страницы dc.group_levels.
func ValidateGetGroupLevels(req *tablesv1.GetGroupLevelsRequest) error {
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

// ValidateGetGroupLevelsFiltered проверяет запрос страницы dc.group_levels.
func ValidateGetGroupLevelsFiltered(req *tablesv1.GetGroupLevelsFilteredRequest) error {
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
