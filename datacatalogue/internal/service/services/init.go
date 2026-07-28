// Package services собирает все сервисы домена в один объект,
// который api-слой получает при создании.
package services

import (
	"database/sql"

	tablesservice "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/tables_service"
	userdomainrolesservice "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/user_domain_roles_service"
	userservice "github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/service/user_service"
)

// ServiceLayer — агрегат сервисов по числу sqlc-пакетов и gRPC-сервисов.
type ServiceLayer struct {
	TablesService          *tablesservice.TablesService
	UserDomainRolesService *userdomainrolesservice.UserDomainRolesService
	UserService            *userservice.UserService
}

// New создаёт все сервисы поверх одного соединения с базой.
func New(db *sql.DB) *ServiceLayer {
	return &ServiceLayer{
		TablesService:          tablesservice.New(db),
		UserDomainRolesService: userdomainrolesservice.New(db),
		UserService:            userservice.New(db),
	}
}
