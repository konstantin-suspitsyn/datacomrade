// Package userapiv1 реализует gRPC-сервис user.v1.UserService.
package userapiv1

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	userv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user/v1"
)

// UserApiV1 — реализация user.v1.UserService.
type UserApiV1 struct {
	userv1.UnimplementedUserServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *UserApiV1 {
	return &UserApiV1{
		services: service,
	}
}
