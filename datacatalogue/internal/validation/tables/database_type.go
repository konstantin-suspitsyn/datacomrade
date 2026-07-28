package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// databaseTypeWritableFields проверяет поля, общие для вставки и обновления dc.database_type.
func databaseTypeWritableFields(
	v *validator.Validator,
	name string,
	dbVersion string,
	userId int64,
) {
	v.StringVarchar("name", name, databaseTypeNameMaxLen)
	v.StringVarchar("db_version", dbVersion, databaseTypeDbVersionMaxLen)
	v.Int64ID("user_id", userId)
}

// ValidateCreateDatabaseType проверяет запрос на вставку строки dc.database_type.
func ValidateCreateDatabaseType(req *tablesv1.CreateDatabaseTypeRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	databaseTypeWritableFields(
		v,
		req.GetName(),
		req.GetDbVersion(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateDatabaseTypeById проверяет запрос на обновление строки dc.database_type.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateDatabaseTypeById(req *tablesv1.UpdateDatabaseTypeByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	databaseTypeWritableFields(
		v,
		req.GetName(),
		req.GetDbVersion(),
		req.GetUserId(),
	)

	return v.Err()
}
