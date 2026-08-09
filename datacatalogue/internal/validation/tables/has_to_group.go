package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateHasToGroup проверяет запрос на вставку строки dc.has_to_group.
func ValidateCreateHasToGroup(req *tablesv1.CreateHasToGroupRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("column_id_a", req.GetColumnIdA())
	v.Int64ID("column_id_b", req.GetColumnIdB())
	v.StringVarchar("description", req.GetDescription(), hasToGroupDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateHasToGroupById проверяет запрос на обновление строки dc.has_to_group.
func ValidateUpdateHasToGroupById(req *tablesv1.UpdateHasToGroupByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("column_id_a", req.GetColumnIdA())
	v.Int64ID("column_id_b", req.GetColumnIdB())
	v.StringVarchar("description", req.GetDescription(), hasToGroupDescriptionMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteHasToGroupById проверяет запрос на мягкое удаление строки dc.has_to_group.
func ValidateDeleteHasToGroupById(req *tablesv1.DeleteHasToGroupByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteHasToGroupById проверяет запрос на обратное удаление строки dc.has_to_group.
func ValidateUndeleteHasToGroupById(req *tablesv1.UndeleteHasToGroupByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetHasToGroupsSearchName проверяет запрос страницы dc.has_to_group.
func ValidateGetHasToGroupsSearchName(req *tablesv1.GetHasToGroupsSearchNameRequest) error {
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

// ValidateGetHasToGroups проверяет запрос страницы dc.has_to_group.
func ValidateGetHasToGroups(req *tablesv1.GetHasToGroupsRequest) error {
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

// ValidateGetHasToGroupsDeleted проверяет запрос страницы dc.has_to_group.
func ValidateGetHasToGroupsDeleted(req *tablesv1.GetHasToGroupsDeletedRequest) error {
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
