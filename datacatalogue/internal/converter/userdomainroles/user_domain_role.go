package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// UserDomainRoleToProto переводит строку dc.user_domain_roles в сущность gRPC.
func UserDomainRoleToProto(row user_domain_roles.DcUserDomainRole) *userdomainrolesv1.UserDomainRole {
	return &userdomainrolesv1.UserDomainRole{
		Id:            row.ID,
		UserId:        row.UserID,
		DomainRolesId: row.DomainRolesID,
		CreatedAt:     converter.TimeToProto(row.CreatedAt),
		UpdatedAt:     converter.TimeToProto(row.UpdatedAt),
		IsDeleted:     row.IsDeleted,
		DomainId:      row.DomainID,
	}
}

// UserDomainRolesToProto переводит список строк dc.user_domain_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func UserDomainRolesToProto(rows []user_domain_roles.DcUserDomainRole) []*userdomainrolesv1.UserDomainRole {
	items := make([]*userdomainrolesv1.UserDomainRole, 0, len(rows))

	for _, row := range rows {
		items = append(items, UserDomainRoleToProto(row))
	}

	return items
}

// ToCreateUserDomainRoleParams собирает параметры вставки dc.user_domain_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateUserDomainRoleParams(req *userdomainrolesv1.CreateUserDomainRoleRequest) user_domain_roles.CreateUserDomainRoleParams {
	return user_domain_roles.CreateUserDomainRoleParams{
		UserID:        req.GetUserId(),
		DomainRolesID: req.GetDomainRolesId(),
		DomainID:      req.GetDomainId(),
	}
}

// ToUpdateUserDomainRoleByIdParams собирает параметры обновления dc.user_domain_roles из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateUserDomainRoleByIdParams(req *userdomainrolesv1.UpdateUserDomainRoleByIdRequest) user_domain_roles.UpdateUserDomainRoleByIdParams {
	return user_domain_roles.UpdateUserDomainRoleByIdParams{
		ID:            req.GetId(),
		UserID:        req.GetUserId(),
		DomainRolesID: req.GetDomainRolesId(),
		DomainID:      req.GetDomainId(),
	}
}
