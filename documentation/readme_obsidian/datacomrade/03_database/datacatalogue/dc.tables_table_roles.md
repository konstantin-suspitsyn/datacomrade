# dc.tables_table_roles

Таблица связей таблиц и ролей таблиц (many-to-many).

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `table_cat_id` | `bigint` | no | `none` | Таблица. |
| `table_roles_id` | `bigint` | no | `none` | Роль таблицы. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |

## Keys & Constraints
- Primary key: `tables_user_table_roles_pk` (`id`) (имя в create отличается от down-мigration).

## Foreign Keys
- `tables_table_roles_table_cat_fk`: `table_cat_id` -> [[dc.table_cat]](`dc.table_cat.id`) `ON UPDATE CASCADE`.
- `tables_table_roles_table_roles_fk`: `table_roles_id` -> [[dc.table_roles]](`dc.table_roles.id`) `ON DELETE RESTRICT ON UPDATE CASCADE`.

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000033_create_table_tables_table_roles.sql`
- FK/Alter: `datacatalogue/db/migrations/000000036_add_fk_tables_table_roles.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.table_cat]] [[dc.table_roles]]
