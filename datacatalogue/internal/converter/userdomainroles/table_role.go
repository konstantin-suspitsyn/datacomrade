package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// TableRoleToProto переводит строку dc.table_roles в сущность gRPC.
func TableRoleToProto(row user_domain_roles.DcTableRole) *userdomainrolesv1.TableRole {
	return &userdomainrolesv1.TableRole{
		Id:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   converter.TimeToProto(row.CreatedAt),
		UpdatedAt:   converter.TimeToProto(row.UpdatedAt),
		IsDeleted:   row.IsDeleted,
	}
}

// TableRolesToProto переводит список строк dc.table_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func TableRolesToProto(rows []user_domain_roles.DcTableRole) []*userdomainrolesv1.TableRole {
	items := make([]*userdomainrolesv1.TableRole, 0, len(rows))

	for _, row := range rows {
		items = append(items, TableRoleToProto(row))
	}

	return items
}

// ToCreateTableRoleParams собирает параметры вставки dc.table_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateTableRoleParams(req *userdomainrolesv1.CreateTableRoleRequest) user_domain_roles.CreateTableRoleParams {
	return user_domain_roles.CreateTableRoleParams{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
}

// ToUpdateTableRoleByIdParams собирает параметры обновления dc.table_roles из запроса gRPC.
func ToUpdateTableRoleByIdParams(req *userdomainrolesv1.UpdateTableRoleByIdRequest) user_domain_roles.UpdateTableRoleByIdParams {
	return user_domain_roles.UpdateTableRoleByIdParams{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}
}
