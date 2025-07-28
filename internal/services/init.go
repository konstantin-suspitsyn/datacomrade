// This is a service layer
// Will be user to contain all sevices as struct

package services

import (
	"database/sql"
)

type ServiceLayer struct {
	db *sql.DB
}

func New(db *sql.DB) *ServiceLayer {
	return &ServiceLayer{
		db: db,
	}
}
