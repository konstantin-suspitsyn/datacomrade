// This is a service layer
// Will be user to contain all sevices as struct

package services

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/internal/accesscontrol"
	"github.com/konstantin-suspitsyn/datacomrade/internal/roles"
	"github.com/konstantin-suspitsyn/datacomrade/internal/shareddata"
	"github.com/konstantin-suspitsyn/datacomrade/internal/users"
)

type ServiceLayer struct {
	UserService          *users.UserService
	SharedDataService    *shareddata.SharedDataService
	AccessControlService *accesscontrol.AccessControlService
	RoleService          *roles.RoleService
}

func New(db *sql.DB) *ServiceLayer {
	return &ServiceLayer{
		UserService:          users.New(db),
		SharedDataService:    shareddata.New(db),
		AccessControlService: accesscontrol.New(db),
		RoleService:          roles.New(db),
	}
}
