package tables

// Ограничения длины строковых колонок. Значения взяты из
// datacatalogue/db/sqlc/tables_model/schema.sql и должны меняться вместе
// с ней: валидация обязана отсекать слишком длинную строку раньше,
// чем её отвергнет Postgres.
const (
	// dc.alias
	aliasNameMaxLen        = 255  // name character varying(255)
	aliasDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.user
	userNameMaxLen = 512 // name character varying(512)

	// dc.host
	hostNameMaxLen        = 255  // name character varying(255)
	hostDescriptionMaxLen = 1000 // description character varying(1000)
	hostHostEnvMaxLen     = 255  // host_env character varying(255)
	hostPortEnvMaxLen     = 255  // port_env character varying(255)
	hostUsernameEnvMaxLen = 255  // username_env character varying(255)
	hostPasswordEnvMaxLen = 255  // password_env character varying(255)
)
