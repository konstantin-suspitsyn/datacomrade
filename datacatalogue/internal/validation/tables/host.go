package tables

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/validation"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// ValidateCreateHost проверяет запрос на вставку строки dc.host.
func ValidateCreateHost(req *tablesv1.CreateHostRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringVarchar("name", req.GetName(), hostNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), hostDescriptionMaxLen)
	v.StringVarchar("host_env", req.GetHostEnv(), hostHostEnvMaxLen)
	v.StringVarchar("port_env", req.GetPortEnv(), hostPortEnvMaxLen)
	v.StringVarchar("username_env", req.GetUsernameEnv(), hostUsernameEnvMaxLen)
	v.StringVarchar("password_env", req.GetPasswordEnv(), hostPasswordEnvMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateUpdateHostById проверяет запрос на обновление строки dc.host.
func ValidateUpdateHostById(req *tablesv1.UpdateHostByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())
	v.StringVarchar("name", req.GetName(), hostNameMaxLen)
	v.StringVarchar("description", req.GetDescription(), hostDescriptionMaxLen)
	v.StringVarchar("host_env", req.GetHostEnv(), hostHostEnvMaxLen)
	v.StringVarchar("port_env", req.GetPortEnv(), hostPortEnvMaxLen)
	v.StringVarchar("username_env", req.GetUsernameEnv(), hostUsernameEnvMaxLen)
	v.StringVarchar("password_env", req.GetPasswordEnv(), hostPasswordEnvMaxLen)
	v.StringUUID("user_id", req.GetExternalId())

	return v.Err()
}

// ValidateDeleteHostById проверяет запрос на мягкое удаление строки dc.host.
func ValidateDeleteHostById(req *tablesv1.DeleteHostByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateUndeleteHostById проверяет запрос на обратное удаление строки dc.host.
func ValidateUndeleteHostById(req *tablesv1.UndeleteHostByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.StringUUID("user_id", req.GetExternalId())
	v.Int64ID("id", req.GetId())

	return v.Err()
}

// ValidateGetHosts проверяет запрос страницы dc.host.
func ValidateGetHosts(req *tablesv1.GetHostsRequest) error {
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

// ValidateGetHostsSearchName проверяет запрос страницы dc.host.
func ValidateGetHostsSearchName(req *tablesv1.GetHostsSearchNameRequest) error {
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

// ValidateGetHostDeleted проверяет запрос страницы dc.host.
func ValidateGetHostDeleted(req *tablesv1.GetHostDeletedRequest) error {
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
