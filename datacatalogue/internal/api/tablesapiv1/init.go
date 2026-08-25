// Package tablesapiv1 реализует gRPC-сервисы tables.v1 — домен tables_model
// собран из отдельного сервиса на таблицу (AliasService, HostService,
// TableCatService и остальные), а не одного общего, как в других доменах.
//
// Каждый метод устроен одинаково:
//
//	validation.Validate<Rpc>(req) → converter.To<Rpc>Params(req)
//	    → services.TablesService.<Rpc>(ctx, params)
//	    → converter.<Entity>ToProto(row) → Response
//
// Встроенные Unimplemented*ServiceServer отдают Unimplemented по тем RPC,
// которые ещё не реализованы.
package tablesapiv1

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// TablesApiV1 — реализация всех gRPC-сервисов tables.v1: по сервису на
// таблицу домена tables_model.
type TablesApiV1 struct {
	tablesv1.UnimplementedAliasServiceServer
	tablesv1.UnimplementedUserServiceServer
	tablesv1.UnimplementedHostServiceServer
	tablesv1.UnimplementedDatabaseTypeServiceServer
	tablesv1.UnimplementedDatabaseCatServiceServer
	tablesv1.UnimplementedColumnCatServiceServer
	tablesv1.UnimplementedColumnTypeServiceServer
	tablesv1.UnimplementedDatabaseCalculationServiceServer
	tablesv1.UnimplementedDomainCatServiceServer
	tablesv1.UnimplementedFollowingCalculationServiceServer
	tablesv1.UnimplementedGroupLevelsServiceServer
	tablesv1.UnimplementedHasToGroupServiceServer
	tablesv1.UnimplementedSchemaCatServiceServer
	tablesv1.UnimplementedTableCatServiceServer
	tablesv1.UnimplementedTableTypeServiceServer
	tablesv1.UnimplementedCalculationTypeServiceServer
	tablesv1.UnimplementedTablesJoinsServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *TablesApiV1 {
	return &TablesApiV1{
		services: service,
	}
}
