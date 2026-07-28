package userdomainroles

// Ограничения длины строковых колонок. Значения взяты из
// datacatalogue/db/sqlc/user_domain_roles/schema.sql и должны меняться вместе
// с ней: валидация обязана отсекать слишком длинную строку раньше,
// чем её отвергнет Postgres.
const (
	// dc.domain_roles
	domainRoleNameMaxLen        = 128  // name character varying(128)
	domainRoleDescriptionMaxLen = 2000 // description character varying(2000)

	// dc.table_roles
	tableRoleNameMaxLen        = 128  // name character varying(128)
	tableRoleDescriptionMaxLen = 2000 // description character varying(2000)
)
