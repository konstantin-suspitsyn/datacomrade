package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateTablesJoins проверяет запрос на вставку строки dc.tables_joins.
func ValidateCreateTablesJoins(req *tablesv1.CreateTablesJoinsRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("table_cat_id_left", req.GetTableCatIdLeft())
	v.Int64ID("table_cat_id_right", req.GetTableCatIdRight())
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateTablesJoinsById проверяет запрос на обновление строки dc.tables_joins.
func ValidateUpdateTablesJoinsById(req *tablesv1.UpdateTablesJoinsByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("table_cat_id_left", req.GetTableCatIdLeft())
	v.Int64ID("table_cat_id_right", req.GetTableCatIdRight())
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteTablesJoinsById проверяет запрос на мягкое удаление строки dc.tables_joins.
func ValidateDeleteTablesJoinsById(req *tablesv1.DeleteTablesJoinsByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetTablesJoins проверяет запрос страницы dc.tables_joins.
func ValidateGetTablesJoins(req *tablesv1.GetTablesJoinsRequest) error {
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
