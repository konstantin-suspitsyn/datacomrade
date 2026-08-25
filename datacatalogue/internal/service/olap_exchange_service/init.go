// Package olapexchangeservice содержит бизнес-логику поверх sqlc-репозитория
// olap_exchange.
package olapexchangeservice

import (
	"database/sql"

	"github.com/konstantin-suspitsyn/datacomrade/datacatalogue/internal/repository/olapexchange"
)

// OlapExchangeService отдаёт сквозные выборки поверх нескольких таблиц для
// OLAP-обмена (например, список алиасов, доступных пользователю по его
// ролям), без своего CRUD.
type OlapExchangeService struct {
	OlapExchangeRepository *olapexchange.Queries
}

// New создаёт сервис поверх открытого соединения с базой.
func New(db *sql.DB) *OlapExchangeService {
	return &OlapExchangeService{
		OlapExchangeRepository: olapexchange.New(db),
	}
}
