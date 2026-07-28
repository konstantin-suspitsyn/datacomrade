// Package userdomainrolesapiv1 реализует gRPC-сервис
// user_domain_roles.v1.UserDomainRolesService.
package userdomainrolesapiv1

import (
	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/services"
	userdomainrolesv1 "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/proto/user_domain_roles/v1"
)

// UserDomainRolesApiV1 — реализация user_domain_roles.v1.UserDomainRolesService.
type UserDomainRolesApiV1 struct {
	userdomainrolesv1.UnimplementedUserDomainRolesServiceServer

	services *services.ServiceLayer
}

// New создаёт хендлеры поверх слоя сервисов.
func New(service *services.ServiceLayer) *UserDomainRolesApiV1 {
	return &UserDomainRolesApiV1{
		services: service,
	}
}
