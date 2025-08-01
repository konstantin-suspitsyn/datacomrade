// This is a service layer
// Will be user to contain all sevices as struct

package services

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/internal/users"
)

type ServiceLayer struct {
	UserService *users.UserService
}

func New(db *sql.DB) *ServiceLayer {
	return &ServiceLayer{
		UserService: users.New(db),
	}
}
