package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// domainCatWritableFields проверяет поля, общие для вставки и обновления dc.domain_cat.
func domainCatWritableFields(
	v *validator.Validator,
	domainName string,
	userExternalId string,
) {
	v.StringVarchar("domain_name", domainName, domainCatDomainNameMaxLen)
	v.StringUUID("user_id", userExternalId)
}

// ValidateCreateDomainCat проверяет запрос на вставку строки dc.domain_cat.
func ValidateCreateDomainCat(req *tablesv1.CreateDomainCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	domainCatWritableFields(
		v,
		req.GetDomainName(),
		req.GetUserExternalId(),
	)

	return v.Err()
}

// ValidateUpdateDomainCatById проверяет запрос на обновление строки dc.domain_cat.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateDomainCatById(req *tablesv1.UpdateDomainCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	domainCatWritableFields(
		v,
		req.GetDomainName(),
		req.GetUserExternalId(),
	)

	return v.Err()
}
