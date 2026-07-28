package userdomainroles

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// TablesTableRoleToProto переводит строку dc.tables_table_roles в сущность gRPC.
func TablesTableRoleToProto(row user_domain_roles.DcTablesTableRole) *userdomainrolesv1.TablesTableRole {
	return &userdomainrolesv1.TablesTableRole{
		Id:           row.ID,
		TableCatId:   row.TableCatID,
		TableRolesId: row.TableRolesID,
		CreatedAt:    converter.TimeToProto(row.CreatedAt),
		UpdatedAt:    converter.TimeToProto(row.UpdatedAt),
		IsDeleted:    row.IsDeleted,
	}
}

// TablesTableRolesToProto переводит список строк dc.tables_table_roles в список сущностей gRPC.
// Для пустого входа возвращается пустой, а не nil-слайс.
func TablesTableRolesToProto(rows []user_domain_roles.DcTablesTableRole) []*userdomainrolesv1.TablesTableRole {
	items := make([]*userdomainrolesv1.TablesTableRole, 0, len(rows))

	for _, row := range rows {
		items = append(items, TablesTableRoleToProto(row))
	}

	return items
}

// ToCreateTablesTableRoleParams собирает параметры вставки dc.tables_table_roles из запроса gRPC.
// id, is_deleted, created_at и updated_at не переносятся — их выставляет SQL.
func ToCreateTablesTableRoleParams(req *userdomainrolesv1.CreateTablesTableRoleRequest) user_domain_roles.CreateTablesTableRoleParams {
	return user_domain_roles.CreateTablesTableRoleParams{
		TableCatID:   req.GetTableCatId(),
		TableRolesID: req.GetTableRolesId(),
	}
}

// ToUpdateTablesTableRoleByIdParams собирает параметры обновления dc.tables_table_roles из запроса gRPC.
// updated_at выставляет SQL, is_deleted через обновление не меняется.
func ToUpdateTablesTableRoleByIdParams(req *userdomainrolesv1.UpdateTablesTableRoleByIdRequest) user_domain_roles.UpdateTablesTableRoleByIdParams {
	return user_domain_roles.UpdateTablesTableRoleByIdParams{
		ID:           req.GetId(),
		TableCatID:   req.GetTableCatId(),
		TableRolesID: req.GetTableRolesId(),
	}
}
