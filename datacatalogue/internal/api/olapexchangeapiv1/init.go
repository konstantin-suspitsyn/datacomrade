// Package olapexchangeapiv1 реализует gRPC-сервис olap_exchange.v1.AliasService.
package olapexchangeapiv1

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	olapexchangev1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/olap_exchange/v1"
)

// OlapExchangeApiV1 — реализация olap_exchange.v1.AliasService.
type OlapExchangeApiV1 struct {
	olapexchangev1.UnimplementedAliasServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *OlapExchangeApiV1 {
	return &OlapExchangeApiV1{
		services: service,
	}
}
