package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateDomainCat проверяет запрос на вставку строки dc.domain_cat.
func ValidateCreateDomainCat(req *tablesv1.CreateDomainCatRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("domain_name", req.GetDomainName(), domainCatDomainNameMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateDomainCatById проверяет запрос на обновление строки dc.domain_cat.
func ValidateUpdateDomainCatById(req *tablesv1.UpdateDomainCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("domain_name", req.GetDomainName(), domainCatDomainNameMaxLen)
	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateDeleteDomainCatById проверяет запрос на мягкое удаление строки dc.domain_cat.
func ValidateDeleteDomainCatById(req *tablesv1.DeleteDomainCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteDomainCatById проверяет запрос на обратное удаление строки dc.domain_cat.
func ValidateUndeleteDomainCatById(req *tablesv1.UndeleteDomainCatByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetDomainCatsSearchName проверяет запрос страницы dc.domain_cat.
func ValidateGetDomainCatsSearchName(req *tablesv1.GetDomainCatsSearchNameRequest) error {
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

// ValidateGetDomainCats проверяет запрос страницы dc.domain_cat.
func ValidateGetDomainCats(req *tablesv1.GetDomainCatsRequest) error {
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

// ValidateGetDomainCatsDeleted проверяет запрос страницы dc.domain_cat.
func ValidateGetDomainCatsDeleted(req *tablesv1.GetDomainCatsDeletedRequest) error {
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
