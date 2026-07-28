package user

// Ограничения длины строковых колонок. Значения взяты из
// datacatalogue/db/sqlc/user_model/schema.sql и должны меняться вместе
// с ней: валидация обязана отсекать слишком длинную строку раньше,
// чем её отвергнет Postgres.
const (
	// dc.user
	userNameMaxLen = 512 // name character varying(512)
)
