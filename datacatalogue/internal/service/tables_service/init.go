// Package tablesservice содержит бизнес-логику поверх sqlc-репозитория
// tables_model: проверки, которые нельзя выразить одним запросом,
// и обогащение ошибок репозитория контекстом.
package tablesservice

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/tables_model"
)

// TablesService обслуживает 15 таблиц каталога из схемы dc.
type TablesService struct {
	TablesRepository *tables_model.Queries
}

// New создаёт сервис поверх открытого соединения с базой.
func New(db *sql.DB) *TablesService {
	return &TablesService{
		TablesRepository: tables_model.New(db),
	}
}
