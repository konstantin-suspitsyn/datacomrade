// Package userdomainrolesservice содержит бизнес-логику поверх sqlc-репозитория
// user_domain_roles.
package userdomainrolesservice

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_domain_roles"
)

// UserDomainRolesService обслуживает 6 таблиц ролей: dc.domain_roles,
// dc.table_roles и связующие к ним.
type UserDomainRolesService struct {
	UserDomainRolesRepository *user_domain_roles.Queries
}

// New создаёт сервис поверх открытого соединения с базой.
func New(db *sql.DB) *UserDomainRolesService {
	return &UserDomainRolesService{
		UserDomainRolesRepository: user_domain_roles.New(db),
	}
}
