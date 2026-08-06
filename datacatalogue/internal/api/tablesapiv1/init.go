// Package tablesapiv1 реализует gRPC-сервисы tables.v1.AliasService,
// tables.v1.UserService и tables.v1.HostService — домен tables_model собран
// из отдельного сервиса на таблицу, а не одного общего, как в остальных
// доменах.
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

// TablesApiV1 — реализация tables.v1.AliasService, tables.v1.UserService и
// tables.v1.HostService.
type TablesApiV1 struct {
	tablesv1.UnimplementedAliasServiceServer
	tablesv1.UnimplementedUserServiceServer
	tablesv1.UnimplementedHostServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *TablesApiV1 {
	return &TablesApiV1{
		services: service,
	}
}
