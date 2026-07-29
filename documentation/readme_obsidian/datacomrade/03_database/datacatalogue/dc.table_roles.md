# dc.table_roles

Справочник ролей на уровне таблицы.

## Columns

| Name | Type | Nullable | Default | Description |
| --- | --- | --- | --- | --- |
| `id` | `bigserial` | no | `none` | PK. |
| `name` | `varchar(128)` | no | `none` | Уникальная роль. |
| `description` | `varchar(2000)` | no | `none` | Описание. |
| `created_at` | `timestamp` | no | `now()` | Создание. |
| `updated_at` | `timestamp` | no | `now()` | Обновление. |
| `is_deleted` | `bool` | no | `false` | Флаг удаления. |

## Keys & Constraints
- Primary key: `table_roles_pk` (`id`).
- Unique: `table_roles_name_unique` (`name`).

## Foreign Keys
- Исходящие FK отсутствуют.

## Indexes
- Явно не создаются отдельными миграциями.

## Seed data
- Миграция `000000040_insert_table_roles.sql` заполняет справочник тремя ролями: `can_read`, `can_write`, `can_grant`.

## Migration source
- Create: `datacatalogue/db/migrations/000000030_create_table_table_roles.sql`
- Seed: `datacatalogue/db/migrations/000000040_insert_table_roles.sql`

## References
- Tags: #database #datacatalogue #table
- Links: [[dc.user_table_roles]]
