package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// DomainsDomainRoleToProto переводит строку dc.domains_domain_roles в сущность gRPC.
func DomainsDomainRoleToProto(row user_domain_roles.DcDomainsDomainRole) *userdomainrolesv1.DomainsDomainRole {
	return &userdomainrolesv1.DomainsDomainRole{
		Id:            row.ID,
		DomainCatId:   row.DomainCatID,
		DomainRolesId: row.DomainRolesID,
		CreatedAt:     converter.TimeToProto(row.CreatedAt),
		UpdatedAt:     converter.TimeToProto(row.UpdatedAt),
		IsDeleted:     row.IsDeleted,
	}
}

// DomainsDomainRolesToProto переводит список строк dc.domains_domain_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func DomainsDomainRolesToProto(rows []user_domain_roles.DcDomainsDomainRole) []*userdomainrolesv1.DomainsDomainRole {
	items := make([]*userdomainrolesv1.DomainsDomainRole, 0, len(rows))

	for _, row := range rows {
		items = append(items, DomainsDomainRoleToProto(row))
	}

	return items
}

// ToCreateDomainsDomainRoleParams собирает параметры вставки dc.domains_domain_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateDomainsDomainRoleParams(req *userdomainrolesv1.CreateDomainsDomainRoleRequest) user_domain_roles.CreateDomainsDomainRoleParams {
	return user_domain_roles.CreateDomainsDomainRoleParams{
		DomainCatID:   req.GetDomainCatId(),
		DomainRolesID: req.GetDomainRolesId(),
	}
}

// ToUpdateDomainsDomainRoleByIdParams собирает параметры обновления dc.domains_domain_roles из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateDomainsDomainRoleByIdParams(req *userdomainrolesv1.UpdateDomainsDomainRoleByIdRequest) user_domain_roles.UpdateDomainsDomainRoleByIdParams {
	return user_domain_roles.UpdateDomainsDomainRoleByIdParams{
		ID:            req.GetId(),
		DomainCatID:   req.GetDomainCatId(),
		DomainRolesID: req.GetDomainRolesId(),
	}
}
