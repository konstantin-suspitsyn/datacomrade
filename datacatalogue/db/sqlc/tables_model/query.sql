-- =========================================================
-- dc.alias
-- =========================================================

-- name: GetAliasById :one
SELECT *
FROM dc.alias
WHERE id = $1
AND is_deleted = false;

-- name: GetAliases :many
SELECT *
FROM dc.alias
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedAliasById :one
SELECT *
FROM dc.alias
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedAliases :many
SELECT *
FROM dc.alias
WHERE is_deleted = true
ORDER BY id;

-- name: CreateAlias :one
INSERT INTO dc.alias
("name", description, created_at, updated_at, is_deleted, user_id)
VALUES($1, $2, now(), now(), false, (SELECT u.id FROM dc."user" u WHERE u.external_id = $3))
RETURNING *;

-- name: UpdateAliasById :one
UPDATE dc.alias
SET "name"=$2, description=$3, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $4), updated_at=now()
WHERE dc.alias.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteAliasById :exec
UPDATE dc.alias
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteAliasById :exec
UPDATE dc.alias
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.column_cat
-- =========================================================

-- name: GetColumnCatById :one
SELECT *
FROM dc.column_cat
WHERE id = $1
AND is_deleted = false;

-- name: GetColumnCats :many
SELECT *
FROM dc.column_cat
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedColumnCatById :one
SELECT *
FROM dc.column_cat
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedColumnCats :many
SELECT *
FROM dc.column_cat
WHERE is_deleted = true
ORDER BY id;

-- name: GetColumnCatsByTableId :many
SELECT *
FROM dc.column_cat
WHERE table_id = $1
AND is_deleted = false
ORDER BY id;

-- name: GetColumnCatsByAliasId :many
SELECT *
FROM dc.column_cat
WHERE alias_id = $1
AND is_deleted = false
ORDER BY id;

-- name: CreateColumnCat :one
INSERT INTO dc.column_cat
(table_id, "name", alias_id, column_type_id, description, calculation_type_id, is_deleted, show_in_ui, created_at, updated_at, user_id)
VALUES($1, $2, $3, $4, $5, $6, false, $7, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $8))
RETURNING *;

-- name: UpdateColumnCatById :one
UPDATE dc.column_cat
SET table_id=$2, "name"=$3, alias_id=$4, column_type_id=$5, description=$6, calculation_type_id=$7, show_in_ui=$8, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $9), updated_at=now()
WHERE dc.column_cat.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteColumnCatById :exec
UPDATE dc.column_cat
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteColumnCatById :exec
UPDATE dc.column_cat
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.column_type
-- =========================================================

-- name: GetColumnTypeById :one
SELECT *
FROM dc.column_type
WHERE id = $1
AND is_deleted = false;

-- name: GetColumnTypes :many
SELECT *
FROM dc.column_type
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedColumnTypeById :one
SELECT *
FROM dc.column_type
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedColumnTypes :many
SELECT *
FROM dc.column_type
WHERE is_deleted = true
ORDER BY id;

-- name: CreateColumnType :one
INSERT INTO dc.column_type
("name", description, is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $3))
RETURNING *;

-- name: UpdateColumnTypeById :one
UPDATE dc.column_type
SET "name"=$2, description=$3, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $4), updated_at=now()
WHERE dc.column_type.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteColumnTypeById :exec
UPDATE dc.column_type
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteColumnTypeById :exec
UPDATE dc.column_type
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.calculation_type
-- =========================================================

-- name: GetCalculationTypeById :one
SELECT *
FROM dc.calculation_type
WHERE id = $1
AND is_deleted = false;

-- name: GetCalculationTypes :many
SELECT *
FROM dc.calculation_type
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedCalculationTypeById :one
SELECT *
FROM dc.calculation_type
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedCalculationTypes :many
SELECT *
FROM dc.calculation_type
WHERE is_deleted = true
ORDER BY id;

-- name: CreateCalculationType :one
INSERT INTO dc.calculation_type
("name", description, created_at, updated_at, is_deleted)
VALUES($1, $2, now(), now(), false)
RETURNING *;

-- name: UpdateCalculationTypeById :one
UPDATE dc.calculation_type
SET "name"=$2, description=$3, updated_at=now()
WHERE id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteCalculationTypeById :exec
UPDATE dc.calculation_type
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteCalculationTypeById :exec
UPDATE dc.calculation_type
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.database_calculation
-- =========================================================

-- name: GetDatabaseCalculationById :one
SELECT *
FROM dc.database_calculation
WHERE id = $1
AND is_deleted = false;

-- name: GetDatabaseCalculations :many
SELECT *
FROM dc.database_calculation
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedDatabaseCalculationById :one
SELECT *
FROM dc.database_calculation
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedDatabaseCalculations :many
SELECT *
FROM dc.database_calculation
WHERE is_deleted = true
ORDER BY id;

-- name: CreateDatabaseCalculation :one
INSERT INTO dc.database_calculation
(database_cat_id, calculation_type_id, created_at, updated_at, is_deleted, user_id)
VALUES($1, $2, now(), now(), false, (SELECT u.id FROM dc."user" u WHERE u.external_id = $3))
RETURNING *;

-- name: UpdateDatabaseCalculationById :one
UPDATE dc.database_calculation
SET database_cat_id=$2, calculation_type_id=$3, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $4), updated_at=now()
WHERE dc.database_calculation.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteDatabaseCalculationById :exec
UPDATE dc.database_calculation
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteDatabaseCalculationById :exec
UPDATE dc.database_calculation
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.database_cat
-- =========================================================

-- name: GetDatabaseCatById :one
SELECT *
FROM dc.database_cat
WHERE id = $1
AND is_deleted = false;

-- name: GetDatabaseCats :many
SELECT *
FROM dc.database_cat
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedDatabaseCatById :one
SELECT *
FROM dc.database_cat
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedDatabaseCats :many
SELECT *
FROM dc.database_cat
WHERE is_deleted = true
ORDER BY id;

-- name: GetDatabaseCatsByHostId :many
SELECT *
FROM dc.database_cat
WHERE host_id = $1
AND is_deleted = false
ORDER BY id;

-- name: GetDatabaseCatsByDatabaseTypeId :many
SELECT *
FROM dc.database_cat
WHERE database_type_id = $1
AND is_deleted = false
ORDER BY id;

-- name: CreateDatabaseCat :one
INSERT INTO dc.database_cat
("name", host_id, database_type_id, description, is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, $3, $4, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $5))
RETURNING *;

-- name: UpdateDatabaseCatById :one
UPDATE dc.database_cat
SET "name"=$2, host_id=$3, database_type_id=$4, description=$5, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $6), updated_at=now()
WHERE dc.database_cat.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteDatabaseCatById :exec
UPDATE dc.database_cat
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteDatabaseCatById :exec
UPDATE dc.database_cat
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.database_type
-- =========================================================

-- name: GetDatabaseTypeById :one
SELECT *
FROM dc.database_type
WHERE id = $1
AND is_deleted = false;

-- name: GetDatabaseTypes :many
SELECT *
FROM dc.database_type
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedDatabaseTypeById :one
SELECT *
FROM dc.database_type
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedDatabaseTypes :many
SELECT *
FROM dc.database_type
WHERE is_deleted = true
ORDER BY id;

-- name: CreateDatabaseType :one
INSERT INTO dc.database_type
("name", db_version, is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $3))
RETURNING *;

-- name: UpdateDatabaseTypeById :one
UPDATE dc.database_type
SET "name"=$2, db_version=$3, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $4), updated_at=now()
WHERE dc.database_type.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteDatabaseTypeById :exec
UPDATE dc.database_type
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteDatabaseTypeById :exec
UPDATE dc.database_type
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.domain_cat
-- =========================================================

-- name: GetDomainCatById :one
SELECT *
FROM dc.domain_cat
WHERE id = $1
AND is_deleted = false;

-- name: GetDomainCats :many
SELECT *
FROM dc.domain_cat
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedDomainCatById :one
SELECT *
FROM dc.domain_cat
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedDomainCats :many
SELECT *
FROM dc.domain_cat
WHERE is_deleted = true
ORDER BY id;

-- name: CreateDomainCat :one
INSERT INTO dc.domain_cat
(domain_name, is_deleted, created_at, updated_at, user_id)
VALUES($1, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $2))
RETURNING *;

-- name: UpdateDomainCatById :one
UPDATE dc.domain_cat
SET domain_name=$2, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $3), updated_at=now()
WHERE dc.domain_cat.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteDomainCatById :exec
UPDATE dc.domain_cat
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteDomainCatById :exec
UPDATE dc.domain_cat
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.following_calculation
-- =========================================================

-- name: GetFollowingCalculationById :one
SELECT *
FROM dc.following_calculation
WHERE id = $1
AND is_deleted = false;

-- name: GetFollowingCalculations :many
SELECT *
FROM dc.following_calculation
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedFollowingCalculationById :one
SELECT *
FROM dc.following_calculation
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedFollowingCalculations :many
SELECT *
FROM dc.following_calculation
WHERE is_deleted = true
ORDER BY id;

-- name: CreateFollowingCalculation :one
INSERT INTO dc.following_calculation
(column_cat_id, calculation_type_id, created_at, updated_at, is_deleted, user_id)
VALUES($1, $2, now(), now(), false, (SELECT u.id FROM dc."user" u WHERE u.external_id = $3))
RETURNING *;

-- name: UpdateFollowingCalculationById :one
UPDATE dc.following_calculation
SET column_cat_id=$2, calculation_type_id=$3, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $4), updated_at=now()
WHERE dc.following_calculation.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteFollowingCalculationById :exec
UPDATE dc.following_calculation
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteFollowingCalculationById :exec
UPDATE dc.following_calculation
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.group_levels
-- =========================================================

-- name: GetGroupLevelById :one
SELECT *
FROM dc.group_levels
WHERE id = $1
AND is_deleted = false;

-- name: GetGroupLevels :many
SELECT *
FROM dc.group_levels
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedGroupLevelById :one
SELECT *
FROM dc.group_levels
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedGroupLevels :many
SELECT *
FROM dc.group_levels
WHERE is_deleted = true
ORDER BY id;

-- name: CreateGroupLevel :one
INSERT INTO dc.group_levels
(column_id, parent_column_id, "level", description, created_at, updated_at, is_deleted, user_id)
VALUES($1, $2, $3, $4, now(), now(), false, (SELECT u.id FROM dc."user" u WHERE u.external_id = $5))
RETURNING *;

-- name: UpdateGroupLevelById :one
UPDATE dc.group_levels
SET column_id=$2, parent_column_id=$3, "level"=$4, description=$5, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $6), updated_at=now()
WHERE dc.group_levels.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteGroupLevelById :exec
UPDATE dc.group_levels
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteGroupLevelById :exec
UPDATE dc.group_levels
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.has_to_group
-- =========================================================

-- name: GetHasToGroupById :one
SELECT *
FROM dc.has_to_group
WHERE id = $1
AND is_deleted = false;

-- name: GetHasToGroups :many
SELECT *
FROM dc.has_to_group
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedHasToGroupById :one
SELECT *
FROM dc.has_to_group
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedHasToGroups :many
SELECT *
FROM dc.has_to_group
WHERE is_deleted = true
ORDER BY id;

-- name: CreateHasToGroup :one
INSERT INTO dc.has_to_group
(column_id_a, column_id_b, description, is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, $3, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $4))
RETURNING *;

-- name: UpdateHasToGroupById :one
UPDATE dc.has_to_group
SET column_id_a=$2, column_id_b=$3, description=$4, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $5), updated_at=now()
WHERE dc.has_to_group.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteHasToGroupById :exec
UPDATE dc.has_to_group
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteHasToGroupById :exec
UPDATE dc.has_to_group
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.host
-- =========================================================

-- name: GetHostById :one
SELECT *
FROM dc.host
WHERE id = $1
AND is_deleted = false;

-- name: GetHosts :many
SELECT *
FROM dc.host
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedHostById :one
SELECT *
FROM dc.host
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedHosts :many
SELECT *
FROM dc.host
WHERE is_deleted = true
ORDER BY id;

-- name: CreateHost :one
INSERT INTO dc.host
("name", description, host_env, port_env, username_env, password_env, is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, $3, $4, $5, $6, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $7))
RETURNING *;

-- name: UpdateHostById :one
UPDATE dc.host
SET "name"=$2, description=$3, host_env=$4, port_env=$5, username_env=$6, password_env=$7, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $8), updated_at=now()
WHERE dc.host.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteHostById :exec
UPDATE dc.host
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteHostById :exec
UPDATE dc.host
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.schema_cat
-- =========================================================

-- name: GetSchemaCatById :one
SELECT *
FROM dc.schema_cat
WHERE id = $1
AND is_deleted = false;

-- name: GetSchemaCats :many
SELECT *
FROM dc.schema_cat
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedSchemaCatById :one
SELECT *
FROM dc.schema_cat
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedSchemaCats :many
SELECT *
FROM dc.schema_cat
WHERE is_deleted = true
ORDER BY id;

-- name: GetSchemaCatsByDatabaseId :many
SELECT *
FROM dc.schema_cat
WHERE database_id = $1
AND is_deleted = false
ORDER BY id;

-- name: CreateSchemaCat :one
INSERT INTO dc.schema_cat
(database_id, "name", is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $3))
RETURNING *;

-- name: UpdateSchemaCatById :one
UPDATE dc.schema_cat
SET database_id=$2, "name"=$3, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $4), updated_at=now()
WHERE dc.schema_cat.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteSchemaCatById :exec
UPDATE dc.schema_cat
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteSchemaCatById :exec
UPDATE dc.schema_cat
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.table_cat
-- =========================================================

-- name: GetTableCatById :one
SELECT *
FROM dc.table_cat
WHERE id = $1
AND is_deleted = false;

-- name: GetTableCats :many
SELECT *
FROM dc.table_cat
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedTableCatById :one
SELECT *
FROM dc.table_cat
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedTableCats :many
SELECT *
FROM dc.table_cat
WHERE is_deleted = true
ORDER BY id;

-- name: GetTableCatsBySchemaId :many
SELECT *
FROM dc.table_cat
WHERE schema_id = $1
AND is_deleted = false
ORDER BY id;

-- name: GetTableCatsByTableTypeId :many
SELECT *
FROM dc.table_cat
WHERE table_type_id = $1
AND is_deleted = false
ORDER BY id;

-- name: GetTableCatsByDomainId :many
SELECT *
FROM dc.table_cat
WHERE domain_id = $1
AND is_deleted = false
ORDER BY id;

-- name: CreateTableCat :one
INSERT INTO dc.table_cat
("name", description, schema_id, table_type_id, domain_id, is_deleted, created_at, updated_at, is_get_dict, user_id)
VALUES($1, $2, $3, $4, $5, false, now(), now(), $6, (SELECT u.id FROM dc."user" u WHERE u.external_id = $7))
RETURNING *;

-- name: UpdateTableCatById :one
UPDATE dc.table_cat
SET "name"=$2, description=$3, schema_id=$4, table_type_id=$5, domain_id=$6, is_get_dict=$7, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $8), updated_at=now()
WHERE dc.table_cat.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteTableCatById :exec
UPDATE dc.table_cat
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteTableCatById :exec
UPDATE dc.table_cat
SET is_deleted=false, updated_at=now()
WHERE id=$1;

-- =========================================================
-- dc.table_type
-- =========================================================

-- name: GetTableTypeById :one
SELECT *
FROM dc.table_type
WHERE id = $1
AND is_deleted = false;

-- name: GetTableTypes :many
SELECT *
FROM dc.table_type
WHERE is_deleted = false
ORDER BY id;

-- name: GetDeletedTableTypeById :one
SELECT *
FROM dc.table_type
WHERE id = $1
AND is_deleted = true;

-- name: GetDeletedTableTypes :many
SELECT *
FROM dc.table_type
WHERE is_deleted = true
ORDER BY id;

-- name: CreateTableType :one
INSERT INTO dc.table_type
("name", description, is_deleted, created_at, updated_at, user_id)
VALUES($1, $2, false, now(), now(), (SELECT u.id FROM dc."user" u WHERE u.external_id = $3))
RETURNING *;

-- name: UpdateTableTypeById :one
UPDATE dc.table_type
SET "name"=$2, description=$3, user_id=(SELECT u.id FROM dc."user" u WHERE u.external_id = $4), updated_at=now()
WHERE dc.table_type.id=$1
AND is_deleted = false
RETURNING *;

-- name: DeleteTableTypeById :exec
UPDATE dc.table_type
SET is_deleted=true, updated_at=now()
WHERE id=$1;

-- name: UndeleteTableTypeById :exec
UPDATE dc.table_type
SET is_deleted=false, updated_at=now()
WHERE id=$1;
