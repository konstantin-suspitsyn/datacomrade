package tables

// Ограничения длины строковых колонок. Значения взяты из
// datacatalogue/db/sqlc/tables_model/schema.sql и должны меняться вместе
// с ней: валидация обязана отсекать слишком длинную строку раньше,
// чем её отвергнет Postgres.
const (
	// dc.alias
	aliasNameMaxLen        = 255  // name character varying(255)
	aliasDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.column_cat
	columnCatNameMaxLen        = 256  // name character varying(256)
	columnCatDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.column_type
	columnTypeNameMaxLen        = 128  // name character varying(128)
	columnTypeDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.calculation_type
	calculationTypeNameMaxLen        = 52   // name character varying(52)
	calculationTypeDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.database_cat
	databaseCatNameMaxLen        = 255  // name character varying(255)
	databaseCatDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.database_type
	databaseTypeNameMaxLen      = 128 // name character varying(128)
	databaseTypeDbVersionMaxLen = 512 // db_version character varying(512)

	// dc.domain_cat
	domainCatDomainNameMaxLen = 100 // domain_name character varying(100)

	// dc.group_levels
	groupLevelDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.has_to_group
	hasToGroupDescriptionMaxLen = 1000 // description character varying(1000)

	// dc.host
	hostNameMaxLen        = 255  // name character varying(255)
	hostDescriptionMaxLen = 1000 // description character varying(1000)
	hostHostEnvMaxLen     = 255  // host_env character varying(255)
	hostPortEnvMaxLen     = 255  // port_env character varying(255)
	hostUsernameEnvMaxLen = 255  // username_env character varying(255)
	hostPasswordEnvMaxLen = 255  // password_env character varying(255)

	// dc.schema_cat
	schemaCatNameMaxLen = 128 // name character varying(128)

	// dc.table_cat
	tableCatNameMaxLen        = 128  // name character varying(128)
	tableCatDescriptionMaxLen = 2000 // description character varying(2000)

	// dc.table_type
	tableTypeNameMaxLen        = 128  // name character varying(128)
	tableTypeDescriptionMaxLen = 1000 // description character varying(1000)
)
