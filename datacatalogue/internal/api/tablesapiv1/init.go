// Package tablesapiv1 реализует gRPC-сервис tables.v1.TableService.
//
// Каждый метод устроен одинаково:
//
//	validation.Validate<Rpc>(req) → converter.To<Rpc>Params(req)
//	    → services.TablesService.<Rpc>(ctx, params)
//	    → converter.<Entity>ToProto(row) → Response
//
// Встроенный UnimplementedTableServiceServer отдаёт Unimplemented по тем
// RPC, которые ещё не реализованы.
package tablesapiv1

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	tablesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/tables/v1"
)

// TablesApiV1 — реализация tables.v1.TableService.
type TablesApiV1 struct {
	tablesv1.UnimplementedTableServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *TablesApiV1 {
	return &TablesApiV1{
		services: service,
	}
}
