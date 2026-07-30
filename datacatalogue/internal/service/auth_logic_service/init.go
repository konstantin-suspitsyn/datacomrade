// Package authlogicservice содержит бизнес-логику поверх sqlc-репозитория
// auth_logic.
package authlogicservice

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/auth_logic"
)

// AuthLogicService отдаёт id таблиц dc.table_cat, доступных пользователю
// по его ролям — сквозные выборки поверх нескольких таблиц, без своего CRUD.
type AuthLogicService struct {
	AuthLogicRepository *auth_logic.Queries
}

// New создаёт сервис поверх открытого соединения с базой.
func New(db *sql.DB) *AuthLogicService {
	return &AuthLogicService{
		AuthLogicRepository: auth_logic.New(db),
	}
}
