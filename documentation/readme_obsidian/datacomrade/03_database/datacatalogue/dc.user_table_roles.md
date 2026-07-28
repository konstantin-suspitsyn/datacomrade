# dc.user_table_roles

Связка пользователей и ролей таблиц.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `user_id` | `bigint` | no | `none` | Пользователь. |
| `table_roles_id` | `bigint` | no | `none` | Роль таблицы. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |

## Keys & Constraints
- Primary key: `user_table_roles_pk` (`id`).

## Foreign Keys
- `user_table_roles_user_fk`: `user_id` -> [[dc.user]](`dc."user".id`).
- `user_table_roles_table_roles_fk`: `table_roles_id` -> [[dc.table_roles]](`dc.table_roles.id`).

## Indexes
- Явно не создаются отдельными миграциями.

## Migration source
- Create: `datacatalogue/db/migrations/000000032_create_table_user_table_roles.sql`
- FK/Alter: `datacatalogue/db/migrations/000000038_add_fk_user_table_roles.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user]] [[dc.table_roles]]
