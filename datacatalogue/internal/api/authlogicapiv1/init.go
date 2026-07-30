// Package authlogicapiv1 реализует gRPC-сервис auth_logic.v1.AuthLogicService.
package authlogicapiv1

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	authlogicv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/auth_logic/v1"
)

// AuthLogicApiV1 — реализация auth_logic.v1.AuthLogicService.
type AuthLogicApiV1 struct {
	authlogicv1.UnimplementedAuthLogicServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *AuthLogicApiV1 {
	return &AuthLogicApiV1{
		services: service,
	}
}
