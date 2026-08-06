// Package validation проверяет поля входящих запросов gRPC до того,
// как они дойдут до конвертера и репозитория.
//
// Правила по сущностям лежат в подпакетах, повторяющих деление
// на sqlc-пакеты и gRPC-сервисы:
//
//	validation/tables            — dc.alias, dc.host, dc.table_cat и остальные 15 таблиц
//	validation/user              — dc.user
//	validation/userdomainroles   — dc.domain_roles и остальные 6 таблиц
//
// Каждая функция возвращает *validator.ValidationError со списком проблем
// по каждому полю или nil, если запрос корректен. Границы длин строк
// в подпакетах вынесены в limits.go и повторяют schema.sql.
package validation

import "github.com/konstantin-suspitsyn/datacomrade/shared/pkg/validator"

// Границы постраничных выборок. Значения — политика сервиса, а не описание
// запроса, поэтому живут здесь, а не в schema.sql/schema.json.
const (
	DefaultPageSize = 50
	MaxPageSize     = 200
)

// ValidateID проверяет идентификатор записи для запросов, у которых
// в теле только id: Get*ById, GetDeleted*ById, Delete*ById, Undelete*ById.
// Идентификаторы — bigserial, поэтому допустимы только положительные значения.
func ValidateID(id int64) error {
	v := validator.New()
	v.Int64ID("id", id)

	return v.Err()
}
