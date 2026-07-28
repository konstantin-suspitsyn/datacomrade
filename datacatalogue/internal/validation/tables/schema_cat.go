package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// schemaCatWritableFields проверяет поля, общие для вставки и обновления dc.schema_cat.
func schemaCatWritableFields(
	v *validator.Validator,
	databaseId int64,
	name string,
	userId int64,
) {
	v.Int64ID("database_id", databaseId)
	v.StringVarchar("name", name, schemaCatNameMaxLen)
	v.Int64ID("user_id", userId)
}

// ValidateCreateSchemaCat проверяет запрос на вставку строки dc.schema_cat.
func ValidateCreateSchemaCat(req *tablesv1.CreateSchemaCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	schemaCatWritableFields(
		v,
		req.GetDatabaseId(),
		req.GetName(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateSchemaCatById проверяет запрос на обновление строки dc.schema_cat.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateSchemaCatById(req *tablesv1.UpdateSchemaCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	schemaCatWritableFields(
		v,
		req.GetDatabaseId(),
		req.GetName(),
		req.GetUserId(),
	)

	return v.Err()
}
