package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// DomainRoleToProto переводит строку dc.domain_roles в сущность gRPC.
func DomainRoleToProto(row user_domain_roles.DcDomainRole) *userdomainrolesv1.DomainRole {
	return &userdomainrolesv1.DomainRole{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		IsDeleted:   row.IsDeleted,
	}
}

// DomainRolesToProto переводит список строк dc.domain_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DomainRolesToProto(rows []user_domain_roles.DcDomainRole) []*userdomainrolesv1.DomainRole {
	items := make([]*userdomainrolesv1.DomainRole, 0, len(rows))

	for _, row := range rows {
		items = append(items, DomainRoleToProto(row))
	}

	return items
}

// ToCreateDomainRoleParams собирает параметры вставки dc.domain_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDomainRoleParams(req *userdomainrolesv1.CreateDomainRoleRequest) user_domain_roles.CreateDomainRoleParams {
	return user_domain_roles.CreateDomainRoleParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
}

// ToUpdateDomainRoleByIdParams собирает параметры обновления dc.domain_roles из запроса gRPC.
func ToUpdateDomainRoleByIdParams(req *userdomainrolesv1.UpdateDomainRoleByIdRequest) user_domain_roles.UpdateDomainRoleByIdParams {
	return user_domain_roles.UpdateDomainRoleByIdParams{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
}
