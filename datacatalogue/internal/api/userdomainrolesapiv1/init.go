// Package userdomainrolesapiv1 реализует gRPC-сервисы файла
// user_domain_roles.v1: UserDomainRolesService и UserTableRolesService.
package userdomainrolesapiv1

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// UserDomainRolesApiV1 — реализация обоих сервисов файла: у каждой таблицы
// свой service в .proto, но обслуживает их одна Go-структура.
type UserDomainRolesApiV1 struct {
	userdomainrolesv1.UnimplementedUserDomainRolesServiceServer
	userdomainrolesv1.UnimplementedUserTableRolesServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *UserDomainRolesApiV1 {
	return &UserDomainRolesApiV1{
		services: service,
	}
}
