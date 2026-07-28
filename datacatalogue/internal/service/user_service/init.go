// Package userservice содержит бизнес-логику поверх sqlc-репозитория user_model.
package userservice

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/user_model"
)

// UserService обслуживает таблицу dc."user".
type UserService struct {
	UserRepository *user_model.Queries
}

// New создаёт сервис поверх открытого соединения с базой.
func New(db *sql.DB) *UserService {
	return &UserService{
		UserRepository: user_model.New(db),
	}
}
