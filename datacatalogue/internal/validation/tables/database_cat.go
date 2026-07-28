package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// databaseCatWritableFields проверяет поля, общие для вставки и обновления dc.database_cat.
func databaseCatWritableFields(
	v *validator.Validator,
	name string,
	hostId int64,
	databaseTypeId int64,
	description string,
	userId int64,
) {
	v.StringVarchar("name", name, databaseCatNameMaxLen)
	v.Int64ID("host_id", hostId)
	v.Int64ID("database_type_id", databaseTypeId)
	v.StringVarchar("description", description, databaseCatDescriptionMaxLen)
	v.Int64ID("user_id", userId)
}

// ValidateCreateDatabaseCat проверяет запрос на вставку строки dc.database_cat.
func ValidateCreateDatabaseCat(req *tablesv1.CreateDatabaseCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	databaseCatWritableFields(
		v,
		req.GetName(),
		req.GetHostId(),
		req.GetDatabaseTypeId(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateDatabaseCatById проверяет запрос на обновление строки dc.database_cat.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateDatabaseCatById(req *tablesv1.UpdateDatabaseCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	databaseCatWritableFields(
		v,
		req.GetName(),
		req.GetHostId(),
		req.GetDatabaseTypeId(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}
