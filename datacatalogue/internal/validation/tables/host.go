package tables

import (
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
	"github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"
)

// hostWritableFields проверяет поля, общие для вставки и обновления dc.host.
func hostWritableFields(
	v *validator.Validator,
	name string,
	description string,
	hostEnv string,
	portEnv string,
	usernameEnv string,
	passwordEnv string,
	userExternalId string,
) {
	v.StringVarchar("name", name, hostNameMaxLen)
	v.StringVarchar("description", description, hostDescriptionMaxLen)
	v.StringVarchar("host_env", hostEnv, hostHostEnvMaxLen)
	v.StringVarchar("port_env", portEnv, hostPortEnvMaxLen)
	v.StringVarchar("username_env", usernameEnv, hostUsernameEnvMaxLen)
	v.StringVarchar("password_env", passwordEnv, hostPasswordEnvMaxLen)
	v.StringUUID("user_id", userExternalId)
}

// ValidateCreateHost проверяет запрос на вставку строки dc.host.
func ValidateCreateHost(req *tablesv1.CreateHostRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	hostWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetHostEnv(),
		req.GetPortEnv(),
		req.GetUsernameEnv(),
		req.GetPasswordEnv(),
		req.GetUserExternalId(),
	)

	return v.Err()
}

// ValidateUpdateHostById проверяет запрос на обновление строки dc.host.
// К изменяемым полям добавляется id обновляемой записи.
func ValidateUpdateHostById(req *tablesv1.UpdateHostByIdRequest) error {
	v := validator.New()

	if req == nil {
		v.AddError("request", validator.MsgRequired)
		return v.Err()
	}

	v.Int64ID("id", req.GetId())

	hostWritableFields(
		v,
		req.GetName(),
		req.GetDescription(),
		req.GetHostEnv(),
		req.GetPortEnv(),
		req.GetUsernameEnv(),
		req.GetPasswordEnv(),
		req.GetUserExternalId(),
	)

	return v.Err()
}
