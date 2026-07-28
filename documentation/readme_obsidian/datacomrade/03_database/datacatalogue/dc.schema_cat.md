# dc.schema_cat

Каталог схем БД.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `database_id` | `int8` | no | `none` | База данных. |
| `name` | `varchar(128)` | no | `none` | Имя схемы. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `user_id` | `int8` | no | `none` | Пользователь. |

## Keys & Constraints
- Primary key: `schema_cat_pk` (`id`).

## Foreign Keys
- `schema_cat_database_cat_fk`: `database_id` -> [[dc.database_cat]](`dc.database_cat.id`).
- `schema_cat_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000010_create_table_schema_cat.sql`
- FK/Alter: `datacatalogue/db/migrations/000000015_add_fk_schema_cat.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.database_cat]] [[dc.user]] [[dc.table_cat]]
