package olapexchange

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/converter"
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/olapexchange"
	olapexchangev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/olap_exchange/v1"
)

// ToGetOlapAliasesParams собирает параметры olap_exchange.GetOlapAliases из запроса gRPC.
func ToGetOlapAliasesParams(req *olapexchangev1.GetOlapAliasesRequest) olapexchange.GetOlapAliasesParams {
	return olapexchange.GetOlapAliasesParams{
		ExternalID: converter.ProtoToUUID(req.GetExternalId()),
		Order:      req.GetOrder(),
	}
}

// GetOlapAliasesRowToProto переводит строку запроса olap_exchange.GetOlapAliases в сообщение gRPC.
func GetOlapAliasesRowToProto(row olapexchange.GetOlapAliasesRow) *olapexchangev1.GetOlapAliasesRow {
	return &olapexchangev1.GetOlapAliasesRow{
		Id:                    row.ID,
		Name:                  row.Name,
		Description:           row.Description,
		IsDeleted:             row.IsDeleted,
		ColumnTypeId:          row.ColumnTypeID,
		ColumnTypeName:        row.ColumnTypeName,
		ColumnTypeDescription: row.ColumnTypeDescription,
	}
}

// GetOlapAliasesToProto переводит результат olap_exchange.GetOlapAliases в ответ gRPC.
func GetOlapAliasesToProto(rows []olapexchange.GetOlapAliasesRow) *olapexchangev1.GetOlapAliasesResponse {
	items := make([]*olapexchangev1.GetOlapAliasesRow, 0, len(rows))

	for _, row := range rows {
		items = append(items, GetOlapAliasesRowToProto(row))
	}

	return &olapexchangev1.GetOlapAliasesResponse{Rows: items}
}
