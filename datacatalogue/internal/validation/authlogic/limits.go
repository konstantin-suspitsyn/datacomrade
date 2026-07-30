package authlogic

// Ограничения длины строковых аргументов. Значения взяты из
// datacatalogue/db/sqlc/auth_logic/schema.sql и должны меняться вместе
// с ней: валидация обязана отсекать слишком длинную строку раньше,
// чем её отвергнет Postgres.
const (
	nameMaxLen = 128 // name character varying(128)
)
