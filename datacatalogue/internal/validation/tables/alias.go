package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// aliasWritableFields проверяет поля, общие для вставки и обновления dc.alias.
func aliasWritableFields(
	v *validator.Validator,
	name string,
	description string,
	userId int64,
) {
	v.StringVarchar("name", name, aliasNameMaxLen)
	v.StringVarchar("description", description, aliasDescriptionMaxLen)
	v.Int64ID("user_id", userId)
}

// ValidateCreateAlias проверяет запрос на вставку строки dc.alias.
func ValidateCreateAlias(req *tablesv1.CreateAliasRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	aliasWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}

// ValidateUpdateAliasById проверяет запрос на обновление строки dc.alias.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateAliasById(req *tablesv1.UpdateAliasByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	aliasWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetUserId(),
	)

	return v.Err()
}
